- Different subnets for each PDU Session can be defined in the SMF and UPF.
- In the UPF the interface attributed to each must be defined.
- Apparently, the UPF configuration does not allow for Open5GS to automatically create the tun interfaces.
	- [x] Tun interfaces must be created manually: `ip tuntap add dev intun0 mode tun`
	- [x] Tun interfaces must be brought up manually: `ip link set intun0 up`
	- [x] Open5GS automatically attributes the IP based on the config

# Design a topology diagram for my POC
## Core
- Interfaces
	- enp0s8
		- Tunnel Interfaces
			- intun0
			- intun1
			- intun2
- Slice 1
	- PDU Sessions
		- PDU Session 0
		- PDU Session 1
		- PDU Session 2
## gNB
- Slices
	- SST 1
- Interfaces
	- enp0s8
	- enp0s9

# Implement EAP-TLS
 - [ ] Implement EAP-TLS with NAUN3 as suplicant and UE as authenticator (and could also be the authentication server for now).