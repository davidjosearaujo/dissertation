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
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	dnsmasqLeaseCheckInterval = 10 * time.Second
)

// Lease holds DHCP lease information for a device.
type Lease struct {
	expiration int
	counter    int
	duration   time.Duration
}

// AllowMAC adds a MAC address to the dnsmasq allowed list.
func AllowMAC(allowedMACsFilePath string, macAddress string, leaseTime string) error {
	f, err := os.OpenFile(allowedMACsFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening %s: %w", allowedMACsFilePath, err)
	}
	defer f.Close()

	entry := fmt.Sprintf("dhcp-host=%s,%s,set:known\n", macAddress,leaseTime)
	if _, err := f.WriteString(entry); err != nil {
		return fmt.Errorf("writing MAC %s to %s: %w", macAddress, allowedMACsFilePath, err)
	}
	logger.Printf("MAC %s added to %s.", macAddress, allowedMACsFilePath)
	return nil
}

// DisallowMAC removes a MAC address from the dnsmasq allowed list.
func DisallowMAC(allowedMACsFilePath string, macAddress string) error {
	fileContent, err := os.ReadFile(allowedMACsFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Printf("%s not found, cannot remove %s.", allowedMACsFilePath, macAddress)
			return nil
		}
		return fmt.Errorf("reading %s: %w", allowedMACsFilePath, err)
	}

	lines := strings.Split(string(fileContent), "\n")
	var newLines []string
	found := false
	macSearchString1 := "dhcp-host=" + macAddress + ","
	macSearchString2 := "dhcp-host=" + macAddress + " "

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if !(strings.HasPrefix(trimmedLine, macSearchString1) ||
			strings.HasPrefix(trimmedLine, macSearchString2) ||
			trimmedLine == "dhcp-host="+macAddress) {
			newLines = append(newLines, line)
		} else {
			logger.Printf("Removing MAC %s line: '%s'", macAddress, trimmedLine)
			found = true
		}
	}

	if !found {
		logger.Printf("MAC %s not found in %s for removal.", macAddress, allowedMACsFilePath)
		return nil
	}

	finalContentBuilder := strings.Builder{}
	for i, line := range newLines {
		if line == "" && (i == len(newLines)-1 || newLines[i+1] == "") {
			continue
		}
		finalContentBuilder.WriteString(line)
		if i < len(newLines)-1 {
			finalContentBuilder.WriteString("\n")
		}
	}
	newContent := strings.TrimSuffix(finalContentBuilder.String(), "\n\n")
	if newContent != "" && !strings.HasSuffix(newContent, "\n") {
		newContent += "\n"
	}

	err = os.WriteFile(allowedMACsFilePath, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("writing updated %s: %w", allowedMACsFilePath, err)
	}
	logger.Printf("MAC %s removed from %s.", macAddress, allowedMACsFilePath)
	return nil
}

// RestartDnsmasq executes the systemctl command to restart the dnsmasq service.
func RestartDnsmasq() error {
	cmd := exec.Command("systemctl", "restart", "dnsmasq")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("restarting dnsmasq: %w. Output: %s", err, string(output))
	}
	logger.Println("dnsmasq restarted.")
	return nil
}

// DnsmasqListener monitors the dnsmasq lease file for changes.
func DnsmasqListener(allowedMACsFilePath string, leasesFilePath string, ueIMSI string, leaseTime string, quit <-chan struct{}) {
	defer wg.Done()
	logger.Printf("Monitoring %s every %s", leasesFilePath, dnsmasqLeaseCheckInterval)

	var lastModTime time.Time
	initialStat, err := os.Stat(leasesFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			logger.Printf("Initial stat for %s failed: %v.", leasesFilePath, err)
			// Consider if this is a fatal error for the listener. For now, it will try to continue.
		}
		logger.Printf("%s not found initially. Will check periodically.", leasesFilePath)
	} else {
		lastModTime = initialStat.ModTime()
	}

	ticker := time.NewTicker(dnsmasqLeaseCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			logger.Println("Received quit signal. Stopping...")
			return
		case <-ticker.C:
			currentStat, statErr := os.Stat(leasesFilePath)
			if statErr != nil {
				if os.IsNotExist(statErr) {
					if lastModTime != (time.Time{}) { // Log only if it previously existed
						logger.Printf("Lease file %s not found.", leasesFilePath)
						lastModTime = time.Time{} // Reset to detect recreation
					}
				} else {
					logger.Printf("Error retrieving lease file stats for %s: %v", leasesFilePath, statErr)
				}
				continue
			}

			if currentStat.ModTime().After(lastModTime) {
				logger.Printf("Lease file %s changed. Processing...", leasesFilePath)
				file, openErr := os.Open(leasesFilePath)
				if openErr != nil {
					logger.Printf("Error opening lease file %s: %v", leasesFilePath, openErr)
					lastModTime = currentStat.ModTime() // Update time to avoid reprocessing error immediately
					continue
				}

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					fields := strings.Fields(scanner.Text())
					if len(fields) < 3 { // Need at least Expiry, MAC, IP
						continue
					}

					expirationStr, macAddressLease, ipAddress := fields[0], fields[1], fields[2]
					if _, parseErr := net.ParseMAC(macAddressLease); parseErr != nil {
						logger.Printf("Invalid MAC address '%s' in lease line: %s", macAddressLease, scanner.Text())
						continue
					}

					expiration, convErr := strconv.Atoi(expirationStr)
					if convErr != nil {
						logger.Printf("Error converting expiration '%s' to int for MAC %s: %v", expirationStr, macAddressLease, convErr)
						continue
					}

					if device, exists := allowedDevices[macAddressLease]; exists {
						if device.lease.expiration != expiration || device.state == "AUTHENTICATED" {
							device.lease.counter++
							device.lease.expiration = expiration
							leaseDuration, err := time.ParseDuration(leaseTime)
							if err != nil {
								logger.Printf("Error converting lease duration '%s' to time.Duration: %v" , leaseTime, err)
								continue
							}
							device.lease.duration = leaseDuration
							pduAddr := "N/A"
							if device.pduSession != nil { // Nil check for pduSession
								pduAddr = device.pduSession.Address
							}

							if device.state == "AUTHENTICATED" {
								device.state = "LEASED"
								logger.Printf("Device %s (MAC: %s, IP: %s) transitioned to LEASED state.", pduAddr, macAddressLease, ipAddress)
							}
							allowedDevices[macAddressLease] = device
							logger.Printf("Lease updated for %s (MAC: %s, IP: %s). Exp: %d, Count: %d. State: %s", pduAddr, macAddressLease, ipAddress, expiration, device.lease.counter, device.state)
						}
					}
				}
				file.Close()
				if scanErr := scanner.Err(); scanErr != nil {
					logger.Printf("Error scanning lease file %s: %v", leasesFilePath, scanErr)
				}
				lastModTime = currentStat.ModTime()
			}
		}
	}
}
