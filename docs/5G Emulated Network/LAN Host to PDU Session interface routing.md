Traffic needs to be routed to specific interfaces based on the source host.

![[topology-LAN-Host-to-PDU-Session-Interface-Routing.png]]

# Requirements
- Mark packets originating from `192.168.60.100` with a specific mark.
- Mark packets originating from `192.168.60.200` with a different mark.
- Create custom routing tables for `pdu1` and `pdu2`.
- Create routing rules that send packets with the appropriate marks to the corresponding routing tables.
- NAT will be needed for both WAN interfaces.
# How to recreate
## Interfaces per VM
### Hosts
One interface each: `192.168.60.x`
### UE
Two interfaces, one for LAN and the other for WAN (Core)
- LAN interface: `lan0`
	- `192.168.60.2`
- WAN interface: `wan0`
	- `192.168.56.100`
	- 2 VLAN interfaces:
		- `wan0.pdu1`: `10.46.0.100`
		- `wan0.pdu2`: `10.46.0.200`

```bash
#Enable ip forwarding
echo 1 > /proc/sys/net/ipv4/ip_forward

# Mangle table: mark packets
iptables -t mangle -A PREROUTING -s 192.168.60.100 -j MARK --set-mark 10
iptables -t mangle -A PREROUTING -s 192.168.60.200 -j MARK --set-mark 20

# Create routing tables (add to /etc/iproute2/rt_tables)
# 100 table10
# 200 table20

# Add routes to tables (replace gateway IPs)
ip route add default via 10.46.0.1 dev wan0.pdu1 table table10
ip route add default via 10.46.0.2 dev wan0.pdu2 table table20

# Add routing rules
ip rule add fwmark 10 table table10
ip rule add fwmark 20 table table20
```
### CORE
One interface:
- `wan0`:`192.168.56.2`
- 2 VLAN interfaces:
	- `wan0.pdu1`: `10.46.0.1`
	- `wan0.pdu1`: `10.46.0.2`