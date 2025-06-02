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
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Session represents a single PDU session.
type Session struct {
	ID          int    `yaml:"id"`
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

const (
	pduSessionEstablishRetries  = 20
	pduSessionEstablishInterval = 3 * time.Second
	pduSessionCmdEstablish      = "ps-establish IPv4 --sst 1"
	pduSessionCmdList           = "ps-list"
	pduSessionCmdRelease        = "ps-release"
	pduSessionStateActive       = "PS-ACTIVE"
)

// NewPDUSession establishes a new PDU session using nr-cli and waits for it to become active.
func NewPDUSession(ueIMSI string, dnn string) (*Session, error) {
	logger.Printf("IMSI %s establishing...", ueIMSI)
	args := fmt.Sprintf(pduSessionCmdEstablish+" --dnn %s", dnn)
	cmd := exec.Command("nr-cli", ueIMSI, "--exec", args)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("PDU establish for %s failed: %w. Output: %s", ueIMSI, err, string(output))
	}
	logger.Printf("IMSI %s requested. Output: %s. Waiting activation...", ueIMSI, strings.TrimSpace(string(output)))

	var session *Session
	for i := 0; i < pduSessionEstablishRetries; i++ {
		time.Sleep(pduSessionEstablishInterval)
		session, err = LastPDUSession(ueIMSI)
		if err != nil {
			logger.Printf("Retry %d/%d IMSI %s: get status failed: %v", i+1, pduSessionEstablishRetries, ueIMSI, err)
			continue
		}

		if session != nil && session.State == pduSessionStateActive && session.Address != "" {
			logger.Printf("PDU ID %d IMSI %s ACTIVE (State: %s, Addr: %s).", session.ID, ueIMSI, session.State, session.Address)
			return session, nil
		}

		sID, sState, sAddr := -1, "N/A", "N/A"
		if session != nil {
			sID, sState, sAddr = session.ID, session.State, session.Address
		}
		logger.Printf("Retry %d/%d IMSI %s: waiting (ID: %d, State: %s, Addr: %s)...", i+1, pduSessionEstablishRetries, ueIMSI, sID, sState, sAddr)
	}

	if session != nil {
		logger.Printf("PDU ID %d IMSI %s not active after %d retries. Releasing.", session.ID, ueIMSI, pduSessionEstablishRetries)
		if releaseErr := ReleasePDUSession(ueIMSI, session.ID); releaseErr != nil {
			logger.Printf("Release PDU ID %d IMSI %s failed: %v", session.ID, ueIMSI, releaseErr)
		}
	} else {
		logger.Printf("No PDU for IMSI %s or not active after %d retries.", ueIMSI, pduSessionEstablishRetries)
	}
	return nil, fmt.Errorf("PDU for IMSI %s not active after %d retries", ueIMSI, pduSessionEstablishRetries)
}

// LastPDUSession retrieves the most recent PDU session from nr-cli ps-list output.
func LastPDUSession(ueIMSI string) (*Session, error) {
	cmd := exec.Command("nr-cli", ueIMSI, "--exec", pduSessionCmdList)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("PDU list for %s failed: %w. Output: %s", ueIMSI, err, out.String())
	}

	var pduSessionsYAML map[string]yaml.Node
	if err := yaml.Unmarshal(out.Bytes(), &pduSessionsYAML); err != nil {
		if strings.TrimSpace(out.String()) == "" {
			return nil, fmt.Errorf("no PDU sessions for %s (empty nr-cli output)", ueIMSI)
		}
		return nil, fmt.Errorf("parse PDU list YAML for %s failed: %w. Content: %s", ueIMSI, err, out.String())
	}

	if len(pduSessionsYAML) == 0 {
		return nil, fmt.Errorf("no PDU sessions for %s (empty map from nr-cli)", ueIMSI)
	}

	var sessionKeys []int
	sessionMap := make(map[int]Session)

	for key, node := range pduSessionsYAML {
		sessionNumStr := strings.TrimPrefix(key, "PDU Session")
		sessionNum, err := strconv.Atoi(strings.TrimSpace(sessionNumStr))
		if err != nil {
			logger.Printf("LastPDUSession: Invalid session ID '%s' (key '%s') IMSI %s. Skipping.", sessionNumStr, key, ueIMSI)
			continue
		}

		var currentSession Session
		if err := node.Decode(&currentSession); err != nil {
			logger.Printf("LastPDUSession: Decode session key '%s' (ID %d) IMSI %s failed: %v. Skipping.", key, sessionNum, ueIMSI, err)
			continue
		}
		currentSession.ID = sessionNum
		sessionMap[sessionNum] = currentSession
		sessionKeys = append(sessionKeys, sessionNum)
	}

	if len(sessionKeys) == 0 {
		return nil, fmt.Errorf("no valid PDU sessions parsed for %s from nr-cli YAML", ueIMSI)
	}

	sort.Ints(sessionKeys)
	lastSession := sessionMap[sessionKeys[len(sessionKeys)-1]]
	return &lastSession, nil
}

// ReleasePDUSession instructs nr-cli to release a specific PDU session.
func ReleasePDUSession(ueIMSI string, pduID int) error {
	logger.Printf("ReleasePDUSession: PDU ID %d IMSI %s releasing...", pduID, ueIMSI)
	cmd := exec.Command("nr-cli", ueIMSI, "--exec", fmt.Sprintf("%s %d", pduSessionCmdRelease, pduID))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("release PDU #%d IMSI %s failed: %w. Output: %s", pduID, ueIMSI, err, string(output))
	}
	logger.Printf("ReleasePDUSession: PDU ID %d IMSI %s released. Output: %s", pduID, ueIMSI, strings.TrimSpace(string(output)))
	return nil
}
