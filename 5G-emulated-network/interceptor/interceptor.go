package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/vishvananda/netlink"
	"gopkg.in/yaml.v3"
)

type HostapdInterceptor struct {
	attached bool
	conn     *net.UnixConn
	remote   *net.UnixAddr
	quit     chan struct{} // Channel to signal termination
}

// Session represents a single PDU session
type Session struct {
	Id          int
	State       string `yaml:"state"`
	SessionType string `yaml:"session-type"`
	APN         string `yaml:"apn"`
	SNSSAI      struct {
		SST string `yaml:"sst"`
		SD  string `yaml:"sd"`
	} `yaml:"s-nssai"`
	Emergency   bool   `yaml:"emergency"`
	Address     string `yaml:"address"`
	AMBR        string `yaml:"ambr"`
	DataPending bool   `yaml:"data-pending"`
}

type Lease struct {
	expiration int
	counter    int
}

type Device struct {
	state       string
	lease       Lease
	pdu_session *Session
}

var (
	logger              *log.Logger
	allowed_devices     = make(map[string]Device)
	wg                  sync.WaitGroup
	hostapd_interceptor *HostapdInterceptor
)

func SetLogging(mode string) {
	if mode == "debug" {
		logger = log.New(os.Stdout, "", log.LstdFlags) // Use stdout
	} else {
		w, err := syslog.New(syslog.LOG_INFO|syslog.LOG_DAEMON, "interceptor")
		if err != nil {
			fmt.Println("Failed to connect to syslog: ", err)
			os.Exit(1)
		}
		logger = log.New(w, "", 0) // Use syslog
	}
}

func NewInterceptor(hostapdSocketPath string) (*HostapdInterceptor, error) {
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
	interceptor := &HostapdInterceptor{
		conn:     localConn,  // Local socket
		remote:   remoteAddr, // Remote address
		attached: false,
		quit:     make(chan struct{}), // Initialize quit channel
	}

	return interceptor, nil
}

func (i *HostapdInterceptor) Request(command []byte) ([]byte, error) {
	_, err := i.conn.WriteToUnix(command, i.remote)
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

func (i *HostapdInterceptor) Attach() error {
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

func (i *HostapdInterceptor) Deauth(mac_address string) error {
	_, err := hostapd_interceptor.Request([]byte(fmt.Sprintf("DEAUTHENTICATE %s", mac_address)))
	if err != nil {
		return err
	}
	return nil
}

// Shutdown cleans up resources gracefully
func (i *HostapdInterceptor) Shutdown() {
	close(i.quit)  // Signal the listener to stop
	i.conn.Close() // Close the socket
}

func AllowMac(allowed_macs_file string, mac string) error {
	// Open file in append mode, create if not exists
	f, err := os.OpenFile(allowed_macs_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	// Write to file
	entry := fmt.Sprintf("dhcp-host=%s,set:known\n", mac)
	if _, err := f.WriteString(entry); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	log.Printf("MAC address %v added to allowed list successfully!", mac)

	return nil
}

func DisallowMac(allowed_macs_file string, mac string) error {
	// Check if file exists, if not return error
	fileContent, err := os.ReadFile(allowed_macs_file)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	// Split the file content into lines
	lines := strings.Split(string(fileContent), "\n")

	// Filter out the line to remove
	var newLines []string
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if !strings.Contains(trimmedLine, mac) {
			newLines = append(newLines, line)
		} else {
			newLines = append(newLines, "")
		}
	}

	// Join the lines back into a single string
	newContent := strings.Join(newLines, "\n")

	// Write the new content back to the file
	err = os.WriteFile(allowed_macs_file, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

// restartDnsmasq runs the systemctl command to restart dnsmasq
func RestartDnsmasq() error {
	cmd := exec.Command("systemctl", "restart", "dnsmasq")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restart dnsmasq: %w", err)
	}
	logger.Println("dnsmasq restarted successfully")
	return nil
}

func ForgetDevice(allowed_macs_file string, leases_file string, mac_address string, ue_imsi string) error {
	// Disallow mac for DHCP offers
	logger.Println("Disallowing device: ", mac_address)
	DisallowMac(allowed_macs_file, mac_address)

	// Delete lease
	DisallowMac(leases_file, mac_address)

	// Restarting dnsmasq
	err := RestartDnsmasq()
	if err != nil {
		return fmt.Errorf("error restarting dnsmasq: %w", err)
	}

	//De-auth device from hostapd
	logger.Println("Deauthenticating device: ", mac_address)
	err = hostapd_interceptor.Deauth(mac_address)
	if err != nil {
		return fmt.Errorf("DEAUTHENTICATE request for %v failed: %v", mac_address, err)
	}

	// Releasing PDU Sessions
	logger.Println("Releasing PDU Session from device: ", mac_address)
	err = ReleasePDUSession(ue_imsi, allowed_devices[mac_address].pdu_session.Id)
	if err != nil {
		logger.Printf("RELEASE request for %v failed: %v", mac_address, err)
	}

	// Forgetting device
	logger.Println("Forgetting device: ", mac_address)
	delete(allowed_devices, mac_address)

	return nil
}

// HostapdListener continuously listens for incoming messages
func HostapdListener(allowed_macs_file string, ue_imsi string) {
	defer wg.Done() // Mark as done when exiting

	logger.Println("Listening for incoming messages...")

	buf := make([]byte, 4096)
	for {
		select {
		case <-hostapd_interceptor.quit:
			logger.Println("Stopping listener...")
			return // Exit loop on quit signal

		default:
			// Blocking read
			n, _, err := hostapd_interceptor.conn.ReadFromUnix(buf)
			if err != nil {
				select {
				case <-hostapd_interceptor.quit:
					// Suppress error if quitting
					return
				default:
					logger.Println("Error receiving message: ", err)
				}
			} else {
				logger.Println(string(buf[3:n]))
				if strings.Contains(string(buf[:n]), "CTRL-EVENT-EAP-SUCCESS") {
					// Extract MAC address
					parts := strings.Split(string(buf[:n]), " ")
					if len(parts) < 2 {
						logger.Println("Could not extract MAC address from event message.")
						return
					}
					mac_address := parts[1]

					logger.Println("Authentication success for client: ", mac_address)

					logger.Println("Requesting PDU Session for client: ", mac_address)
					allowed_devices[mac_address] = Device{
						state: "AUTHENTICATED",
						lease: Lease{
							counter:    0,
							expiration: 0,
						},
						pdu_session: func() *Session {
							session, err := NewPDUSession(ue_imsi)
							if err != nil {
								logger.Println("Error creating PDU session: ", err)
								return nil
							}
							return session
						}(),
					}

					if allowed_devices[mac_address].pdu_session != nil {
						// Add MAC to allowed allowed_devices list
						AllowMac(allowed_macs_file, mac_address)
					} else {
						logger.Println("Couldn't provide PDU Session for device: ", mac_address)
						//De-auth device from hostapd
						logger.Println("Deauthenticating device: ", mac_address)
						err := hostapd_interceptor.Deauth(mac_address)
						if err != nil {
							logger.Printf("DEAUTHENTICATE request for %v failed: %v\n", mac_address, err)
						}
						delete(allowed_devices, mac_address)
					}

					// Restart dnsmasq service
					err = RestartDnsmasq()
					if err != nil {
						logger.Println("Error restarting dnsmasq: ", err)
					}
				}
			}
		}
	}
}

// Listener in the lease file (/var/lib/misc/dnsmasq.leases)
func DnsmasqListener(allowed_macs_file string, leases_file string, ue_imsi string) {
	defer wg.Done()

	start_info, err := os.Stat(leases_file)
	if err != nil {
		logger.Println("Error retrieving leases file stats: ", err)
	}

	for {
		temp_info, err := os.Stat(leases_file)

		if err != nil {
			logger.Println("Error retrieving leases file stats: ", err)
			return
		}

		if temp_info.ModTime() != start_info.ModTime() {
			// Check if file exists, if not return error
			file, err := os.Open(leases_file) // Open the file
			if err != nil {
				logger.Println("Error opening leases file: ", err)
			}

			scanner := bufio.NewScanner(file) // Create a new scanner
			for scanner.Scan() {              // Read line by line
				fields := strings.Fields(scanner.Text())
				if len(fields) < 2 {
					continue
				}

				expirationStr, mac_address := fields[0], fields[1]

				expiration, err := strconv.Atoi(expirationStr)
				if err != nil {
					logger.Println("Error converting expiration to int: ", err)
					continue
				}

				if device, exists := allowed_devices[mac_address]; exists {
					if device.lease.expiration != expiration {
						device.lease.counter++ // Update counter
						device.lease.expiration = expiration
						allowed_devices[mac_address] = device // Save back to map

						logger.Printf("Lease #%v device: %v", device.lease.counter, mac_address)
					}
				}
			}
			file.Close() // Explicit close
			start_info = temp_info
		}
	}
}

func HostDisconnectListener(allowed_macs_file string, leases_file string, ue_imsi string, link netlink.Link) {
	defer wg.Done()

	for {
		// Get neighbor entries for that link
		neighs, err := netlink.NeighList(link.Attrs().Index, netlink.FAMILY_V4)
		if err != nil {
			logger.Println(err)
			return
		}

		for _, neigh := range neighs {
			mac := neigh.HardwareAddr.String()

			if device, exists := allowed_devices[mac]; exists {
				if neigh.State == netlink.NUD_REACHABLE && device.state == "AUTHENTICATED" {
					logger.Println("Device authenticated and reachable:", mac)
					// Update device state
					device.state = "REACHABLE"
					allowed_devices[mac] = device
				} else if neigh.State == netlink.NUD_STALE {
					// Device is not reachable
					logger.Println("Device is stale:", mac)
					err = ForgetDevice(allowed_macs_file, leases_file, mac, ue_imsi)
					if err != nil {
						logger.Println("Error forgetting device:", err)
					}
				}
			}
		}
	}
}

// NewPDUSession establishes a new session and waits until it is active
func NewPDUSession(ue_imsi string) (*Session, error) {
	// Establish a new PDU session
	cmd := exec.Command("nr-cli", ue_imsi, "--exec", "ps-establish IPv4 --sst 1 --dnn clients")
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to establish new PDU Session: %w", err)
	}
	log.Println("New PDU Session requested successfully, waiting for activation...")

	// Wait for session to become active
	const maxRetries = 20
	const sleepInterval = 3 * time.Second

	var session *Session
	var err error

	for i := 0; i < maxRetries; i++ {
		session, err = LastPDUSession(ue_imsi)
		if err != nil {
			log.Printf("Retrying (%d/%d): Failed to retrieve session: %v", i+1, maxRetries, err)
		} else if session != nil && session.State == "PS-ACTIVE" && session.Address != "" {
			log.Println("PDU Session is now active!")
			return session, nil
		} else {
			log.Printf("Retrying (%d/%d): Waiting for session activation (State: %s, Address: %s)...", i+1, maxRetries, session.State, session.Address)
		}

		time.Sleep(sleepInterval)
	}

	if session != nil {
		ReleasePDUSession(ue_imsi, session.Id)
	}

	return nil, fmt.Errorf("PDU Session did not become active within timeout")
}

// LastPDUSession retrieves the latest session from ps-list
func LastPDUSession(ue_imsi string) (*Session, error) {
	cmd := exec.Command("nr-cli", ue_imsi, "--exec", "ps-list")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to retrieve PDU Session listing: %w", err)
	}

	// Parse output as YAML
	var pduSessions map[string]Session
	if err := yaml.Unmarshal(out.Bytes(), &pduSessions); err != nil {
		return nil, fmt.Errorf("failed to parse PDU Session listing as YAML: %w", err)
	}

	// Extract session IDs and sort them
	var sessionKeys []int
	sessionMap := make(map[int]Session)

	for key, session := range pduSessions {
		sessionNumStr := strings.TrimPrefix(key, "PDU Session")
		sessionNum, err := strconv.Atoi(sessionNumStr)
		if err != nil {
			return nil, fmt.Errorf("invalid session ID: %s", key)
		}

		session.Id = sessionNum
		sessionMap[sessionNum] = session
		sessionKeys = append(sessionKeys, sessionNum)
	}

	// If no sessions were found, return an error
	if len(sessionKeys) == 0 {
		return nil, fmt.Errorf("no PDU sessions found")
	}

	// Sort and get the last session
	sort.Ints(sessionKeys)
	lastSessionID := sessionKeys[len(sessionKeys)-1]
	lastSession := sessionMap[lastSessionID]

	return &lastSession, nil
}

// Release PDU Session
func ReleasePDUSession(ue_imsi string, pdu_id int) error {
	// Establish a new PDU session
	cmd := exec.Command("nr-cli", ue_imsi, "--exec", "ps-release ", strconv.Itoa(pdu_id))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to release PDU Session #%v: %w", pdu_id, err)
	}
	log.Println("PDU Session released successfully")

	return nil
}

func main() {
	mode := flag.String("mode", "syslog", "Logging mode: syslog or debug")
	hostpad_int := flag.String("interface", "/var/run/hostapd/enp0s9", "Hostapd socket path")
	allowed_macs_file := flag.String("allowed", "/etc/allowed-macs.conf", "Dnsmasq allowed MACs file")
	leases_file := flag.String("leases", "/var/lib/misc/dnsmasq.leases", "Dnsmasq DHCP leases files")
	ue_imsi := flag.String("imsi", "imsi-999700000000001", "UE IMSI")
	flag.Parse() // Parse command-line flags

	SetLogging(*mode)

	// Create link to the network interface, using the provided interface name
	sl := strings.Split(*hostpad_int, "/")
	linkName := sl[len(sl)-1]
	link, err := netlink.LinkByName(linkName)
	if err != nil {
		log.Println("Error retrieving link by name: ", err)
		return
	}

	// Create interceptor
	hostapd_interceptor, err = NewInterceptor(*hostpad_int)
	if err != nil {
		log.Println("Error starting interceptor: ", err)
		return
	}
	defer hostapd_interceptor.Shutdown()
	logger.Println("Interceptor started!")

	// Attach to the external service
	err = hostapd_interceptor.Attach()
	if err != nil {
		logger.Println("Error attaching to the hostapd: ", err)
		return
	}
	logger.Println("Attached to ", *hostpad_int)

	// Testing hostapd_cli request
	_, err = hostapd_interceptor.Request([]byte(fmt.Sprintf("STATUS")))
	if err != nil {
		logger.Printf("Request for STATUS failed: %v", err)
	}

	// Start listening in a separate goroutine
	wg.Add(1)
	go HostapdListener(*allowed_macs_file, *ue_imsi)

	// Start listening for DHCP lease renewals
	wg.Add(1)
	go DnsmasqListener(*allowed_macs_file, *leases_file, *ue_imsi)

	// Start listening for host disconnects
	wg.Add(1)
	go HostDisconnectListener(*allowed_macs_file, *leases_file, *ue_imsi, link)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for termination signal
	<-sigChan
	logger.Println("Shutting down gracefully...")
}
