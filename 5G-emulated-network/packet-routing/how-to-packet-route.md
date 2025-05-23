# Policy-Based Router Configuration Guide
T
his guide outlines the steps to configure a Linux machine to act as a router and implement policy-based routing (PBR). The goal is to route internet traffic from a specific device on your Local Area Network (LAN), identified by its MAC address, through a designated WAN/PDU (Packet Data Unit) interface (e.g., uesimtunx), while other traffic might use a different default route.

This configuration is based on the logic within the configure_router_logic function of the provided test script.

## I. Prerequisites & Placeholders

Before you begin, identify the following for your router:
- **LAN Interface** (YOUR_LAN_IF): `enp0s9`
- **PDU Interface** (YOUR_PDU_IF): `uesimtunx`
  - **Note**: Replace `x` in `uesimtunx` with the actual PDU session number (e.g., `uesimtun0', `uesimtun5`).
- **Client MAC Address** (YOUR_CLIENT_MAC): The MAC address of the client device on your LAN whose traffic you want to route (e.g., `08:00:27:24:96:ec`).
- **PDU Gateway IP** (YOUR_PDU_GATEWAY_IP): `10.46.0.1`
- **Firewall Mark Value** (YOUR_PDU_SESSION_ID_AS_MARK): This is the actual numerical PDU Session ID.
  - **Example**: If your PDU Session ID is `5`, then YOUR_PDU_SESSION_ID_AS_MARK will be `5`.
- **Custom Routing Table ID** (YOUR_TABLE_ID_NUMERIC): 200_PLUS_PDU_ID (This must be a number)
  - **Note**: Replace 200_PLUS_PDU_ID with an actual unique numerical ID for your routing table, typically derived from the PDU Session ID. For example, if your PDU Session ID is `5`, you might use `205` (i.e., 200 + 5).
- **Custom Routing Table Name** (YOUR_TABLE_NAME_SYMBOLIC): TABLE_FOR_PDU_ID
  - **Note**: Replace TABLE_FOR_PDU_ID with a descriptive symbolic name for your routing table. For example, if your PDU Session ID is 5, you might use table_pdu_5.

Important: Ensure your LAN (enp0s9) and PDU (uesimtunx) interfaces are up and have IP addresses configured.

## II. Configuration Steps
Execute these commands on your Linux router as root or with sudo. Remember to replace placeholders like YOUR_CLIENT_MAC, YOUR_PDU_SESSION_ID_AS_MARK (with the actual PDU Session ID number), YOUR_TABLE_ID_NUMERIC (with a number like 200 + PDU Session ID), and YOUR_TABLE_NAME_SYMBOLIC (with a name like table_pdu_X) as per the notes above.

### Phase 1: Enable IP Forwarding
This allows your system to pass packets between interfaces.

```bash
# Enable IPv4 forwarding
sudo sysctl -w net.ipv4.ip_forward=1

# To make it persistent across reboots, edit /etc/sysctl.conf (or a file in /etc/sysctl.d/)
# and add or uncomment
# net.ipv4.ip_forward=1
# Then apply with: sudo sysctl -p
```

### Phase 2: Setup Custom Routing Infrastructure

This involves creating a new routing table and a rule to direct marked traffic to it.

1. Define the Custom Routing Table:

Add an entry to /etc/iproute2/rt_tables. This file maps table numbers to names.

```bash
# Replace YOUR_TABLE_ID_NUMERIC (e.g., 205) and YOUR_TABLE_NAME_SYMBOLIC (e.g., table_pdu_5)
# Example: echo "205 table_pdu_5" | sudo tee -a /etc/iproute2/rt_tables
if ! grep -q -E "^\s*YOUR_TABLE_ID_NUMERIC\s+YOUR_TABLE_NAME_SYMBOLIC\s*$" /etc/iproute2/rt_tables; then
  echo "YOUR_TABLE_ID_NUMERIC YOUR_TABLE_NAME_SYMBOLIC" | sudo tee -a /etc/iproute2/rt_tables
else
  echo "Entry YOUR_TABLE_ID_NUMERIC YOUR_TABLE_NAME_SYMBOLIC already exists in /etc/iproute2/rt_tables."
fi
```

2. Add a Default Route to the Custom Table:

This route specifies that traffic using this table should go out via your PDU interface (uesimtunx), through its gateway (10.46.0.1).

```bash
# Replace YOUR_TABLE_NAME_SYMBOLIC. (And ensure 'x' in uesimtunx is correct)
# Example: sudo ip route add default via 10.46.0.1 dev uesimtun5 table table_pdu_5
sudo ip route add default via 10.46.0.1 dev uesimtunx table YOUR_TABLE_NAME_SYMBOLIC
```

3. Create an IP Rule:

This rule tells the system to use your custom routing table for packets marked with YOUR_PDU_SESSION_ID_AS_MARK.

```bash
# Replace YOUR_PDU_SESSION_ID_AS_MARK (with the actual PDU Session ID number) and YOUR_TABLE_NAME_SYMBOLIC
# Example: sudo ip rule add fwmark 5 table table_pdu_5
sudo ip rule add fwmark YOUR_PDU_SESSION_ID_AS_MARK table YOUR_TABLE_NAME_SYMBOLIC
```

### Phase 3: Packet Marking (iptables Mangle Table)

These rules mark incoming packets from your specific client MAC address on the enp0s9 interface with the PDU Session ID.

```bash
# Replace YOUR_CLIENT_MAC and YOUR_PDU_SESSION_ID_AS_MARK (with the actual PDU Session ID number)
# Example: sudo iptables -t mangle -A PREROUTING -i enp0s9 -m mac --mac-source 08:00:27:24:96:ec -j MARK --set-mark 5
# Mark packets from the specific client MAC address
sudo iptables -t mangle -A PREROUTING -i enp0s9 -m mac --mac-source YOUR_CLIENT_MAC -j MARK --set-mark YOUR_PDU_SESSION_ID_AS_MARK
```

> Optional: The script included logging rules here for debugging. You can add them if needed:

```bash
# sudo iptables -t mangle -A PREROUTING -i enp0s9 -m mac --mac-source YOUR_CLIENT_MAC -j LOG --log-prefix "MANGLE_PRE_PBR: " --log-level 4
# sudo iptables -t mangle -A PREROUTING -i enp0s9 -m mac --mac-source YOUR_CLIENT_MAC -m mark --mark YOUR_PDU_SESSION_ID_AS_MARK -j LOG --log-prefix "MANGLE_POST_PBR: " --log-level 4
```

### Phase 4: Firewall Forwarding Rules (iptables Filter Table)

These rules control which packets are allowed to be forwarded between interfaces.

1. Set Default FORWARD Policy to DROP (Recommended for Security):

This ensures only explicitly allowed traffic is forwarded.

```bash
sudo iptables -P FORWARD DROP
```

2. Allow Established and Related Connections:

This is crucial for allowing return traffic for connections initiated from your LAN.

```bash
sudo iptables -A FORWARD -m state --state RELATED,ESTABLISHED -j ACCEPT
```

> Optional: Logging for this rule:

```bash
# sudo iptables -A FORWARD -m state --state RELATED,ESTABLISHED -j LOG --log-prefix "FORWARD_REL_EST_PBR: " --log-level 4
```

3. Allow Forwarding for Marked Traffic from Specific Client to PDU Interface:

This rule explicitly permits the marked traffic (using PDU Session ID as mark) from your client (YOUR_CLIENT_MAC) on enp0s9 to be forwarded out the uesimtunx interface.

```bash
# Replace YOUR_CLIENT_MAC and YOUR_PDU_SESSION_ID_AS_MARK (with the actual PDU Session ID number). (And ensure 'x' in uesimtunx is correct)
# Example: sudo iptables -A FORWARD -i enp0s9 -o uesimtun5 -m mac --mac-source 08:00:27:24:96:ec -m mark --mark 5 -j ACCEPT
sudo iptables -A FORWARD -i enp0s9 -o uesimtunx -m mac --mac-source YOUR_CLIENT_MAC -m mark --mark YOUR_PDU_SESSION_ID_AS_MARK -j ACCEPT
```

> Optional: Logging for this rule:

```bash
# sudo iptables -A FORWARD -i enp0s9 -o uesimtunx -m mac --mac-source YOUR_CLIENT_MAC -m mark --mark YOUR_PDU_SESSION_ID_AS_MARK -j LOG --log-prefix "FORWARD_SPECIFIC_PBR: " --log-level 4
```

### Phase 5: NAT Rule (iptables NAT Table)

This rule applies Network Address Translation (NAT) to traffic going out your uesimtunx interface, making it appear to originate from the router's uesimtunx interface IP address. This is usually necessary for internet access.

```bash
# (Ensure 'x' in uesimtunx is correct)
# Example: sudo iptables -t nat -A POSTROUTING -o uesimtun5 -j MASQUERADE
sudo iptables -t nat -A POSTROUTING -o uesimtunx -j MASQUERADE
```

> Optional: Logging for this rule:

```bash
# sudo iptables -t nat -A POSTROUTING -o uesimtunx -j LOG --log-prefix "NAT_PRE_MASQ_PBR: " --log-level 4
```

## III. Making Rules Persistent

The sysctl, ip route, ip rule, and iptables commands applied above are not persistent across reboots by default.

- `sysctl` (IP Forwarding): Edit `/etc/sysctl.conf` or a file in `/etc/sysctl.d/` and run `sudo sysctl -p`.
- `/etc/iproute2/rt_tables`: This file is usually read at boot, so changes here are typically persistent once saved.
- `ip route` and `ip rule`: These need to be reapplied on boot. This can be done using:
  - Network manager dispatcher scripts (e.g., for NetworkManager or systemd-networkd).
  - Scripts in `/etc/network/if-up.d/` if using traditional Debian networking.
  - A custom systemd service unit.
- iptables Rules:
  - Use the iptables-persistent package (Debian/Ubuntu):

```bash  
sudo apt-get install iptables-persistent`
# After setting rules
sudo netfilter-persistent save
# Or: sudo iptables-save > /etc/iptables/rules.v4
```

## IV. Verification

After applying these rules:

1. Check `ip rule list` to see your fwmark rule.
2. Check `sudo ip route show table YOUR_TABLE_NAME_SYMBOLIC` (e.g., sudo ip route show table table_pdu_5) for the default route.
3. Check `sudo iptables -t mangle -L PREROUTING -v -n`, `sudo iptables -L FORWARD -v -n`, and `sudo iptables -t nat -L POSTROUTING -v -n` to see your iptables rules and their packet/byte counters.
4. Test connectivity from the client device (YOUR_CLIENT_MAC). Its traffic should now egress via uesimtunx. You can use tools like traceroute (or mtr) from the client to verify the path.
5. Check kernel logs (dmesg or journalctl -fk) if you included the optional LOG rules.

This guide should provide a solid foundation for setting up your policy-based router. Remember to adapt the placeholder values to your specific environment.t
