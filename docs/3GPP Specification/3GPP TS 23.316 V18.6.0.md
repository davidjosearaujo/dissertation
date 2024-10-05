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

4.3.2 Connection **management*