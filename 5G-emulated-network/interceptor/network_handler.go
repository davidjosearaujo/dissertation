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
	"time"

	"github.com/vishvananda/netlink"
)

const (
	hostDisconnectCheckInterval = 30 * time.Second
)

// Device struct represents a connected client device.
// It now includes a field to store applied iptables rules.
type Device struct {
	state                string
	lease                Lease    // Lease type is defined in dnsmasq_handler.go
	pduSession           *Session // Session type is defined in pdu_handling.go
	AppliedIPTablesRules []AppliedRuleDetail // Stores rules applied by routing_handler.go
}

// ForgetDevice performs a best-effort cleanup for a device.
// It now uses stored iptables rules for removal.
func ForgetDevice(allowedMACsFilePath string, leasesFilePath string, macAddress string, ueIMSI string) {
	logger.Printf("ForgetDevice: Initiating for MAC %s (IMSI: %s)", macAddress, ueIMSI)

	// Retrieve device details, including applied iptables rules
	device, exists := allowedDevices[macAddress]
	if !exists {
		logger.Printf("ForgetDevice: Device %s not found in tracking. Cannot perform full cleanup.", macAddress)
		// Attempt MAC disallow and Deauth as a best effort if device is not tracked but MAC is known
		if err := DisallowMAC(allowedMACsFilePath, macAddress); err != nil {
			logger.Printf("ForgetDevice: Error disallowing MAC %s from %s (device not tracked): %v", macAddress, allowedMACsFilePath, err)
		}
		if err := DisallowMAC(leasesFilePath, macAddress); err != nil {
			logger.Printf("ForgetDevice: Error removing lease for MAC %s from %s (device not tracked): %v", macAddress, leasesFilePath, err)
		}
		if err := Deauth(macAddress); err != nil {
			logger.Printf("ForgetDevice: Deauth %s failed (device not tracked): %v", macAddress, err)
		}
		return
	}

	// Standard cleanup steps
	if err := DisallowMAC(allowedMACsFilePath, macAddress); err != nil {
		logger.Printf("ForgetDevice: Error disallowing MAC %s from %s: %v", macAddress, allowedMACsFilePath, err)
	}
	if err := DisallowMAC(leasesFilePath, macAddress); err != nil {
		logger.Printf("ForgetDevice: Error removing lease for MAC %s from %s: %v", macAddress, leasesFilePath, err)
	}
	if err := RestartDnsmasq(); err != nil {
		logger.Printf("ForgetDevice: Error restarting dnsmasq for %s: %v", macAddress, err)
	}

	logger.Printf("ForgetDevice: Deauthenticating %s via hostapd", macAddress)
	if err := Deauth(macAddress); err != nil {
		logger.Printf("ForgetDevice: Deauth %s failed: %v", macAddress, err)
	}

	// Remove stored iptables rules
	if ruleManager != nil && len(device.AppliedIPTablesRules) > 0 { // ruleManager is assumed global
		if err := ruleManager.RemoveRulesForDevice(macAddress, device.AppliedIPTablesRules); err != nil {
			// RemoveRulesForDevice already logs details, this is a summary log
			logger.Printf("ForgetDevice: Issues encountered removing iptables rules for MAC %s: %v", macAddress, err)
		}
	} else if len(device.AppliedIPTablesRules) == 0 {
		logger.Printf("ForgetDevice: No stored iptables rules to remove for MAC %s.", macAddress)
	} else {
		logger.Printf("ForgetDevice: ruleManager not initialized, cannot remove iptables rules for MAC %s.", macAddress)
	}


	if device.pduSession != nil {
		logger.Printf("ForgetDevice: Releasing PDU ID %d for %s", device.pduSession.ID, macAddress)
		if err := ReleasePDUSession(ueIMSI, device.pduSession.ID); err != nil {
			logger.Printf("ForgetDevice: Release PDU ID %d for %s failed: %v", device.pduSession.ID, macAddress, err)
		}
	} else {
		logger.Printf("ForgetDevice: No PDU session for %s to release.", macAddress)
	}

	logger.Printf("ForgetDevice: Removing %s from internal tracking map.", macAddress)
	delete(allowedDevices, macAddress)
	logger.Printf("ForgetDevice: Completed for MAC %s", macAddress)
}

func HostDisconnectListener(allowedMACsFilePath string, leasesFilePath string, ueIMSI string, link netlink.Link, quit <-chan struct{}) {
	defer wg.Done()
	logger.Printf("HostDisconnectListener: Monitoring link %s (Index %d) every %s", link.Attrs().Name, link.Attrs().Index, hostDisconnectCheckInterval)

	ticker := time.NewTicker(hostDisconnectCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			logger.Println("HostDisconnectListener: Stopping...")
			return
		case <-ticker.C:
			neighs, err := netlink.NeighList(link.Attrs().Index, netlink.FAMILY_V4)
			if err != nil {
				logger.Printf("HostDisconnectListener: NeighList for %s failed: %v.", link.Attrs().Name, err)
				continue 
			}

			activeMACsOnLink := make(map[string]bool)
			for _, neigh := range neighs {
				if neigh.HardwareAddr == nil {
					continue
				}
				macAddress := neigh.HardwareAddr.String()
				activeMACsOnLink[macAddress] = true

				device, deviceExists := allowedDevices[macAddress]
				if !deviceExists {
					continue 
				}

				isReachable := (neigh.State & netlink.NUD_REACHABLE) != 0
				isStale := (neigh.State & netlink.NUD_STALE) != 0
				isFailed := (neigh.State & netlink.NUD_FAILED) != 0
			
				if isReachable {
					if device.state != "REACHABLE" {
						pduAddr := "N/A"
						if device.pduSession != nil { pduAddr = device.pduSession.Address }
						logger.Printf("HostDisconnectListener: Device %s (MAC: %s) state -> REACHABLE (was %s)", pduAddr, macAddress, device.state)
						device.state = "REACHABLE"
						allowedDevices[macAddress] = device
					}
				} else if (isStale || isFailed) && time.Until(time.Unix(int64(device.lease.expiration), 0)) < (device.lease.duration * time.Second * 3 / 4) {
					if device.state == "REACHABLE" || device.state == "LEASED" {
						pduAddr := "N/A"
						if device.pduSession != nil { pduAddr = device.pduSession.Address }
						logger.Printf("HostDisconnectListener: Device %s (MAC: %s) -> STALE (was %s) -> Forgetting.", pduAddr, macAddress, device.state)
						ForgetDevice(allowedMACsFilePath, leasesFilePath, macAddress, ueIMSI)
						continue 
					}
				}
			} 

			macsToForget := []string{}
			for trackedMAC, device := range allowedDevices {
				if _, stillTracked := allowedDevices[trackedMAC]; !stillTracked {
					continue
				}
				if !activeMACsOnLink[trackedMAC] {
					if device.state == "REACHABLE" || device.state == "LEASED" || device.state == "STALE" {
						pduAddr := "N/A"
						if device.pduSession != nil { pduAddr = device.pduSession.Address }
						logger.Printf("HostDisconnectListener: Tracked device %s (MAC: %s, State: %s) no longer in ARP list. Scheduling for forget.", pduAddr, trackedMAC, device.state)
						macsToForget = append(macsToForget, trackedMAC)
					}
				}
			}
			for _, macToForget := range macsToForget {
                if _, ok := allowedDevices[macToForget]; ok {
				    ForgetDevice(allowedMACsFilePath, leasesFilePath, macToForget, ueIMSI)
                }
			}
		}
	}
}
