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
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/vishvananda/netlink"
)

var (
	logger                     *log.Logger
	allowedDevices             = make(map[string]Device)
	wg                         sync.WaitGroup
	hostapdInterceptor         *HostapdInterceptor
	ruleManager                *RuleManager
	quitHostapdInterceptor     chan struct{}
	quitHostDisconnectListener chan struct{}
	quitDnsmasqListener        chan struct{}
)


func SetLogging(mode string) {
	if mode == "debug" {
		logger = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lshortfile)
	} else {
		w, err := syslog.New(syslog.LOG_INFO|syslog.LOG_DAEMON, "interceptor")
		if err != nil {
			fmt.Printf("Syslog connection failed: %v. Exiting.\n", err)
			os.Exit(1)
		}
		logger = log.New(w, "", 0)
	}
}

func main() {
	logModeFlag := flag.String("mode", "syslog", "Logging mode: 'syslog' or 'debug'")
	hostapdSocketPathFlag := flag.String("interface", "/var/run/hostapd/wlan0", "Hostapd control socket path")
	allowedMACsFileFlag := flag.String("allowed", "/etc/dnsmasq.d/allowed-macs.conf", "Dnsmasq allowed MACs file")
	leasesFileFlag := flag.String("leases", "/var/lib/misc/dnsmasq.leases", "Dnsmasq DHCP leases file")
	ueIMSIFlag := flag.String("imsi", "imsi-999700000000001", "UE IMSI for PDU sessions")
	lanIFFlag := flag.String("lan-if", "enp0s9", "LAN interface name for iptables rules")
	dnnFlag := flag.String("dnn", "clients", "DNN for PDU sessions (default: 'clients')")
	pduGatewayIPFlag := flag.String("pdu-gw-ip", "10.46.0.1", "SMF Session Gateway IP for PDU sessions")
	leasesTimeFlag := flag.String("lease-time", "2m", "Lease time for DHCP leases (e.g., '12h', '2m')")

	flag.Parse()

	logMode := *logModeFlag
	hostapdSocketPath := *hostapdSocketPathFlag
	allowedMACsFilePath := *allowedMACsFileFlag
	leasesFilePath := *leasesFileFlag
	ueIMSI := *ueIMSIFlag
	lanIF := *lanIFFlag
	pduGatewayIP := *pduGatewayIPFlag
	leaseTime := *leasesTimeFlag
	dnn := *dnnFlag

	SetLogging(logMode)
	logger.Println("Interceptor starting...")
	logger.Printf("Config - Mode: %s, Socket: %s, AllowedMACs: %s, Leases: %s, IMSI: %s, LAN_IF: %s, PDU_GW_IP: %s",
		logMode, hostapdSocketPath, allowedMACsFilePath, leasesFilePath, ueIMSI, lanIF, pduGatewayIP)


	var errRuleManager error
	ruleManager, errRuleManager = NewRuleManager() 
	if errRuleManager != nil {
		logger.Fatalf("Failed to initialize RuleManager: %v. Ensure program has necessary privileges (e.g., root for iptables).", errRuleManager)
	}
	logger.Println("RuleManager initialized successfully (global rules applied).")


	linkName := filepath.Base(hostapdSocketPath)
	if hostapdSocketPath == "/var/run/hostapd/global" {
		logger.Printf("Warning: Using global hostapd socket. Derived link '%s' might be invalid for disconnect monitoring.", linkName)
	}

	link, err := netlink.LinkByName(linkName)
	if err != nil {
		logger.Fatalf("LinkByName '%s' (from socket '%s') failed: %v. Check interface name or use dedicated flag.", linkName, hostapdSocketPath, err)
	}
	logger.Printf("Monitoring link %s (Index %d) for disconnects.", link.Attrs().Name, link.Attrs().Index)

	var errInterceptor error
	hostapdInterceptor, errInterceptor = NewInterceptor(hostapdSocketPath)
	if errInterceptor != nil {
		logger.Fatalf("NewInterceptor for %s failed: %v", hostapdSocketPath, errInterceptor)
	}
	defer hostapdInterceptor.Shutdown()
	logger.Println("Hostapd interceptor created.")

	if err := hostapdInterceptor.Attach(); err != nil {
		logger.Fatalf("Attach to hostapd %s failed: %v", hostapdSocketPath, err)
	}
	logger.Printf("Attached to hostapd: %s", hostapdSocketPath)

	if _, err := hostapdInterceptor.Request([]byte("PING")); err != nil {
		logger.Printf("Warning: Initial PING to hostapd failed: %v", err)
	} else {
		logger.Println("Hostapd PING successful.")
	}

	quitHostapdInterceptor = make(chan struct{})
	quitDnsmasqListener = make(chan struct{})
	quitHostDisconnectListener = make(chan struct{})

	wg.Add(1)

	go HostapdListener(allowedMACsFilePath, ueIMSI, dnn, lanIF, pduGatewayIP, leaseTime, quitHostapdInterceptor) 
	
	wg.Add(1)
	go DnsmasqListener(allowedMACsFilePath, leasesFilePath, ueIMSI, leaseTime, quitDnsmasqListener)

	wg.Add(1)
	go HostDisconnectListener(allowedMACsFilePath, leasesFilePath, ueIMSI, link, quitHostDisconnectListener)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	logger.Println("Application running. Ctrl+C to exit.")

	receivedSignal := <-sigChan
	logger.Printf("Signal %s received. Shutting down...", receivedSignal)

	logger.Println("Signaling HostapdListener to stop...")
	close(quitHostapdInterceptor)

	logger.Println("Signaling DnsmasqListener to stop...")
	close(quitDnsmasqListener)

	logger.Println("Signaling HostDisconnectListener to stop...")
	close(quitHostDisconnectListener)

	waitTimeout := 5 * time.Second
	waitDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitDone)
	}()

	select {
	case <-waitDone:
		logger.Println("All signaled listeners have shut down gracefully.")
	case <-time.After(waitTimeout):
		logger.Println("Warning: Timeout waiting for listeners to shut down. Forcing exit.")
	}
	logger.Println("Shutdown process complete. Exiting.")
}
