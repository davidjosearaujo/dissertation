package main

import (
	"bufio"
	"fmt"
	"log"
	"log/syslog"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type Interception interface {
	// Intercept packets
	Attach() error
	Request(data []byte) ([]byte, error)
	Receive() ([]byte, error)
	Shutdown()
}

type Interceptor struct {
	attached bool
	conn    *net.UnixConn
	remote  *net.UnixAddr
	quit    chan struct{} // Channel to signal termination
	wg      sync.WaitGroup
}

var logger *syslog.Writer

func init() {
	var err error
	logger, err = syslog.New(syslog.LOG_INFO|syslog.LOG_DAEMON, "interceptor")
	if err != nil {
		fmt.Println("Failed to connect to syslog:", err)
		os.Exit(1)
	}
	log.SetOutput(logger) // Redirect standard Go log package to syslog
}

func NewInterceptor(hostapdSocketPath string) (*Interceptor, error) {
	// Create a unique local socket path
	localPath := fmt.Sprintf("/tmp/interceptor_%d.sock", os.Getpid())

	// Remove any existing socket file
	if err := os.Remove(localPath); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error removing old socket: %v", err)
	}

	// Resolve local address
	localAddr, err := net.ResolveUnixAddr("unixgram", localPath)
	if err != nil {
		return nil, fmt.Errorf("error resolving local address: %v", err)
	}

	// Resolve remote address
	remoteAddr, err := net.ResolveUnixAddr("unixgram", hostapdSocketPath)
	if err != nil {
		return nil, fmt.Errorf("error resolving remote address: %v", err)
	}

	// Create and bind the local socket
	localConn, err := net.ListenUnixgram("unixgram", localAddr)
	if err != nil {
		return nil, fmt.Errorf("error binding local socket: %v", err)
	}

	// Return the interceptor
	interceptor := &Interceptor{
		conn:    localConn,  // Local socket
		remote:  remoteAddr, // Remote address
		attached: false,
		quit:    make(chan struct{}), // Initialize quit channel
	}

	return interceptor, nil
}


func (i *Interceptor) Request(data []byte) ([]byte, error) {
	_, err := i.conn.WriteToUnix(data, i.remote)
	if err != nil {
		return nil, fmt.Errorf("error sending command: %w", err)
	}

	// Read response with timeout handling
	reader := bufio.NewReader(i.conn)
	buf := make([]byte, 4096)

	n, err := reader.Read(buf)
	if err != nil {
		if os.IsTimeout(err) {
			return nil, fmt.Errorf("timeout reached, no response received")
		} else {
			return nil, fmt.Errorf("error reading response: %w", err)
		}
	}

	return buf[:n], nil
}

func (i *Interceptor) Receive() ([]byte, error) {
	buf := make([]byte, 4096)
	n, _, err := i.conn.ReadFromUnix(buf)
	if err != nil {
		return nil, fmt.Errorf("error receiving response: %w", err)
	}

	return buf[:n], nil
}

func (i *Interceptor) Attach() error {
	// Check if already attached
	if i.attached {
		return nil
	}

	// Send "ATTACH" command and get response
	res, err := i.Request([]byte("ATTACH"))
	if err != nil {
		return fmt.Errorf("ATTACH request failed: %w", err)
	}

	// Check if response contains "OK"
	if strings.Contains(string(res), "OK") {
		i.attached = true
		return nil
	}

	// If response is not "OK", return an error
	return fmt.Errorf("ATTACH failed")
}

// Shutdown cleans up resources gracefully
func (i *Interceptor) Shutdown() {
	close(i.quit)     // Signal the listener to stop
	i.wg.Wait()       // Wait for goroutines to finish
	i.conn.Close()   // Close the socket
}

func UpdateAllowedMACs(mac string) error {
	// Open file in append mode, create if not exists
	f, err := os.OpenFile("/etc/allowed-macs.conf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	// Write to file
	entry := fmt.Sprintf("dhcp-host=%s,set:known\n", mac)
	if _, err := f.WriteString(entry); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	log.Println("MAC address written successfully:", entry)

	return nil
}

// restartDnsmasq runs the systemctl command to restart dnsmasq
func RestartDnsmasq() error {
	cmd := exec.Command("systemctl", "restart", "dnsmasq")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ListenLoop continuously listens for incoming messages
func (i *Interceptor) ListenLoop() {
	defer i.wg.Done() // Mark as done when exiting

	log.Println("Listening for incoming messages...")

	buf := make([]byte, 4096)
	for {
		select {
		case <-i.quit:
			log.Println("Stopping listener...")
			return // Exit loop on quit signal

		default:
			// Blocking read
			n, _, err := i.conn.ReadFromUnix(buf)
			if err != nil {
				select {
				case <-i.quit:
					// Suppress error if quitting
					return
				default:
					log.Println("Error receiving message:", err)
				}
			} else {
				log.Println("Received:", string(buf[:n]))
				if strings.Contains(string(buf[:n]), "CTRL-EVENT-EAP-SUCCESS") {
					log.Println("EAP success detected for client: ", string(buf[:n]))

					// Extract MAC address
					parts := strings.Split(string(buf[:n]), " ")
					if len(parts) < 2 {
						log.Println("Could not extract MAC address from event message.")
						return
					}
					
					UpdateAllowedMACs(parts[1])

					err = RestartDnsmasq()
					if err != nil {
						log.Println("Error restarting dnsmasq:", err)
					}

					log.Println("dnsmasq service restarted successfully.")
				}
			}
		}
	}
}

func main() {
	// Define paths
	serverPath := "/var/run/hostapd/enp0s10" // Example external service

	// Create interceptor
	interceptor, err := NewInterceptor(serverPath)
	if err != nil {
		log.Println("Error creating interceptor:", err)
		return
	}
	defer interceptor.Shutdown()
	log.Println("Interceptor created")

	// Attach to the external service
	err = interceptor.Attach()
	if err != nil {
		log.Println("Error attaching to the external service:", err)
		return
	}
	log.Println("Attached to", serverPath)

	// Start listening in a separate goroutine
	interceptor.wg.Add(1)
	go interceptor.ListenLoop()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for termination signal
	<-sigChan
	log.Println("Shutting down gracefully...")
}
