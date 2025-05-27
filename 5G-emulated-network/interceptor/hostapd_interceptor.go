// Copyright 2025 David Araújo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	hostapdRequestTimeout      = 5 * time.Second
	hostapdListenerReadTimeout = 1 * time.Second
	hostapdSocketBufferSize    = 4096
	hostapdCmdAttach           = "ATTACH"
	hostapdCmdDeauthenticate   = "DEAUTHENTICATE"
	hostapdEventEAPSuccess     = "CTRL-EVENT-EAP-SUCCESS"
	hostapdResponseOK          = "OK"
)

// HostapdInterceptor manages communication with the hostapd control interface.
type HostapdInterceptor struct {
	attached bool
	conn     *net.UnixConn
	remote   *net.UnixAddr
}

// NewInterceptor creates and initializes a new HostapdInterceptor.
func NewInterceptor(hostapdSocketPath string) (*HostapdInterceptor, error) {
	localPath := fmt.Sprintf("/tmp/interceptor_%d.sock", os.Getpid())
	if err := os.Remove(localPath); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("removing old socket %s: %v", localPath, err)
	}

	localAddr, err := net.ResolveUnixAddr("unixgram", localPath)
	if err != nil {
		return nil, fmt.Errorf("resolving local addr %s: %v", localPath, err)
	}
	remoteAddr, err := net.ResolveUnixAddr("unixgram", hostapdSocketPath)
	if err != nil {
		return nil, fmt.Errorf("resolving remote addr %s: %v", hostapdSocketPath, err)
	}

	localConn, err := net.ListenUnixgram("unixgram", localAddr)
	if err != nil {
		return nil, fmt.Errorf("binding local socket %s: %v", localPath, err)
	}

	return &HostapdInterceptor{
		conn:     localConn,
		remote:   remoteAddr,
		attached: false,
	}, nil
}

// Request sends a command to hostapd and returns the response.
func (i *HostapdInterceptor) Request(command []byte) ([]byte, error) {
	if i.conn == nil {
		return nil, fmt.Errorf("interceptor conn not initialized")
	}
	_, err := i.conn.WriteToUnix(command, i.remote)
	if err != nil {
		return nil, fmt.Errorf("sending cmd to hostapd: %w", err)
	}

	buf := make([]byte, hostapdSocketBufferSize)
	if err := i.conn.SetReadDeadline(time.Now().Add(hostapdRequestTimeout)); err != nil {
		return nil, fmt.Errorf("setting read deadline: %w", err)
	}

	n, err := i.conn.Read(buf)
	_ = i.conn.SetReadDeadline(time.Time{})

	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, fmt.Errorf("timeout reading response for '%s'", string(command))
		}
		return nil, fmt.Errorf("reading response from hostapd: %w", err)
	}
	return buf[:n], nil
}

// Attach connects this interceptor to the hostapd control interface.
func (i *HostapdInterceptor) Attach() error {
	if i.attached {
		logger.Println("HostapdInterceptor: Already attached.")
		return nil
	}
	res, err := i.Request([]byte(hostapdCmdAttach))
	if err != nil {
		return fmt.Errorf("ATTACH request failed: %w", err)
	}
	if strings.Contains(string(res), hostapdResponseOK) {
		i.attached = true
		logger.Println("HostapdInterceptor: Attached successfully.")
		return nil
	}
	return fmt.Errorf("ATTACH failed, response: %s", string(res))
}

// Deauth sends a deauthentication request for the given MAC address.
func Deauth(macAddress string) error {
	if hostapdInterceptor == nil { // Assumes hostapdInterceptor is a global
		return fmt.Errorf("hostapd interceptor not initialized for DEAUTH")
	}
	cmd := fmt.Sprintf("%s %s", hostapdCmdDeauthenticate, macAddress)
	logger.Printf("HostapdInterceptor: Sending DEAUTH for %s", macAddress)
	_, err := hostapdInterceptor.Request([]byte(cmd))
	if err != nil {
		return fmt.Errorf("DEAUTH for %s failed: %w", macAddress, err)
	}
	logger.Printf("HostapdInterceptor: DEAUTH for %s sent.", macAddress)
	return nil
}

// Shutdown gracefully closes the interceptor's connection.
func (i *HostapdInterceptor) Shutdown() {
	logger.Println("HostapdInterceptor: Shutting down (closing connection)...")
	if i.conn != nil {
		if err := i.conn.Close(); err != nil {
			logger.Printf("HostapdInterceptor: Error closing connection: %v", err)
		} else {
			logger.Println("HostapdInterceptor: Connection closed.")
		}
		i.conn = nil
	}
}

// HostapdListener listens for messages from hostapd.
func HostapdListener(
	allowedMACsFilePath string,
	ueIMSI string,
	dnn string,
	lanIF string,
	pduGatewayIP string,
	leaseTime string,
	quit <-chan struct{},
) {
	defer wg.Done()
	logger.Println("HostapdListener: Started.")

	buf := make([]byte, hostapdSocketBufferSize)
	for {
		select {
		case <-quit:
			logger.Println("HostapdListener: Received quit signal. Stopping...")
			return
		default:
			if hostapdInterceptor == nil || hostapdInterceptor.conn == nil {
				logger.Println("HostapdListener: Interceptor or its connection is nil. Stopping.")
				return
			}

			err := hostapdInterceptor.conn.SetReadDeadline(time.Now().Add(hostapdListenerReadTimeout))
			if err != nil {
				logger.Printf("HostapdListener: SetReadDeadline failed: %v. Retrying.", err)
				time.Sleep(hostapdListenerReadTimeout)
				continue
			}

			n, _, err := hostapdInterceptor.conn.ReadFromUnix(buf)
			_ = hostapdInterceptor.conn.SetReadDeadline(time.Time{})

			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				select {
				case <-quit:
					logger.Println("HostapdListener: Quitting during read error.")
					return
				default:
					logger.Printf("HostapdListener: ReadFromUnix error: %v", err)
					if strings.Contains(err.Error(), "use of closed network connection") {
						logger.Println("HostapdListener: Connection closed externally. Stopping.")
						return
					}
					time.Sleep(1 * time.Second)
				}
				continue
			}

			message := string(buf[:n])
			logPayload := message
			if idx := strings.Index(message, ">"); strings.HasPrefix(message, "<") && idx != -1 && idx+1 < len(message) {
				logPayload = message[idx+1:]
			}
			logger.Printf("HostapdListener: Event: %s", logPayload)

			if strings.Contains(message, hostapdEventEAPSuccess) {
				parts := strings.Fields(message)
				var macAddress string
				for _, part := range parts {
					if hwAddr, parseErr := net.ParseMAC(part); parseErr == nil {
						macAddress = hwAddr.String()
						break
					}
				}

				if macAddress == "" {
					logger.Printf("HostapdListener: No valid MAC in EAP-SUCCESS: %s", message)
					continue
				} else if _, exists := allowedDevices[macAddress]; exists {
					logger.Printf("HostapdListener: Device %s already authenticated.", macAddress)
					continue
				}

				logger.Printf("HostapdListener: Auth success for %s", macAddress)
				logger.Printf("HostapdListener: Requesting PDU for %s (IMSI: %s)", macAddress, ueIMSI)

				session, pduErr := NewPDUSession(ueIMSI, dnn)
				if pduErr != nil {
					logger.Printf("HostapdListener: PDU session for %s failed: %v", macAddress, pduErr)
					logger.Printf("HostapdListener: Deauthenticating %s due to PDU failure.", macAddress)
					if deauthErr := Deauth(macAddress); deauthErr != nil {
						logger.Printf("HostapdListener: DEAUTH for %s failed: %v", macAddress, deauthErr)
					}
					continue
				}

				pduIF := fmt.Sprintf("uesimtun%d", session.ID-1)
				logger.Printf("HostapdListener: PDU Session ID %d, constructed PDU Interface: %s", session.ID, pduIF)

				var appliedRules []AppliedRuleDetail
				var applyErr error
				if ruleManager != nil { // Use the global ruleManager
					logger.Printf("HostapdListener: Applying iptables rules for MAC %s (LAN: %s, PDU_IF: %s, GW: %s)", macAddress, lanIF, pduIF, pduGatewayIP)
					appliedRules, applyErr = ruleManager.ApplyMappingRules(lanIF, macAddress, pduIF, pduGatewayIP, session.ID)
					if applyErr != nil {
						logger.Printf("HostapdListener: Error applying iptables rules for %s: %v. Proceeding without rules.", macAddress, applyErr)
					} else {
						logger.Printf("HostapdListener: Successfully applied %d iptables rules for %s.", len(appliedRules), macAddress)
					}

					allowedDevices[macAddress] = Device{
						state:                "AUTHENTICATED",
						lease:                Lease{},
						pduSession:           session,
						AppliedIPTablesRules: appliedRules,
					}
					logger.Printf("HostapdListener: Device %s tracked, PDU ID %d. Stored %d iptables rules.", macAddress, session.ID, len(appliedRules))

					if err := AllowMAC(allowedMACsFilePath, macAddress, leaseTime); err != nil {
						logger.Printf("HostapdListener: AllowMAC for %s failed: %v", macAddress, err)
					}
					if err := RestartDnsmasq(); err != nil {
						logger.Printf("HostapdListener: RestartDnsmasq for %s failed: %v", macAddress, err)
					}
				} else {
					logger.Printf("HostapdListener: Global ruleManager not initialized, skipping iptables rule application for %s.", macAddress)
				}
			}
		}
	}
}
