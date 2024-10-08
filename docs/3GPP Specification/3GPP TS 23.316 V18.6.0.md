# Definitions
- **Wireline access Control Plane protocol (W-CP)**: Protocol used to transport AS and NAS signalling between the 5G- RG and the W-AGF over the Y4 reference point. W-CP is specified by BBF and CableLabs. There is no assumption that W-CP refers to only a single protocol or only a specific protocol layer.
- **Wireline access User Plane protocol (W-UP)**: Protocol used to carry PDU Session user plane traffic between the 5G- RG and the W-AGF over the Y4 reference point. W-UP is specified by BBF and CableLabs. There is no assumption that W-UP refers to only a single protocol or only a specific protocol layer.

# 4.2.2 Identification and authentication
In the case of **FN-RG connected via W-5GAN**, the specification defined in [[3GPP TS 23.501 V18.7.0|TS 23.501 clause 5.2.3]] applies with the following differences:
- UE is replaced by FN-RG
- The **W-AGF provides the NAS signalling** connection to the 5GC **on behalf of the FN-RG**.
- The W-5GAN may authenticate the FN-BRG per BBF specification BBF TR-456 and WT-457. The W- 5GAN may authenticate the FN-CRG per [[WR-TR-5WWC-ARCH-V01-190820|CableLabs DOCSIS MULPI]].

# 4.2.3 Authorisation
In the case of FN-RG connected via W-5GAN, the specification defined in TS 23.501 [2] clause 5.2.4 applies with the following differences:
- UE is replaced by FN-RG
- **W-AGF performs the UE Registration procedure on behalf of the FN-RG**

# 4.3.1 Registration management
Registration management when 5G-RG or FN-RG is connected to 5GC via wireline access is described in [[3GPP TS 23.501 V18.7.0#5.5.1 Registration Management|TS 23.501 clause 5.5.1]]

# 4.3.2 Connection management
Connection management when 5G-RG or FN-RG is connected to 5GC via wireline access is described in [[3GPP TS 23.501 V18.7.0#5.5.2 Connection Management|clause 5.5.2 of TS 23.501]]

# 4.4.2 Session management for FN-RG
Session management of FN-RG follows the principle defined in TS 23.501 clause 5.6 with the following differences:
- UE is replaced by W-AGF
- FN-RG is connected to 5GC via wireline access instead of 3GPP access

# 4.10a Non-5G capable device behind 5G-CRG and FN-CRG
non-5G capable (N5GC) devices connecting via W-5GAN can be authenticated by the 5GC using EAP based authentication method(s) as defined in TS 33.501.

**Roaming is not supported for N5GC devices.**

![[2024-10-07_12-39.png]]
1. The W-AGF registers the FN-CRG to 5GC as specified in clause 7.2.1.3 or the 5G-CRG registers to 5GC as specified in [[3GPP TS 23.316 V18.6.0#7.2.1.1 5G-RG Registration via W-5GAN|clause 7.2.1.1]]
2. The CRG is configured as L2 bridge mode and forwards any L2 frame to W-AGF. 802.1x authentication may be triggered. This can be done either by N5GC device sending a EAPOL-start frame to W-AGF or W-AGF receives a frame from an unknown MAC address.
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

# 7.2.1.1 5G-RG Registration via W-5GAN