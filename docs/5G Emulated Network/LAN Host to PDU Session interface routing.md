Traffic needs to be routed to specific interfaces based on the source host.

![[topology-LAN-Host-to-PDU-Session-Interface-Routing.png]]

What is needed:
- Mark packets originating from `192.168.1.2` with a specific mark.
- Mark packets originating from `192.168.1.3` with a different mark.
- Create custom routing tables for `pdu1` and `pdu2`.
- Create routing rules that send packets with the appropriate marks to the corresponding routing tables.
- NAT will be needed for both WAN interfaces.

```bash
#Enable ip forwarding
echo 1 > /proc/sys/net/ipv4/ip_forward

# Mangle table: mark packets
iptables -t mangle -A PREROUTING -s 192.168.1.2 -j MARK --set-mark 10
iptables -t mangle -A PREROUTING -s 192.168.1.3 -j MARK --set-mark 20

# Create routing tables (add to /etc/iproute2/rt_tables)
# 100 table10
# 200 table20

# Add routes to tables (replace gateway IPs)
ip route add default via pdu1 dev eth0 table table10
ip route add default via pdu2 dev eth4 table table20

# Add routing rules
ip rule add fwmark 10 table table10
ip rule add fwmark 20 table table20
```
