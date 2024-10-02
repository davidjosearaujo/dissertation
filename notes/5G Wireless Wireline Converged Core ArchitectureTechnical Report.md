# Abbreviations
| Abbreviation | Meaning                          |
| ------------ | -------------------------------- |
| CM           | Cable modem                      |
| CMTS         | Cable modem termination system   |
| CPE          | Customer premises equipment      |
| CRG          | Cable Residential Gateway        |
| HFC          | hybrid fiber-coax                |
| W-AGF        | Wireline Access Gateway Function |

# 5 Converged Architecture with the 3GPP 5G Core
3GPP has identified reference architecture diagrams for the interworking and integration models of convergence; they are shown in Figure 1 and Figure 2, respectively. The following **hybrid fiber-coax (HFC)** network components are portrayed in the 3GPP reference architectures. ***We are only interested in two***
- The W-5GCAN (Wireline 5G Cable Access Network) combines the HFC infrastructure, primarily the **cable model termination system (CMTS)**, with the **W-AGF interworking function**. The W-5GCAN may also include cable modem (CM) initialization servers, PacketCable Multimedia (PCMM) interfaces, and IP address management components.
- The **W-AGF (Wireline Access Gateway Function)**, as identified by 3GPP, is a layer of **interworking capabilities between the HFC network and the 5G mobile core** infrastructure. It is contained within the W-5GCAN.

## 5.1 3GPP R16 Interworking Model of Convergence

![[Pasted image 20241002154111.png]]

The interworking model for convergence as depicted in Figure 1 places **interworking and translation functions between the 5G core (5GC) and the HFC network within network infrastructure**. There is **no impact to deployed CPE and no change to CM authentication and network admission**. This method provides a means for operators to immediately realize benefits from a shared core while using legacy CMs.

The **N1 reference point supports UE authentication and network admission signaling**, with the 3GPP Non-Access Stratum (NAS) protocol profiled for fixed CPE

The **N2** reference point to the AMF carries access network **control messaging** as specified by 3GPP. This **control messaging is translated by the interfaces between the W-AGF and CMTS**

The **W-AGF acts as a 5G UE on behalf of the CM CPE** in the interworking model. It **manages registration into the 5GC**, data **session management**, and **slice selection** on behalf of the CPE.

## 5.3 Bridged CRG with Convergence of Services and Policy
![[Pasted image 20241002170115.png]]