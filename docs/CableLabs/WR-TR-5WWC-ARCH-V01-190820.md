# Abbreviations
| Abbreviation | Meaning                                         |
| ------------ | ----------------------------------------------- |
| CM           | Cable modem                                     |
| CMTS         | Cable modem termination system                  |
| CPE          | Customer premises equipment                     |
| CRG          | Cable Residential Gateway                       |
| DOCSIS       | Data Over Cable Service Interface Specification |
| eMBB         | Enhanced Mobile Broadband                       |
| HFC          | hybrid fiber-coax                               |
| NAS          | Non-Access Stratum                              |
| PDU          | Packer data unit                                |
| W-AGF        | Wireline Access Gateway Function                |
> **Data Over Cable Service Interface Specification (DOCSIS)** is an international  telecommunications standard that permits the addition of high-bandwidth data transfer to an existing cable television (CATV) system. It is used by many cable television operators to provide cable Internet access over their existing hybrid fiber-coaxial (HFC) infrastructure.

# 5 Converged Architecture with the 3GPP 5G Core
3GPP has identified reference architecture diagrams for the interworking and integration models of convergence; they are shown in Figure 1 and Figure 2, respectively. The following **hybrid fiber-coax (HFC)** network components are portrayed in the 3GPP reference architectures. ***We are only interested in two***
- The **W-5GCAN (Wireline 5G Cable Access Network)** combines the HFC infrastructure, primarily the **cable model termination system (CMTS)**, with the **W-AGF interworking function**. The W-5GCAN may also include cable modem (CM) initialization servers, PacketCable Multimedia (PCMM) interfaces, and IP address management components.
- The **W-AGF (Wireline Access Gateway Function)**, as identified by 3GPP, is a layer of **interworking capabilities between the HFC network and the 5G mobile core** infrastructure. It is contained within the W-5GCAN.

## 5.1 3GPP R16 Interworking Model of Convergence

![[Pasted image 20241002154111.png]]

The interworking model for convergence as depicted in Figure 1 places **interworking and translation functions between the 5G core (5GC) and the HFC network within network infrastructure**. There is **no impact to deployed CPE and no change to CM authentication and network admission**. This method provides a means for operators to immediately realize benefits from a shared core while using legacy CMs.

The **N1 reference point supports UE authentication and network admission signaling**, with the 3GPP Non-Access Stratum (NAS) protocol profiled for fixed CPE

The **N2** reference point to the AMF carries access network **control messaging** as specified by 3GPP. This **control messaging is translated by the interfaces between the W-AGF and CMTS**

The **W-AGF acts as a 5G UE on behalf of the CM CPE** in the interworking model. It **manages registration into the 5GC**, data **session management**, and **slice selection** on behalf of the CPE.

# 7 Requirements
## 7.4 High Level W-AGF Requirements
Each residential gateway may be assigned its own VLAN between the CMTS and the W-AGF to complete the user plane traffic flow between the HFC access network and the W-AGF. The W-AGF must support individual VLANs per residential gateway as received via the CMTS.

## 7.5 Registration, Authentication, and CPE Status
For the FN-CRG, the **W-AGF must detect CM registration into the CMTS**. The W-AGF must report CM registrations and de-registrations.

The W-AGF must be able to **detect CM unreachability by the CMTS** and report the idle status to the 5GC

**The method by which the W-AGF detects CM registration and reachability status from the  CMTS is per vendor implementation.** ^c138a2

When the EAP-5G registration is complete and a security association is established between the5G-CRG and the W-AGF, then the **5G-CRG must support NAS over TLS to the W-AGF for the balance of registration**.

## 7.6 Slicing
The **W-AGF must support network slicing** and NSSAI on behalf of individual FN-CRGs. The W-AGF **must allow the operator to configure the way in which the W-AGF selects slices** for the FN-CRG based on slice IDs received during FN-CRG registrations, DOCSIS flows initiated by the FN-CRG, or FN-CRG application.

Upon FN-CRG registration, when no other information is available for network slice selection, the **W-AGF must select either the enhanced mobile broadband (eMBB) network slice specified by 3GPP for the 5GC to connect to default DOCSIS service flows or an operator-provisioned default slice in the W-AGF for the FN-CRG**. ^0fbf60

The W-AGF must support operator configuration of slice types to FN-CRGs and DOCSIS service flow settings.

## 7.7 Session Management
The W-AGF must support session management procedures for each FN-CRG.

The W-AGF must establish an IP PDU data session within the eMBB slice per operator configuration upon the CMTS completing default DOCSIS service flows for the FN-CRG.

If the operator has configured the FN-CRG to be a bridged gateway, then the W-AGF must establish an Ethernet PDU data session within the eMBB slice.

The W-AGF must be able to support multiple PDU sessions for the FN-CRG and to map an individual PDU session to a specific DOCSIS service flow.

