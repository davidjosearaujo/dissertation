# Definitions
- **Wireline access Control Plane protocol (W-CP)**: Protocol used to transport AS and NAS signalling between the 5G- RG and the W-AGF over the Y4 reference point. W-CP is specified by BBF and CableLabs. There is no assumption that W-CP refers to only a single protocol or only a specific protocol layer.
- **Wireline access User Plane protocol (W-UP)**: Protocol used to carry PDU Session user plane traffic between the 5G- RG and the W-AGF over the Y4 reference point. W-UP is specified by BBF and CableLabs. There is no assumption that W-UP refers to only a single protocol or only a specific protocol layer.

# 4.2.2 Identification and authentication
In the case of **FN-RG connected via W-5GAN**, the specification defined in [[System architecture for the 5G System (5GS) - 23.501|TS 23.501 clause 5.2.3]] applies with the following differences:
- UE is replaced by FN-RG
- The **W-AGF provides the NAS signalling** connection to the 5GC **on behalf of the FN-RG**.
- The W-5GAN may authenticate the FN-BRG per BBF specification BBF TR-456 and WT-457. The W- 5GAN may authenticate the FN-CRG per [[WR-TR-5WWC-ARCH-V01-190820|CableLabs DOCSIS MULPI]].

# 4.2.3 Authorization
In the case of FN-RG connected via W-5GAN, the specification defined in TS 23.501 [2] clause 5.2.4 applies with the following differences:
- UE is replaced by FN-RG
- **W-AGF performs the UE Registration procedure on behalf of the FN-RG**

# 4.3.1 Registration management
Registration management when 5G-RG or FN-RG is connected to 5GC via wireline access is described in [[System architecture for the 5G System (5GS) - 23.501#5.5.1 Registration Management|TS 23.501 clause 5.5.1]]

# 4.3.2 Connection management
Connection management when 5G-RG or FN-RG is connected to 5GC via wireline access is described in [[System architecture for the 5G System (5GS) - 23.501#5.5.2 Connection Management|clause 5.5.2 of TS 23.501]]

# 4.4.2 Session management for FN-RG
Session management of FN-RG follows the principle defined in TS 23.501 clause 5.6 with the following differences:
- UE is replaced by W-AGF
- FN-RG is connected to 5GC via wireline access instead of 3GPP access

# 4.10a Non-5G capable device behind 5G-CRG and FN-CRG
non-5G capable (N5GC) devices connecting via W-5GAN can be authenticated by the 5GC using EAP based authentication method(s) as defined in TS 33.501.

**Roaming is not supported for N5GC devices.**

![[2024-10-07_12-39.png]]
1. The W-AGF registers the FN-CRG to 5GC as specified in [[Wireless and wireline convergence access support for the 5G System (5GS) - 23.316#7.2.1.3 FN-RG Registration via W-5GAN|clause 7.2.1.3]] or the 5G-CRG registers to 5GC as specified in [[Wireless and wireline convergence access support for the 5G System (5GS) - 23.316#7.2.1.1 5G-RG Registration via W-5GAN|clause 7.2.1.1]]
   
2. The CRG is configured as L2 bridge mode and forwards any L2 frame to W-AGF. **802.1x authentication may be triggered**. This can be done either by N5GC device sending a EAPOL-start frame to W-AGF or W-AGF receives a frame from an unknown MAC address.
   
   How the CRG is configured to work in L2 bridge mode and how the W-AGF is triggered to apply procedures for N5GC devices is defined in [[WR-TR-5WWC-ARCH-V01-190820|CableLabs WR-TR-5WWC-ARCH]]
   The N5GC device send an EAP-Resp/Indentity including its Network Access Identifier (NAI) in the form of username@realm.
   
3. **W-AGF, on behalf of the N5GC device, sends a NAS Registration Request message to AMF with a device capability indicator that the device is non-5G capable. For this purpose, the W-AGF creates a NAS Registration Request message containing a SUCI. The W-AGF constructs the SUCI from the NAI received within EAP- Identity from the N5GC device as defined in TS 33.501**
   Over N2 there is a separate NGAP connection per N5GC device served by the W-AGF.
   When it provides (over N2) ULI to be associated with a N5GC device, the W-AGF builds the N5GC's ULI using the GCI (see clause 4.7.9) of the CRG connecting the N5GC device.
   
4. AMF selects a suitable AUSF as specified in TS 23.501 clause 6.3.4.
   
5. **EAP based authentication defined in TS 33.501 is performed between the AUSF and N5GC device.**
   Once the N5GC device has been authenticated, the AUSF provides relevant security related information to the AMF. AUSF shall return the SUPI (this SUPI corresponds to a NAI that contains the username of the N5GC device and a realm as defined in TS 33.501) to AMF only after the authentication is successful.
   
> Each N5GC device is registered to 5GC with its own unique SUPI.

6. The AMF performs other registration procedures as required (see TS 23.502 clause 4.2.2.2.2).
   When providing a PEI for a N5GC device, the **W-AGF shall provide a PEI containing the MAC address of the N5GC device**. The **W-AGF may**, based on operator policy, **encode the MAC address of the N5GC device using the IEEE Extended Unique Identifier EUI-64 format** (see IEE Publication).
   
7. **The AMF sends Registration Accept message to W-AGF.**
   Once the registration procedure is completed, the **W-AGF requests the establishment of a PDU Session on behalf of the N5GC device**. Only one PDU session per N5GC device is supported. The procedure is the same as the PDU Session establishment procedure specified in clause 7.3.4 with the difference as below:

After successful registration, PDU Session establishment/modification/release procedure specified in clause 7.3.4, 7.3.6, and 7.3.7 apply with the difference as below:
- **FN-RG is replaced by N5GC device**.

The W-AGF shall request the release of the NGAP connection for each N5GC device served by a CRG whose NGAP connection has been released.

5G-CRG behaves as FN-CRG (i.e. L2 bridge mode) when handling N5GC devices.

# 4.13 Support of FN-RG
FN-RG is a legacy type of residential gateway that does not support N1 signalling and is not 5GC capable.

Support for FN-RG connectivity to 5GC is provided by means of W-AGF supporting 5G functionality on behalf of the FN-RG

The W-AGF supports the following functionality on behalf of the FN- RG:
- Has access to configuration information, as defined in [[TR-456-AGF-Functional-Requirements.pdf|BBF TR-456]], WT-457 and [[WR-TR-5WWC-ARCH-V01-190820|CableLabs WR-TR- 5WWC-ARCH]], to be able to serve FN-RGs and to construct AS and NAS messages.
- Acting as end-point of N1 towards AMF, including maintaining CM and RM states and related dynamic information received from 5GC
- Mapping between Y5 towards FN-RG and N1/N2 towards 5GC as well as mapping between a Y5 user plane connection and a PDU Session user plane tunnel on N3.

Authentication of FN-RG may be done by the W-AGF, as defined by BBF and [[WR-TR-5WWC-ARCH-V01-190820|Cablelabs]]. The W-AGF provides an indication on N2 that the FN-RG has been authenticated.

# 5.1 Network Function Functional description
## 5.1.1 W-AGF
The functionality of W-AGF in the case of Wireline 5G Access network includes the following:
- Termination of N2 and N3 interfaces
- Handling of N2 signalling related to PDU Sessions and QoS
- Relaying uplink and downlink user-plane packets between the 5G-RG and UPF and between FN-RG and UPF. This involves:
	- Enforcing QoS corresponding to N3 packet marking
	- N3 user-plane packet marking
- Supporting AMF discovery and selection
- Termination of wireline access protocol on Y4 and Y5
- In the case of FN-RG the W-AGF acts as end point of N1 on behalf of the FN-RG

# 6.2.2 Control Plane Protocol Stacks between the FN-RG and the 5GC
![[2024-10-08_11-29.png]]
The control plane protocol stack between FN-RG and AMF is defined in figure 6.2.2-1. The W-AGF acts as an N1 termination point on behalf of FN-RG.

For W-5GBAN, the L-W-CP protocol stack, between FN-BRG and W-AGF is defined in [[TR-456-AGF-Functional-Requirements.pdf|BBF TR-456]] and WT-457. For W-5GCAN, the L-W-CP protocol stack between FN-CRG and W-AGF is defined in [[WR-TR-5WWC-ARCH-V01-190820|WR-TR-5WWC-ARCH]]

# 6.3.2 User Plane Protocol Stacks between the FN-RG and the 5GC
![[2024-10-08_11-32.png]]
The user plane protocol stack between FN-RG and UPF is defined in figure 6.3.2-1.

For W-5GBAN, the L-W-UP protocol stack, between FN-BRG and W-AGF is defined in [[TR-456-AGF-Functional-Requirements.pdf|BBF TR-456]] and WT-457. For W-5GCAN, the L-W-UP protocol stack between FN-CRG and W-AGF is defined in [[WR-TR-5WWC-ARCH-V01-190820|WR-TR-5WWC-ARCH]]

# 7.2.1.3 FN-RG Registration via W-5GAN
![[Pasted image 20241008115450.png]]

1. The **FN-RG connects to a W-AGF (W-5GAN) via a layer-2 (L2)** connection, based on Wireline AN specific procedure.
   The FN-RG is authenticated by the W-5GAN based on Wireline AN specific mechanisms.
   
2. W-AGF selects an AMF based on the AN parameters and local policy. W-AGF may use the Line ID / HFC identifier provided from the Wireline AN to determine the 5GC and AN parameters to be used for the FN-RG registration. How the W-AGF can determine the necessary 5GC and AN parameters is defined in [[TR-456-AGF-Functional-Requirements.pdf|BBF TR-456]], WT-457 or [[WR-TR-5WWC-ARCH-V01-190820|CableLabs WR-TR-5WWC-ARCH]]
   
3. W-AGF performs initial registration on behalf of the FN-RG to the 5GC. The W-AGF **sends a Registration Request** to the selected AMF within an N2 initial UE message (NAS Registration Request, ULI, Establishment cause, UE context request, Allowed NSSAI, Authenticated Indication).
   
   The **NAS Registration Request contains the SUCI (Subscription Concealed Identifier)** or 5G-GUTI of the FN-RG, security parameters/UE security capability, UE MM Core Network Capability, PDU Session Status, Follow-on request, Requested NSSAI. The 5G-GUTI, if available, has been received from the AMF during a previous registration and stored in W-AGF.
   
   The **SUCI is built by the W-AGF** based on:
   - In the case of a BBF access: the **GLI** as defined in clause 4.7.8 **together with an identifier of the Home network** as described in TS 23.003.
   - In the case of a Cable access: the **GCI** as defined in clause 4.7.8 together with an identifier of the Home network as described in TS 23.003.

![[Pasted image 20241008121040.png]]

   The following differences exist, compared to 5G-RG case:
   - The W-AGF use SUCI
   - The Authenticated Indication indicates to AMF and 5GC that the FN-RG has been authenticated by the access network.

4. **If the AMF receives a SUCI, the AMF shall select an AUSF** as specified in TS 23.501 clause 6.3.4 based on SUCI. **If 5G-GUTI** (5G Globally Unique Temporary Identity) is provided, there is no need to map SUCI to SUPI and **steps 5-9 can be skipped**.

5. AMF sends an authentication request to the AUSF in the form of, Nausf_UEAuthentication_Authenticate. It contains the SUCI of the FN-RG. It also contains an indication that the W-5GAN has authenticated the FN-RG.

6. AUSF selects a UDM as described in clause 6.3.8 of TS 23.501 and sends a Nudm_UEAuthentication_Get Request to the UDM. It contains the SUCI of the FN-RG and indication that the W-5GAN has authenticated the FN-RG.

7. UDM invokes the SIDF to map the SUCI to a SUPI.
   
8. UDM sends a Nudm_UEAuthentication_Get Response to the AUSF. It contains the SUPI corresponding to the SUCI. It also contains an indication that authentication is not required for the FN-RG.

9. AUSF sends a Nausf_UEAuthentication_Authenticate Response to the AMF. This response from AUSF indicates that authentication is successful. The response contains the SUPI corresponding to the SUCI.
   
   The procedure described in TS 23.502 [3] clause 4.2.2.2.3 may apply (the AMF decides if the Registration Request needs to be rerouted, where the initial AMF refers to the AMF)

10. 
	1. AMF initiates a NAS security mode command procedure upon successful authentication as defined in TS 33.501
	   
	   The NAS security mode command is sent from the AMF to the W-AGF in a N2 Downlink NAS transport message
	   
	2. W-AGF responds to the AMF with a NAS Security Mode Complete message in a N2 Uplink NAS transport message. A NAS security context is created between W-AGF and AMF

11. The AMF performs steps 11-16 in TS 23.502 clause 4.2.2.2.2.
    The AMF may be configured by local policies to issue EIR check:
    - Only if the PEI is an IMEI; or
    - Only if the PEI is an IMEI or a user device trusted MAC address.
      
    At FN-RG registration to UDM, **the Access Type non-3GPP access is used**. The UDM, based on Access and Mobility Subscription information authorizes the FN-RG to access the 5GC. For FN-CRG, the AMF compares the list of serving area restrictions it receives from the UDM against the ULI from the W-AGF to check if the location information is allowed for the FN-CRG, as defined in clause 9.5.1. The AMF may also interact with the PCF for obtaining the Access and Mobility policy for the FN-RG.

12. 
	1. Upon receiving NAS Security Mode Complete, the AMF shall send an N2 Initial Context Setup Request message as defined in TS 38.413 and TS 29.413 including possibly as additional W-AGF specific parameter the RG Level Wireline Access Characteristics to the W-AGF.
	   
	2. W-AGF notifies to the AMF that the FN-RG context was created by sending a N2 Initial Context Setup Response

13. **The AMF sends the N2 Downlink NAS transport with NAS Registration Accept message** (5GS registration result, 5G-GUTI, Equivalent PLMNs or SNPNs, Non-3GPP TAI, Allowed NSSAI, Rejected NSSAI, Configured NSSAI, 5GS network feature support, network slicing indication, Non-3GPP de-registration timer value, Emergency number lists, SOR transport container, NSSAI inclusion mode) **to the W-AGF**. 

14. The W-AGF sends a N2 Uplink NAS transport message, including a NAS Registration Complete message, back to the AMF when the procedure is completed. The W-AGF shall store the 5G-GUTI to be able to send it in potential later NAS procedures.

15. The AMF performs step 23-24 in TS 23.502 clause 4.2.2.2.2.