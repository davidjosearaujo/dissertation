# 4.2.2 Registration Management procedures
Aspects related to dual registration in 3GPP and non-3GPP access are described in clause [[#4.12 Procedures for Untrusted non-3GPP access|4.12]]. The general Registration call flow in clause [[#4.2.2.2.2 General Registration|4.2.2.2.2]] is also used for the case of registration in 3GPP access when the UE is already registered in a non-3GPP access and vice versa. Registration in 3GPP access when the UE is already registered in a non-3GPP access scenario may require an AMF change, as further detailed in clause 4.12.8.
### 4.2.2.2.2 General Registration
![[4_2_2_2_2_general_registration_procedure.png]]
*General Registration procedure*
# 4.12.2 Registration via Untrusted non-3GPP Access
![[Pasted image 20241017151403.png]]
Clause 4.12.2 specifies how a UE can register to 5GC via an untrusted non-3GPP Access Network. It is based on the Registration procedure specified in clause [[#4.2.2.2.2 General Registration|4.2.2.2.2]] and it uses a vendor-specific EAP method called "EAP-5G". The EAP-5G packets utilize the "Expanded" EAP type and the existing 3GPP Vendor-Id registered with IANA under the SMI Private Enterprise Code registry. The "EAP-5G" method is used between the UE and the N3IWF and is utilized only for encapsulating NAS messages (not for authentication). If the UE needs to be authenticated, mutual authentication is executed between the UE and AUSF. The details of the authentication procedure are specified in [[Security architecture and procedures for 5G system - 33.501|TS 33.501]]
![[registration_via_untrusted_non-3GPP_access.png]]
1. The **UE connects to an untrusted non-3GPP Access Network** with any appropriate authentication procedure and it is assigned an IP address. For example, **a non-3GPP authentication method can be used**, e.g. no authentication (in the case of a free WLAN), EAP with pre-shared key, username/password, etc.
	1. When the UE decides to attach to 5GC network, **the UE** not operating in SNPN access mode for NWu interface **selects an N3IWF in a 5G PLMN**
2. The UE proceeds with the establishment of an IPsec Security Association (SA) with the selected N3IWF by initiating an IKE initial exchange.
3. ...
# 4.12a.2 Registration via Trusted non-3GPP Access
![[Pasted image 20241017151444.png]]
In this case, the "EAP-5G" method is used between the UE and the TNGF and is utilized for encapsulating NAS messages.

In Registration and subsequent Registration procedures via trusted non-3GPP access, the NAS messages are always exchanged between the UE and the AMF. When possible, the UE can be authenticated by reusing the existing UE security context in AMF.
![[Registration via Trusted non-3GPP Access.png]]
#### Step 1
A layer-2 connection is established between the UE and the TNAP. In the case of IEEE Std 802.11, this step corresponds to an 802.11 Association. In the case of PPP, this step corresponds to a PPP LCP negotiation. In other types of non-3GPP access (e.g. Ethernet), this step may not be required.
#### Step 2-3
An EAP procedure is initiated.
#### Step 4-10
An EAP-5G procedure is executed as the one specified in clause 4.12.2.2 for the untrusted non-3GPP access with the following modifications:
- The registration request may contain an indication that the UE supports TNGF selection based on the slices the UE wishes to use over trusted non-3GPP access (i.e. that the UE supports Extended WLANSP rule).
- A TNGF key (instead of an N3IWF key) is created in the UE and in the AMF after the successful authentication. The TNGF key is transferred from the AMF to TNGF in step 10a (within the N2 Initial Context Setup Request). The TNGF derives a TNAP key, which is provided to the TNAP. The TNAP key depends on the non-3GPP access technology (e.g. it is a Pairwise Master Key in the case of IEEE Std 802.11).
# 4.12b  Procedures for devices that do not support 5GC NAS over WLAN access
![[Pasted image 20241017151526.png]]
**Devices that do not support 5GC NAS** signalling over WLAN access (referred to as "Non-5G-Capable over WLAN" devices, or N5CW devices for short), **may access 5GC** in a PLMN or an SNPN **via a trusted WLAN Access Network** that supports a Trusted WLAN Interworking Function (TWIF). The following clause specifies how a N5CW device can be registered to 5GC and how it can send data via a PDU Session.
![[Pasted image 20241017151801.png]]
#### Step 2
The N5CW device provides its Network Access Identity (NAI). The Trusted WLAN Access Point (TWAP) selects a Trusted WLAN Interworking Function (TWIF), e.g. based on the received realm and sends an AAA request to the selected TWIF.
#### Step 3
The TWIF creates a 5GC Registration Request message on behalf of the N5CW device. **The TWIF uses default values to populate the parameters in the Registration Request message, which are the same for all N5CW devices.** The Registration type indicates "Initial Registration".

If the TWIF receives a Decorated NAI, in Registration Request message the TWIF send the NAI which corresponds to the HPLMN by removing the decoration, for example `NAI=type1.rid678.schid0.useriduser17@ nai.5gc-nn.mnc<MNC_Home>.mcc<MCC_Home>.3gppnetwork.org`
