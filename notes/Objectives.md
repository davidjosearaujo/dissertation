Improve and extend an existing prototype of the target Interworking Function (IWF). The IWF corresponds to a BBF AGF for FN-RG but presenting itself to the 5G Core as a gNB, thus **eliminating the need for the 5G Cores to recognize other Access Node types**. It promotes, in a simplified and universal way, convergence between 5G and Ethernet
based accesses, e.g. WLAN.
# Existing prototype

The **existing IWF does not required devices to establish IPSec tunnels** and mobile **terminals to have 5G credentials**. For that purpose, the IWF is **provisioned with a set of 5G credentials**, also provisioned at the 5G core, which are used whenever a non 5G capable terminal requires data access. That process is triggered by Dynamic Host Configuration Protocol (DHCP), requiring the IWF to be the DHCP server of the IP segment the device is connected to, being the presented **MAC address used as the ID of the device and mapped to an available 5G identifier**. This mapping can be made static. Currently, only connection to a single slide is considered.

It follows Software Defined Networking (SDN) principles and mechanisms to map WLAN
connections to 5G sessions, under the common control of a single instance of 5G Core composed of three main blocks, which are depicted in Figure 33:
1. **SDN controller and interworking function emulation**
	1. Responsible for **intercepting protocol packets on the LAN** (e.g. DHCP), **mapping clients** (MAC addresses) into particular SUPIs and **trigger the allocation of UEs** on the 5G network.
	2. Controls the data plane for **bridging traffic into the UPF**.
	3. It is externally controlled via REST, emulates a router/ Broadband Network Gateway (BNG) on the Fixed Network side.
2. **gRPC**
	1. **Responsible for receiving the allocations and deallocations on the 5G network.** Interacts with the RAN, emulating the Next Generation Application Protocol (NGAP) protocol.
3. **Data plane**
	1. Bridges traffic from the LAN to the UPF (layer 3), via standard N3 interface, which requires General Packet Radio Service Tunnelling Protocol (GTP) encapsulation.

From above, there is no security mechanisms to properly authenticate and authorize non 5G
devices. Only the observed MAC address, which can be forged, is used to identify the device and assign to it a 5G identity.
![[Pasted image 20240927124739.png]]
# Functional expansion

Extended to improve security and extend 5G slicing concept with WLAN SSIDs or Ethernet VLAN as described in the following.

## Improve non-5G devices authentication and authorization (AA) mechanisms

N3IWF, TNGF and TWIF require devices to have 5G credentials in order to have a single point of provisioning.
The foreseen solution is based on the interception of the RADIUS traffic between the authenticating entity and the respective RADIUS server. This allows the extraction of the device identity and obtain the result of the process.

## Implement ‘extended slices’

Static mapping between devices IDs and 5G IDs must be supported and provisioned in the IWF.

![[Pasted image 20240927153324.png]]