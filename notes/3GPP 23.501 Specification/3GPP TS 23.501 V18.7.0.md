# 4.2.8 Support of non-3GPP access
The following types of non-3GPP access networks are defined:
- Untrusted non-3GPP access networks;
- Trusted non-3GPP access networks; and
- Wireline access networks.
## General Concepts
The 5G Core Network supports connectivity of UEs via non-3GPP access networks, e.g. WLAN access networks.

Supports both untrusted non-3GPP access networks and trusted non-3GPP access networks (TNANs).

An **untrusted non-3GPP** access network **shall be connected** to the 5G Core Network **via a Non-3GPP InterWorking Function (N3IWF)**, via **N2** interface .

An untrusted non-3GPP access network shall be connected to the 5G Core Network via a Non-3GPP InterWorking Function (N3IWF), whereas a trusted non-3GPP access network shall be connected to the 5G Core Network via aTrusted Non-3GPP Gateway Function (TNGF).

A UE shall establish an IPsec tunnel with the N3IWF or with the TNGF in order to register with the 5G Core Network over non-3GPP access. Further details about the UE registration to 5G Core Network over untrusted non-3GPP access
and over trusted non-3GPP access are described in clause 4.12.2 and in clause 4.12.2a of TS 23.502 [3], respectively.
## 4.2.8.1A General Concepts to support Wireline Access
Wireline 5G Access Network (W-5GAN) shall be connected to the 5G Core Network via a **Wireline Access Gateway Function (W-AGF**). The W-AGF interfaces the 5G Core Network CP and UP functions via N2 and N3 interfaces,
respectively.
# 6.2.9 N3IWF
The functionality of N3IWF in the case of untrusted non-3GPP access includes the following:
- Support of IPsec tunnel establishment with the UE: The N3IWF terminates the IKEv2/IPsec protocols with the UE over NWu and relays over N2 the information needed to authenticate the UE and authorize its access to the 5G Core Network.
- Termination of N2 and N3 interfaces to 5G Core Network for control - plane and user -plane respectively.
- Relaying uplink and downlink control-plane NAS (N1) signalling between the UE and AMF.
- Establishment of IPsec Security Association (IPsec SA) to support PDU Session traffic.
- Relaying uplink and downlink user-plane packets between the UE and UPF. This involves:
	- De-capsulation/ encapsulation of packets for IPSec and N3 tunnelling.
# 6.2.20 WAGF
[SAGERAN W-AGF Device](https://www.sageran.com/products/network-equipments/w-agf.html)

W-AGF(Wireline Access Gateway Function ) is a wired access network function node based on 5G technology, defined by 3GPP and BBF. It acts as an intermediary between the RG and UPF, supporting N2 and N3 access to 5GC.

W-AGF is developed on the SageRAN's 5G Engine™ protocol stack platform and complies with 3GPP standards and the BBF WWC architecture, while adding processing logic for N1(NAS) protocol and SIM card information for network registration.

Additionally, W-AGF includes data routing and forwarding modules for registration and IP data link connections to the 5GC core network, providing standard fixed-mobile convergence network access technology for high-reliability, high-security, and unified authentication.

W-AGF can be widely applied in dedicated network links for industries such as hospitals, banks, campuses, and small and medium enterprises.
# 6.3.6 N3IWF selection
