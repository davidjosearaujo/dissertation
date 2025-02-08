- Different subnets for each PDU Session can be defined in the SMF and UPF.
- In the UPF the interface attributed to each must be defined.
- Apparently, the UPF configuration does not allow for Open5GS to automatically create the tun interfaces.
	- [x] Tun interfaces must be created manually: `ip tuntap add dev intun0 mode tun`
	- [x] Tun interfaces must be brought up manually: `ip link set intun0 up`
	- [x] Open5GS automatically attributes the IP based on the config

 - [ ] Design a topology diagram for my POC
 - [ ] Implement EAP-TLS with NAUN3 as suplicant and UE as authenticator (and could also be the authentication server for now).