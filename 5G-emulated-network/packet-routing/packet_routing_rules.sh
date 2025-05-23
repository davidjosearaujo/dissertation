#!/bin/bash

# PBR Test Script with Namespaces and Revised Non-Colliding IPs
# Simulates Client, Router (root ns), and WAN Gateway in separate namespaces/root.

# Exit on error
set -e

# === Configuration Variables ===
NS_CLIENT="ns_client_rev"
NS_WAN_GW="ns_wan_gw_rev"

# Veth pair for LAN simulation
V_CLI_LAN="v_cli_lan_nsr"    # Client side (in ns_client_rev)
V_LAN_ROOT="v_lan_root_nsr"  # Router side (root ns, our LAN_IF)

# Veth pair for WAN simulation
V_WAN_ROOT="v_wan_root_nsr"  # Router side (root ns, our PDU_IF)
V_GW_WAN="v_gw_wan_nsr"      # Gateway side (in ns_wan_gw_rev)

CLIENT_MAC="08:00:27:24:96:ec"
CLIENT_IP_CIDR="172.20.0.100/24"
CLIENT_IP=$(echo $CLIENT_IP_CIDR | cut -d'/' -f1)

ROUTER_LAN_IP_CIDR="172.20.0.1/24"
ROUTER_LAN_IP=$(echo $ROUTER_LAN_IP_CIDR | cut -d'/' -f1)

ROUTER_WAN_IP_CIDR="172.21.0.100/24"
ROUTER_WAN_IP=$(echo $ROUTER_WAN_IP_CIDR | cut -d'/' -f1)
PDU_IF_MTU="1400"

PDU_GATEWAY_IP_CIDR="172.21.0.1/24"
PDU_GATEWAY_IP=$(echo $PDU_GATEWAY_IP_CIDR | cut -d'/' -f1)

EXTERNAL_TEST_IP="8.8.8.8" # Dummy IP in ns_wan_gw_rev for testing NAT

# Interfaces for the router logic
LAN_IF="$V_LAN_ROOT"
PDU_IF="$V_WAN_ROOT"

# Policy Based Routing settings
MARK_VALUE="1"
TABLE_ID="205" # New ID for this version
TABLE_NAME="pdu_ns_revised_ips_test"

# Log Prefixes (NSR for Namespace Revised)
LOG_MANGLE_PRE="MANGLE_PRE_NSR: "
LOG_MANGLE_POST="MANGLE_POST_NSR: "
LOG_FORWARD_REL_EST="FORWARD_REL_EST_NSR: "
LOG_FORWARD_SPECIFIC="FORWARD_SPECIFIC_NSR: "
LOG_FORWARD_DROP="FORWARD_DROP_NSR: "
LOG_NAT_PRE_MASQ="NAT_PRE_MASQ_NSR: "

# --- Helper Functions ---
check_root() {
  if [ "$(id -u)" -ne 0 ]; then
    echo "This script must be run as root."
    exit 1
  fi
}

cleanup() {
  echo "--- Cleaning up Namespace Revised Test Environment ---"

  # Delete namespaces (implicitly deletes veth ends within them)
  sudo ip netns del "$NS_CLIENT" 2>/dev/null || true
  sudo ip netns del "$NS_WAN_GW" 2>/dev/null || true

  # Delete veth pairs (if one end was not in a deleted namespace, usually not needed if ns deleted first)
  sudo ip link del "$V_LAN_ROOT" 2>/dev/null || true
  sudo ip link del "$V_WAN_ROOT" 2>/dev/null || true
  
  # Remove iptables rules
  echo "Flushing iptables and resetting FORWARD policy to ACCEPT..."
  sudo iptables -t mangle -F PREROUTING 2>/dev/null || true
  sudo iptables -t nat -F POSTROUTING 2>/dev/null || true
  sudo iptables -F FORWARD 2>/dev/null || true
  sudo iptables -P FORWARD ACCEPT 2>/dev/null || true

  # Delete ip rule
  echo "Deleting ip rule for mark $MARK_VALUE table $TABLE_NAME..."
  while sudo ip rule del fwmark "$MARK_VALUE" table "$TABLE_NAME" 2>/dev/null; do :; done

  # Flush custom route table
  echo "Flushing custom route table $TABLE_NAME..."
  sudo ip route flush table "$TABLE_NAME" 2>/dev/null || true

  # Remove entry from /etc/iproute2/rt_tables
  if grep -q -E "^\s*$TABLE_ID\s+$TABLE_NAME\s*$" /etc/iproute2/rt_tables; then
    echo "Removing $TABLE_ID $TABLE_NAME from /etc/iproute2/rt_tables..."
    sudo sed -i -E "/^\s*$TABLE_ID\s+$TABLE_NAME\s*$/d" /etc/iproute2/rt_tables
  fi
  
  echo "Namespace Revised Cleanup complete."
}

setup_environment() {
  echo "--- Setting up Namespace Revised VETH test environment ---"

  # 1. Create Namespaces
  echo "Creating namespaces: $NS_CLIENT, $NS_WAN_GW"
  sudo ip netns add "$NS_CLIENT"
  sudo ip netns add "$NS_WAN_GW"

  # 2. Setup LAN veth pair (V_CLI_LAN <--> V_LAN_ROOT)
  echo "Creating LAN veth pair: $V_CLI_LAN <--> $V_LAN_ROOT"
  sudo ip link add "$V_CLI_LAN" type veth peer name "$V_LAN_ROOT"
  
  sudo ip link set "$V_CLI_LAN" netns "$NS_CLIENT"
  
  echo "Configuring $V_LAN_ROOT (Router LAN IF) in root namespace..."
  sudo ip addr add "$ROUTER_LAN_IP_CIDR" dev "$V_LAN_ROOT"
  sudo ip link set "$V_LAN_ROOT" up
  
  echo "Configuring $V_CLI_LAN in $NS_CLIENT namespace..."
  sudo ip netns exec "$NS_CLIENT" ip addr add "$CLIENT_IP_CIDR" dev "$V_CLI_LAN"
  sudo ip netns exec "$NS_CLIENT" ip link set dev "$V_CLI_LAN" address "$CLIENT_MAC"
  sudo ip netns exec "$NS_CLIENT" ip link set "$V_CLI_LAN" up
  sudo ip netns exec "$NS_CLIENT" ip route add default via "$ROUTER_LAN_IP"

  # 3. Setup WAN veth pair (V_WAN_ROOT <--> V_GW_WAN)
  echo "Creating WAN veth pair: $V_WAN_ROOT <--> $V_GW_WAN"
  sudo ip link add "$V_WAN_ROOT" type veth peer name "$V_GW_WAN"
  
  sudo ip link set "$V_GW_WAN" netns "$NS_WAN_GW"
  
  echo "Configuring $V_WAN_ROOT (Router WAN IF) in root namespace..."
  sudo ip addr add "$ROUTER_WAN_IP_CIDR" dev "$V_WAN_ROOT"
  sudo ip link set "$V_WAN_ROOT" mtu "$PDU_IF_MTU"
  sudo ip link set "$V_WAN_ROOT" up
  
  echo "Configuring $V_GW_WAN in $NS_WAN_GW namespace..."
  sudo ip netns exec "$NS_WAN_GW" ip addr add "$PDU_GATEWAY_IP_CIDR" dev "$V_GW_WAN"
  sudo ip netns exec "$NS_WAN_GW" ip link set "$V_GW_WAN" up
  # Enable IP forwarding in WAN GW namespace
  sudo ip netns exec "$NS_WAN_GW" sysctl -w net.ipv4.ip_forward=1
  # Add dummy external IP in WAN GW namespace's loopback for NAT testing
  sudo ip netns exec "$NS_WAN_GW" ip addr add "$EXTERNAL_TEST_IP/32" dev lo
  sudo ip netns exec "$NS_WAN_GW" ip link set dev lo up # Ensure lo is up

  echo "Namespace Revised Environment setup complete."
}

configure_router_logic() {
  echo "--- Configuring Router Logic (root namespace) ---"

  # Phase 1: Enable IP Forwarding
  echo "Enabling IP forwarding in root namespace..."
  sudo sysctl -w net.ipv4.ip_forward=1

  # Phase 2: Setup Custom Routing Infrastructure
  echo "Setting up custom routing table $TABLE_NAME (ID: $TABLE_ID)..."
  if ! grep -q -E "^\s*$TABLE_ID\s+$TABLE_NAME\s*$" /etc/iproute2/rt_tables; then
    echo "$TABLE_ID $TABLE_NAME" | sudo tee -a /etc/iproute2/rt_tables
  else
    echo "Entry $TABLE_ID $TABLE_NAME already exists in /etc/iproute2/rt_tables."
  fi
  
  sudo ip route add default via "$PDU_GATEWAY_IP" dev "$PDU_IF" table "$TABLE_NAME"
  sudo ip rule add fwmark "$MARK_VALUE" table "$TABLE_NAME"

  # Phase 3: Packet Marking (Mangle Table)
  echo "Setting up Mangle rules for $LAN_IF..."
  sudo iptables -t mangle -A PREROUTING -i "$LAN_IF" -m mac --mac-source "$CLIENT_MAC" -j LOG --log-prefix "$LOG_MANGLE_PRE" --log-level 4
  sudo iptables -t mangle -A PREROUTING -i "$LAN_IF" -m mac --mac-source "$CLIENT_MAC" -j MARK --set-mark "$MARK_VALUE"
  sudo iptables -t mangle -A PREROUTING -i "$LAN_IF" -m mac --mac-source "$CLIENT_MAC" -m mark --mark "$MARK_VALUE" -j LOG --log-prefix "$LOG_MANGLE_POST" --log-level 4

  # Phase 4: Firewall Forwarding Rules
  echo "Setting up FORWARD rules (LAN_IF: $LAN_IF, PDU_IF: $PDU_IF)..."
  sudo iptables -P FORWARD DROP
  sudo iptables -A FORWARD -m state --state RELATED,ESTABLISHED -j LOG --log-prefix "$LOG_FORWARD_REL_EST" --log-level 4
  sudo iptables -A FORWARD -m state --state RELATED,ESTABLISHED -j ACCEPT
  sudo iptables -A FORWARD -i "$LAN_IF" -o "$PDU_IF" -m mac --mac-source "$CLIENT_MAC" -m mark --mark "$MARK_VALUE" -j LOG --log-prefix "$LOG_FORWARD_SPECIFIC" --log-level 4
  sudo iptables -A FORWARD -i "$LAN_IF" -o "$PDU_IF" -m mac --mac-source "$CLIENT_MAC" -m mark --mark "$MARK_VALUE" -j ACCEPT
  sudo iptables -A FORWARD -j LOG --log-prefix "$LOG_FORWARD_DROP" --log-level 4 # Catch-all

  # Phase 5: NAT Rule
  echo "Setting up NAT rule for $PDU_IF..."
  sudo iptables -t nat -A POSTROUTING -o "$PDU_IF" -j LOG --log-prefix "$LOG_NAT_PRE_MASQ" --log-level 4
  sudo iptables -t nat -A POSTROUTING -o "$PDU_IF" -j MASQUERADE
  
  echo "Router Logic configuration complete."
}

show_test_instructions() {
  echo ""
  echo "--- Namespace Revised Test Execution Instructions (New IPs) ---"
  echo "1. Open another terminal and watch kernel logs:"
  echo "   sudo dmesg -wH  (or sudo journalctl -fk)"
  echo ""
  echo "2. In yet another terminal(s), run tcpdump on the router interfaces (root namespace):"
  echo "   sudo tcpdump -i $LAN_IF -n -e 'arp or icmp' -vv"
  echo "   sudo tcpdump -i $PDU_IF -n -e 'arp or icmp' -vv"
  echo ""
  echo "3. Test basic connectivity from client ns to its gateway (router's LAN IF):"
  echo "   sudo ip netns exec $NS_CLIENT ping -c 4 -W 2 $ROUTER_LAN_IP"
  echo "   Expected: Pings succeed. tcpdump on $LAN_IF shows ARP and ICMP."
  echo ""
  echo "4. If step 3 succeeds, test ping from client ns to PDU Gateway ($PDU_GATEWAY_IP):"
  echo "   sudo ip netns exec $NS_CLIENT ping -c 5 $PDU_GATEWAY_IP"
  echo "   Expected: Pings succeed. Logs show MANGLE, FORWARD_SPECIFIC, NAT_PRE_MASQ."
  echo "             tcpdump on $PDU_IF shows ICMP echo requests with source IP $ROUTER_WAN_IP (due to NAT)."
  echo ""
  echo "5. Test ping from client ns to the dummy External IP ($EXTERNAL_TEST_IP):"
  echo "   sudo ip netns exec $NS_CLIENT ping -c 5 $EXTERNAL_TEST_IP"
  echo "   Expected: Pings succeed. Similar logs and tcpdump behavior."
  echo ""
  echo "6. To tear down the environment, run:"
  echo "   sudo $0 cleanup"
  echo "---"
}

# --- Main Script Logic ---
check_root

if [ "$1" == "setup" ]; then
  cleanup # Clean before setup for idempotency
  setup_environment
  configure_router_logic
  show_test_instructions
elif [ "$1" == "cleanup" ]; then
  cleanup
else
  echo "Usage: sudo $0 {setup|cleanup}"
  exit 1
fi

exit 0