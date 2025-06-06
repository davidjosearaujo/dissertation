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
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync" // Added for file writing lock

	"github.com/coreos/go-iptables/iptables"
)

// Global static rule specifications for iptables.
var (
	globalRuleRelatedEstablishedSpec = []string{"-m", "state", "--state", "RELATED,ESTABLISHED", "-j", "ACCEPT"}
	rtTablesPath                     = "/etc/iproute2/rt_tables"
	rtTablesMutex                    sync.Mutex // Mutex to protect concurrent access to rt_tables
)

// RuleType defines the type of rule managed.
type RuleType string

const (
	RuleTypeIPTables     RuleType = "iptables"
	RuleTypeIPRoute      RuleType = "ip_route"
	RuleTypeIPRule       RuleType = "ip_rule"
	RuleTypeRTTableEntry RuleType = "rt_table_entry"
)

// AppliedRuleDetail holds the specifics of a rule that was applied.
type AppliedRuleDetail struct {
	Type     RuleType
	Table    string   // For iptables (e.g., "mangle", "nat", "filter")
	Chain    string   // For iptables (e.g., "PREROUTING", "POSTROUTING", "FORWARD")
	RuleSpec []string // For iptables: rule parts; For ip_route/ip_rule: args after "add"; For rt_table: {tableIDStr, tableNameStr}
	Comment  string   // Optional: A comment for easier identification, e.g., MAC address or PDU session ID
}

// RuleManager holds the iptables client and handles rule application/removal.
type RuleManager struct {
	ipt *iptables.IPTables
}

// NewRuleManager initializes a new iptables rule manager and applies global firewall rules.
func NewRuleManager() (*RuleManager, error) {
	ipt, err := iptables.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize iptables: %w", err)
	}
	rm := &RuleManager{ipt: ipt}

	logger.Println("Initializing global iptables rules...")

	// Set FORWARD policy to DROP
	if err := rm.setForwardPolicy("DROP"); err != nil {
		logger.Printf("[CRITICAL_ERROR] Failed to set FORWARD Policy to DROP during init: %v", err)
		return nil, fmt.Errorf("failed to set initial FORWARD policy: %w", err)
	}
	logger.Println("[SUCCESS] Global: FORWARD chain policy set to DROP.")

	// Allow RELATED,ESTABLISHED traffic in FORWARD chain
	// Adding a generic comment for global rules
	if err := rm.ensureRule("filter", "FORWARD", globalRuleRelatedEstablishedSpec, "global_related_established_interceptor"); err != nil {
		logger.Printf("[CRITICAL_ERROR] Failed to ensure FORWARD RELATED,ESTABLISHED rule during init: %v", err)
		return nil, fmt.Errorf("failed to set initial FORWARD RELATED,ESTABLISHED rule: %w", err)
	}
	logger.Printf("[SUCCESS] Global: FORWARD RELATED,ESTABLISHED rule ensured: %v", globalRuleRelatedEstablishedSpec)

	logger.Println("iptables manager and global rules initialized successfully.")
	return rm, nil
}

// executeCommand runs an external command (like "ip")
func (rm *RuleManager) executeCommand(logPrefix string, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	fullCmdStr := command + " " + strings.Join(args, " ")

	if err != nil {
		// More robust check for "already exists" or "not found" errors
		outStr := string(output)
		isAddCmd := strings.Contains(strings.Join(args, " "), "add")
		isDelCmd := strings.Contains(strings.Join(args, " "), "del") || strings.Contains(strings.Join(args, " "), "delete")

		// "RTNETLINK answers: File exists" for 'ip rule add' or 'ip route add' if already present
		if isAddCmd && (strings.Contains(outStr, "File exists") || strings.Contains(outStr, "Object already exists")) {
			logger.Printf("%s: Command '%s' indicated rule/object likely already exists: %s", logPrefix, fullCmdStr, strings.TrimSpace(outStr))
			return nil // Not an error if rule already exists for an 'add' operation
		}
		// "RTNETLINK answers: No such file or directory" for 'ip rule del' or 'ip route del' if not found
		// "Cannot find device" or "No such process" can also happen for deletion if already gone
		if isDelCmd && (strings.Contains(outStr, "No such file or directory") ||
			strings.Contains(outStr, "Cannot find device") ||
			strings.Contains(outStr, "No such process") ||
			strings.Contains(outStr, "does not exist")) {
			logger.Printf("%s: Command '%s' indicated rule/object likely already deleted/not found: %s", logPrefix, fullCmdStr, strings.TrimSpace(outStr))
			return nil // Not an error if rule already gone for a 'delete' operation
		}

		logger.Printf("%s: Command '%s' failed. Output: %s, Error: %v", logPrefix, fullCmdStr, strings.TrimSpace(outStr), err)
		return fmt.Errorf("executing %s: %w. Output: %s", fullCmdStr, err, outStr)
	}
	logger.Printf("%s: Command '%s' executed successfully. Output: %s", logPrefix, fullCmdStr, strings.TrimSpace(strings.ReplaceAll(string(output), "\n", " ")))
	return nil
}

// manageRTTableEntry adds or removes an entry from /etc/iproute2/rt_tables
func (rm *RuleManager) manageRTTableEntry(tableID int, tableName string, add bool, macAddr string) error {
	rtTablesMutex.Lock()
	defer rtTablesMutex.Unlock()

	logPrefix := fmt.Sprintf("RTTable (%s)", macAddr)
	fileContent, err := os.ReadFile(rtTablesPath)
	if err != nil && !os.IsNotExist(err) { // Allow not exist for initial creation
		logger.Printf("%s: Error reading %s: %v", logPrefix, rtTablesPath, err)
		return fmt.Errorf("%s: reading %s: %w", logPrefix, rtTablesPath, err)
	}

	entryLine := fmt.Sprintf("%d\t%s", tableID, tableName) // Use tab as separator, common in rt_tables
	var newLines []string
	found := false
	modified := false

	scanner := bufio.NewScanner(strings.NewReader(string(fileContent)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			newLines = append(newLines, scanner.Text()) // Preserve original line including leading/trailing spaces for comments
			continue
		}

		// Normalize potential multiple spaces/tabs in existing lines for comparison
		parts := strings.Fields(line)
		normalizedExistingLine := ""
		if len(parts) >= 2 {
			normalizedExistingLine = fmt.Sprintf("%s\t%s", parts[0], parts[1])
		}

		if normalizedExistingLine == entryLine {
			found = true
			if !add { // If removing, mark modified and skip this line
				logger.Printf("%s: Removing line '%s' from %s", logPrefix, scanner.Text(), rtTablesPath)
				modified = true
				continue
			}
		}
		newLines = append(newLines, scanner.Text()) // Preserve original line
	}
	if err := scanner.Err(); err != nil {
		logger.Printf("%s: Error scanning %s content: %v", logPrefix, rtTablesPath, err)
		return fmt.Errorf("%s: scanning %s content: %w", logPrefix, rtTablesPath, err)
	}

	if add {
		if found {
			logger.Printf("%s: Entry '%s' already exists in %s.", logPrefix, entryLine, rtTablesPath)
			return nil // Already exists, no change needed
		}
		logger.Printf("%s: Adding line '%s' to %s", logPrefix, entryLine, rtTablesPath)
		newLines = append(newLines, entryLine)
		modified = true
	} else { // If removing
		if !found {
			logger.Printf("%s: Entry '%s' not found in %s for removal.", logPrefix, entryLine, rtTablesPath)
			return nil // Not found, no change needed
		}
		// 'modified' is already true if found and !add
	}

	if modified {
		// Write back the modified content
		// Ensure the file ends with a newline if it's not empty.
		finalContentBuilder := strings.Builder{}
		for i, line := range newLines {
			finalContentBuilder.WriteString(line)
			// Add newline for all but potentially the last line if it's an empty string from split
			if i < len(newLines)-1 || (i == len(newLines)-1 && line != "") {
				finalContentBuilder.WriteString("\n")
			}
		}

		contentToWrite := finalContentBuilder.String()
		// If the file was empty and we added one line, it might not have a trailing newline yet.
		if contentToWrite != "" && !strings.HasSuffix(contentToWrite, "\n") {
			contentToWrite += "\n"
		}

		if err := os.WriteFile(rtTablesPath, []byte(contentToWrite), 0644); err != nil {
			logger.Printf("%s: Error writing updated %s: %v", logPrefix, rtTablesPath, err)
			return fmt.Errorf("%s: writing updated %s: %w", logPrefix, rtTablesPath, err)
		}
		logger.Printf("%s: Successfully updated %s.", logPrefix, rtTablesPath)
	} else {
		logger.Printf("%s: No changes made to %s.", logPrefix, rtTablesPath)
	}
	return nil
}

// ApplyMappingRules configures MAC-specific iptables, ip route, and ip rule entries.
// pduSessionID is used for fwmark and routing table naming.
func (rm *RuleManager) ApplyMappingRules(lanIF, macAddr, pduIF, pduGatewayIP string, pduSessionID int) ([]AppliedRuleDetail, error) {
	var appliedRuleDetails []AppliedRuleDetail
	var errorsEncountered []string
	// Sanitize MAC for use in comments, replace ':' with '_'
	safeMacForComment := strings.ReplaceAll(macAddr, ":", "_")
	comment := fmt.Sprintf("interceptor_mac_%s_pduid_%d", safeMacForComment, pduSessionID)

	logger.Printf("Applying rules for MAC %s (LAN: %s, PDU_IF: %s, GW: %s, PDU_ID: %d, Comment: %s)", macAddr, lanIF, pduIF, pduGatewayIP, pduSessionID, comment)

	pduSessionIDStr := strconv.Itoa(pduSessionID)
	routingTableID := 200 + pduSessionID 
	routingTableName := fmt.Sprintf("table_pdu_%d", pduSessionID) 

	// Order based on `configure_router_logic` from the bash script:
	// 1. /etc/iproute2/rt_tables entry
	// 2. ip route add ... table ...
	// 3. ip rule add fwmark ... table ...
	// 4. iptables ... -j MARK ...
	// 5. iptables ... FORWARD ... -j ACCEPT
	// 6. iptables ... NAT ... -j MASQUERADE

	// 1. Create Custom Routing Table Entry in /etc/iproute2/rt_tables
	if err := rm.manageRTTableEntry(routingTableID, routingTableName, true, macAddr); err != nil {
		errMsg := fmt.Sprintf("Manage RT Table Entry (%s %s) for %s: %v", strconv.Itoa(routingTableID), routingTableName, macAddr, err)
		errorsEncountered = append(errorsEncountered, errMsg)
		logger.Printf("[ERROR] %s", errMsg)
	} else {
		logger.Printf("[SUCCESS] RT Table Entry %d %s for %s ensured/added.", routingTableID, routingTableName, macAddr)
		appliedRuleDetails = append(appliedRuleDetails, AppliedRuleDetail{Type: RuleTypeRTTableEntry, RuleSpec: []string{strconv.Itoa(routingTableID), routingTableName}, Comment: comment})
	}

	// 2. Add route to custom table
	ipRouteArgs := []string{"default", "via", pduGatewayIP, "dev", pduIF, "table", routingTableName}
	cmdArgsRoute := append([]string{"route", "add"}, ipRouteArgs...)
	if err := rm.executeCommand(fmt.Sprintf("IPRouteAdd (%s)", macAddr), "ip", cmdArgsRoute...); err != nil {
		errMsg := fmt.Sprintf("IP Route Add for %s (table %s): %v", macAddr, routingTableName, err)
		errorsEncountered = append(errorsEncountered, errMsg)
		logger.Printf("[ERROR] %s", errMsg)
	} else {
		logger.Printf("[SUCCESS] IP Route Add for %s (table %s): default via %s dev %s", macAddr, routingTableName, pduGatewayIP, pduIF)
		appliedRuleDetails = append(appliedRuleDetails, AppliedRuleDetail{Type: RuleTypeIPRoute, RuleSpec: ipRouteArgs, Comment: comment})
	}

	// 3. Create rule to use custom routing table
	ipRuleArgs := []string{"fwmark", pduSessionIDStr, "table", routingTableName}
	cmdArgsRule := append([]string{"rule", "add"}, ipRuleArgs...)
	if err := rm.executeCommand(fmt.Sprintf("IPRuleAdd (%s)", macAddr), "ip", cmdArgsRule...); err != nil {
		errMsg := fmt.Sprintf("IP Rule Add for %s (fwmark %s, table %s): %v", macAddr, pduSessionIDStr, routingTableName, err)
		errorsEncountered = append(errorsEncountered, errMsg)
		logger.Printf("[ERROR] %s", errMsg)
	} else {
		logger.Printf("[SUCCESS] IP Rule Add for %s: fwmark %s table %s", macAddr, pduSessionIDStr, routingTableName)
		appliedRuleDetails = append(appliedRuleDetails, AppliedRuleDetail{Type: RuleTypeIPRule, RuleSpec: ipRuleArgs, Comment: comment})
	}

	// 4. Mangle PREROUTING rule to MARK packets
	mangleRuleSpec := []string{"-i", lanIF, "-m", "mac", "--mac-source", macAddr, "-j", "MARK", "--set-mark", pduSessionIDStr}
	if err := rm.ensureRule("mangle", "PREROUTING", mangleRuleSpec, comment); err != nil {
		errMsg := fmt.Sprintf("Mangle MARK rule for %s: %v", macAddr, err)
		errorsEncountered = append(errorsEncountered, errMsg)
		logger.Printf("[ERROR] %s", errMsg)
	} else {
		logger.Printf("[SUCCESS] Mangle PREROUTING MARK for %s: %v", macAddr, mangleRuleSpec)
		appliedRuleDetails = append(appliedRuleDetails, AppliedRuleDetail{Type: RuleTypeIPTables, Table: "mangle", Chain: "PREROUTING", RuleSpec: mangleRuleSpec, Comment: comment})
	}

	// 5. Allow forwarding from LAN to PDU IF based on MAC AND MARK
	forwardMacRuleSpec := []string{
		"-i", lanIF,
		"-o", pduIF,
		"-m", "mac", "--mac-source", macAddr,
		"-m", "mark", "--mark", pduSessionIDStr, 
		"-j", "ACCEPT",
	}
	if err := rm.ensureRule("filter", "FORWARD", forwardMacRuleSpec, comment); err != nil {
		errMsg := fmt.Sprintf("FORWARD allow MAC %s with mark %s to PDU %s: %v", macAddr, pduSessionIDStr, pduIF, err)
		errorsEncountered = append(errorsEncountered, errMsg)
		logger.Printf("[ERROR] %s", errMsg)
	} else {
		logger.Printf("[SUCCESS] FORWARD allow MAC %s from %s with mark %s to %s: %v", macAddr, lanIF, pduSessionIDStr, pduIF, forwardMacRuleSpec)
		appliedRuleDetails = append(appliedRuleDetails, AppliedRuleDetail{Type: RuleTypeIPTables, Table: "filter", Chain: "FORWARD", RuleSpec: forwardMacRuleSpec, Comment: comment})
	}

	// 6. NAT MASQUERADE for PDU interface
	natMasqueradeRuleSpec := []string{"-o", pduIF, "-j", "MASQUERADE"}
	if err := rm.ensureRule("nat", "POSTROUTING", natMasqueradeRuleSpec, comment); err != nil {
		errMsg := fmt.Sprintf("NAT MASQUERADE for PDU %s: %v", pduIF, err)
		errorsEncountered = append(errorsEncountered, errMsg)
		logger.Printf("[ERROR] %s", errMsg)
	} else {
		logger.Printf("[SUCCESS] NAT POSTROUTING MASQUERADE for %s: %v", pduIF, natMasqueradeRuleSpec)
		appliedRuleDetails = append(appliedRuleDetails, AppliedRuleDetail{Type: RuleTypeIPTables, Table: "nat", Chain: "POSTROUTING", RuleSpec: natMasqueradeRuleSpec, Comment: comment})
	}


	if len(errorsEncountered) > 0 {
		logger.Printf("[SUMMARY_ERRORS] Encountered %d error(s) applying rules for MAC %s (PDU_ID: %d):", len(errorsEncountered), macAddr, pduSessionID)
		for i, e := range errorsEncountered {
			logger.Printf("  %d: %s", i+1, e)
		}
		return appliedRuleDetails, fmt.Errorf("encountered %d error(s) applying rules for MAC %s. See logs", len(errorsEncountered), macAddr)
	}
	return appliedRuleDetails, nil
}

// RemoveRulesForDevice removes the specified rules for a device.
// Rules are removed in reverse order of application for dependency management.
func (rm *RuleManager) RemoveRulesForDevice(macAddr string, rulesToRemove []AppliedRuleDetail) error {
	logger.Printf("Removing %d stored rules for MAC %s", len(rulesToRemove), macAddr)
	var errorsEncountered []string
	var successfullyRemovedCount int

	// Iterate in reverse order to handle dependencies
	for i := len(rulesToRemove) - 1; i >= 0; i-- {
		rule := rulesToRemove[i]
		logPrefix := fmt.Sprintf("RemoveRule (MAC: %s, Type: %s, Comment: %s)", macAddr, rule.Type, rule.Comment)
		var err error

		switch rule.Type {
		case RuleTypeIPTables:
			
			ruleSpecForDelete := rule.RuleSpec
			if rule.Comment != "" {
				
				hasComment := false
				for k := 0; k < len(ruleSpecForDelete)-2; k++ {
					if ruleSpecForDelete[k] == "-m" && ruleSpecForDelete[k+1] == "comment" && ruleSpecForDelete[k+2] == "--comment" {
						hasComment = true
						break
					}
				}
				if !hasComment { 
					
					ruleSpecForDelete = append(make([]string, 0, len(rule.RuleSpec)+3), rule.RuleSpec...) 
					ruleSpecForDelete = append(ruleSpecForDelete, "-m", "comment", "--comment", rule.Comment)
				}
			}

			err = rm.ipt.Delete(rule.Table, rule.Chain, ruleSpecForDelete...)
			if err != nil {
				
				errMsgStr := err.Error()
				if strings.Contains(errMsgStr, "No chain/target/match by that name") || 
					strings.Contains(errMsgStr, "does not exist") || 
					strings.Contains(errMsgStr, "No such file or directory") || 
					strings.Contains(strings.ToLower(errMsgStr), "rule not found") ||
					(strings.Contains(errMsgStr, "Bad rule") && (strings.Contains(rule.Table, "mangle") || strings.Contains(rule.Table, "nat") || strings.Contains(rule.Table, "filter"))) { 
					logger.Printf("%s: IPTables rule (Table: %s, Chain: %s, Spec: %v) likely already removed or not found: %v", logPrefix, rule.Table, rule.Chain, ruleSpecForDelete, err)
					
				} else {
					errMsg := fmt.Sprintf("deleting IPTables rule for MAC %s (Table: %s, Chain: %s, Spec: %v): %v", macAddr, rule.Table, rule.Chain, ruleSpecForDelete, err)
					errorsEncountered = append(errorsEncountered, errMsg)
					logger.Printf("[ERROR] %s: %s", logPrefix, errMsg)
				}
			} else {
				logger.Printf("[SUCCESS] %s: Deleted IPTables rule (Table: %s, Chain: %s, Spec: %v)", logPrefix, rule.Table, rule.Chain, ruleSpecForDelete)
				successfullyRemovedCount++
			}

		case RuleTypeIPRoute:
			delArgs := append([]string{"route", "del"}, rule.RuleSpec...)
			err = rm.executeCommand(logPrefix, "ip", delArgs...)
			if err != nil {
				errMsg := fmt.Sprintf("IP Route Del for MAC %s (Spec: %v): %v", macAddr, rule.RuleSpec, err)
				errorsEncountered = append(errorsEncountered, errMsg)
			} else {
				logger.Printf("[SUCCESS] %s: Deleted IP Route (Spec: %v)", logPrefix, rule.RuleSpec)
				successfullyRemovedCount++
			}

		case RuleTypeIPRule:
			delArgs := append([]string{"rule", "del"}, rule.RuleSpec...)
			err = rm.executeCommand(logPrefix, "ip", delArgs...)
			if err != nil {
				errMsg := fmt.Sprintf("IP Rule Del for MAC %s (Spec: %v): %v", macAddr, rule.RuleSpec, err)
				errorsEncountered = append(errorsEncountered, errMsg)
			} else {
				logger.Printf("[SUCCESS] %s: Deleted IP Rule (Spec: %v)", logPrefix, rule.RuleSpec)
				successfullyRemovedCount++
			}

		case RuleTypeRTTableEntry:
			if len(rule.RuleSpec) == 2 {
				tableID, convErr := strconv.Atoi(rule.RuleSpec[0])
				if convErr != nil {
					errMsg := fmt.Sprintf("Invalid table ID in RuleSpec for RTTableEntry removal for MAC %s: %v", macAddr, rule.RuleSpec)
					errorsEncountered = append(errorsEncountered, errMsg)
					logger.Printf("[ERROR] %s: %s", logPrefix, errMsg)
					continue
				}
				tableName := rule.RuleSpec[1]
				err = rm.manageRTTableEntry(tableID, tableName, false, macAddr) // false for remove
				if err != nil {
					errMsg := fmt.Sprintf("Manage RT Table Entry Del for MAC %s (ID: %d, Name: %s): %v", macAddr, tableID, tableName, err)
					errorsEncountered = append(errorsEncountered, errMsg)
					logger.Printf("[ERROR] %s: %s", logPrefix, errMsg)
				} else {
					logger.Printf("[SUCCESS] %s: Removed/Ensured absent RT Table Entry (ID: %d, Name: %s)", logPrefix, tableID, tableName)
					successfullyRemovedCount++
				}
			} else {
				errMsg := fmt.Sprintf("Invalid RuleSpec for RTTableEntry removal for MAC %s: %v", macAddr, rule.RuleSpec)
				errorsEncountered = append(errorsEncountered, errMsg)
				logger.Printf("[ERROR] %s: %s", logPrefix, errMsg)
			}
		default:
			errMsg := fmt.Sprintf("Unknown rule type '%s' for MAC %s, cannot remove.", rule.Type, macAddr)
			errorsEncountered = append(errorsEncountered, errMsg)
			logger.Printf("[ERROR] %s: %s", logPrefix, errMsg)
		}
	}

	if len(errorsEncountered) > 0 {
		logger.Printf("[SUMMARY_ERRORS] Encountered %d error(s) during rule removal for MAC %s. Successfully removed: %d.", len(errorsEncountered), macAddr, successfullyRemovedCount)
		for i, e := range errorsEncountered {
			logger.Printf("  %d: %s", i+1, e)
		}
	} else if len(rulesToRemove) > 0 {
		logger.Printf("Rule removal process completed for MAC %s. Successfully removed/verified absent: %d of %d.", macAddr, successfullyRemovedCount, len(rulesToRemove))
	} else {
		logger.Printf("No rules specified for removal for MAC %s.", macAddr)
	}

	return nil // Best-effort removal
}

// ensureRule appends an iptables rule if it doesn't already exist.
// It now includes a comment for better identification of rules.
func (rm *RuleManager) ensureRule(table, chain string, ruleSpec []string, comment string) error {
	// Create a copy of ruleSpec to avoid modifying the original slice
	finalRuleSpec := make([]string, len(ruleSpec))
	copy(finalRuleSpec, ruleSpec)

	// Add comment if provided and not already present
	if comment != "" {
		hasComment := false
		for i := 0; i < len(finalRuleSpec)-2; i++ { // Check up to len-2 for "-m comment --comment"
			if finalRuleSpec[i] == "-m" && finalRuleSpec[i+1] == "comment" && finalRuleSpec[i+2] == "--comment" {
				hasComment = true
				break
			}
		}
		if !hasComment {
			finalRuleSpec = append(finalRuleSpec, "-m", "comment", "--comment", comment)
		}
	}

	exists, err := rm.ipt.Exists(table, chain, finalRuleSpec...)
	if err != nil {
		return fmt.Errorf("checking rule existence (table: %s, chain: %s, rule: %v): %w", table, chain, finalRuleSpec, err)
	}
	if !exists {
		if err := rm.ipt.Append(table, chain, finalRuleSpec...); err != nil {
			return fmt.Errorf("appending rule (table: %s, chain: %s, rule: %v): %w", table, chain, finalRuleSpec, err)
		}
		logger.Printf("Appended rule to %s %s: %v", table, chain, finalRuleSpec)
	} else {
		logger.Printf("Rule already exists in %s %s: %v", table, chain, finalRuleSpec)
	}
	return nil
}

// setForwardPolicy sets the default policy for the FORWARD chain in the filter table.
func (rm *RuleManager) setForwardPolicy(policy string) error {
	validPolicies := map[string]bool{"ACCEPT": true, "DROP": true, "REJECT": true}
	upperPolicy := strings.ToUpper(policy)
	if !validPolicies[upperPolicy] {
		return fmt.Errorf("invalid policy '%s'. Must be ACCEPT, DROP, or REJECT", policy)
	}

	if err := rm.ipt.ChangePolicy("filter", "FORWARD", upperPolicy); err != nil {
		return fmt.Errorf("setting FORWARD chain policy to %s: %w", upperPolicy, err)
	}
	logger.Printf("FORWARD chain policy set to %s.", upperPolicy)
	return nil
}

// main function for demonstration or testing purposes
// func main() {
// 	logger.Println("Starting PBR Demo...")

// 	rm, err := NewRuleManager()
// 	if err != nil {
// 		logger.Fatalf("Failed to create %v", err)
// 	}

// 	// Example usage:
// 	lanIF := "enp0s9" // Replace with your actual LAN interface
// 	macAddr := "08:00:27:01:02:03" // Replace with a test MAC
// 	pduIF := "uesimtun5" // Replace with a test PDU interface (e.g., uesimtunX)
// 	pduGatewayIP := "10.46.0.1" // Replace with your PDU gateway
// 	pduSessionID := 5 // Example PDU Session ID

// 	logger.Printf("Attempting to apply rules for MAC: %s, PDU Session ID: %d", macAddr, pduSessionID)
// 	appliedRules, err := rm.ApplyMappingRules(lanIF, macAddr, pduIF, pduGatewayIP, pduSessionID)
// 	if err != nil {
// 		logger.Printf("Error applying mapping rules: %v", err)
// 		// Optionally, attempt cleanup even if apply failed partially
// 		// rm.RemoveRulesForDevice(macAddr, appliedRules)
// 	} else {
// 		logger.Printf("Successfully applied %d rules/entries for MAC %s.", len(appliedRules), macAddr)
// 	}

// 	// Example: Wait for a key press before cleaning up
// 	// logger.Println("Rules applied. Press Enter to remove rules and exit...")
// 	// bufio.NewReader(os.Stdin).ReadBytes('\n')

// 	// logger.Printf("Attempting to remove rules for MAC: %s", macAddr)
// 	// if err := rm.RemoveRulesForDevice(macAddr, appliedRules); err != nil {
// 	// 	logger.Printf("Error removing rules: %v", err)
// 	// } else {
// 	// 	logger.Println("Successfully removed rules.")
// 	// }

// 	logger.Println("PBR Demo finished.")
// }
