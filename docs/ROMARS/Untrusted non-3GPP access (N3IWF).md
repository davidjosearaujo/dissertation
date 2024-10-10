In 3GPP release 15 of 5G \[1\], the support for non-3GPP access to the Core Network (i.e., avoiding the RAN part and using other types of access) is described, defining the Non-3GPP Inter-Working Function (N3IWF) Network Function. More specifically, untrusted Wi-Fi (802.11) access is considered, since the large majority of UEs have both 3GPP radio interface and Wi-Fi interface. N3IWF is in charge to “adapt” access and authentication protocols of the non-3GPP WiFi section towards the 5G CN, interfacing with the CN with the same role of the RAN (i.e., enabling N2 and N3).

**Untrusted refers to the fact that non-3GPP access (i.e., the WiFi Access Point) is assumed not managed by the same operator of the 5G network: as an example an hotspot into an airport, or a private home Wi-Fi router**. For this reason an end-to-end security association between UE and N3IWF shall be established, regardless of the one established in the access at layer 2 (i.e., WPA2).

The goal of 3GPP is to define for 5G, since the beginning, an interoperable, straightforward and reliable non-3GPP access mechanism which will most likely be adopted on a wider scale. The resulting architecture is shown in figure, with the UE that can simultaneously enable the NR-Uu access via gNB, or have WiFi as unique access method. In either cases, 3GPP credentials to initiate the access (i.e., USIM/eUICC) are always **required**.

![[Pasted image 20240926175108.png]]

N3IWF mainly provides a secure gateway to operator’s 5G network for non-3GPP access. The interface NWu between UE and N3IWF is based on IPSec/IKE to establish a secure tunnel, by-passing the security mechanisms enforced by the Access Point (if any). As UE is expected to communicate with AMF over the NAS interface, N3IWF has a N2 interface connecting with AMF to enable N1. Then, N3 interface to interact with UPF is included as well in N3IWF.

N3IWF is responsible for setting up the IPSec connection to be used by control plane traffic directed to AMF/SMF, as well as the traffic directed to the UPF for the user plane. As a consequence UE and N3IWF need to establish two IPSec Security Associations (SAs):  
• Signalling (control plane) IPSec SA – it transports NAS messages destined to AMF  
• User plane IPSec SA – it transports packets destined to DN

In the first step of access operations, **UE and N3IWF must establish the main signalling IPSec SA**, using the IKE protocol which allows to **support EAP registration of the UE in the 5G system**,similarly and with the same security requirements of the 3GPP access. Figure 6 presents a control plane protocol’s stack used to establish signalling IPSec tunnel.

![[Pasted image 20240926175210.png]]
*Figure 6: N3IWF protocol stack before SA establishment – control plane*

When the signalling IPSec SA is established, the IPSec tunnel (i.e., ESP in tunnel mode) can deliver EAP-5G messages, used to encapsulate NAS messages between UE and N3IWF. 3GPP specification says that EAP-5G is identified by a specific EAP extension header “Expanded Type” (0xFE) with specified values for “Vendor ID” and “Vendor Type”. At this stage, UE can communicate with AMF to perform end-to-end NAS signalling as presented in Figure 7, as it would normally do via RAN. Note that the IPsec layer includes as well the “Inner IP” header, required to establish the IPSec tunnel.

![[Pasted image 20240926175453.png]]
*Figure 7: N3IWF protocol stack after SA Establishment – control plane*

Once UE is registered in the 5G system and enables its N1/NAS interface, it can negotiate the PDU session establishment which results in the establishment of a child IPSec SA (called user plane IPSec SA) to communicate with the Data Networks (DN). The resulting user plane protocol’s stack is depicted in Figure 8. The GRE (Generic Routing Encapsulation) protocol is used to carry user PDU (IP) between UE and N3IWF. GRE allows in particular to implement a flow-based QoS model as specified in TS 23.501, carrying QFI (QoS Flow Identifier) and QRI (QoS Reflective Identifier) associated with user data packets as defined by 3GPP specifications. GRE is also supporting the other types of PDU foreseen (i.e, Ethernet). Also the user plane IPSec tunnel as well includes an “Inner IP” header.

The adaptation of WLAN to access a 3GPP network is then requiring a relatively higher overhead with regard to direct use of the same technology by the UE: 2 independent IP layers, plus 1 GRE header plus 1 IPSec header with inner IP header (for ESP in tunnel mode) instead than a single IP header.

![[Pasted image 20240926175520.png]]
*Figure 8: N3IWF protocol stack after SA Establishment – user plane*

Note that current implementations of 5G networks are based on rel-15 and on a legacy LTE Core Network (i.e., Non-Stand-Alone, NSA) configuration, therefore the N3IWF is normally neither available in the core network nor available in COTS UEs as of today.