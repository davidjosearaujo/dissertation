# [Introduction](https://www.cablelabs.com/insights/a-comparative-introduction-to-4g-and-5g-authentication)
3GPP defines an Authentication and Key Agreement (AKA) protocol and procedures that support entity authentication, message integrity, and message confidentiality, among other security properties.

The **3GPP AKA protocol is a challenge-and-response authentication protocol based on a symmetric key** shared between a subscriber and a home network. After the mutual authentication between a subscriber and a home network, cryptographic keying materials are derived to protect subsequent communication between a subscriber and a serving network, including both signaling messages and user plane data (e.g., over radio channels).
# 4G Authentication
From an authentication perspective, a cellular network consists of three main components: **UEs**, a **serving network (SN)**, and a **home network (HN)**
![[Pasted image 20241014163827.png]]
Each UE has a universal integrated circuit card (UICC) hosting at least a universal subscriber identity module (USIM) application, which **stores a cryptographic key that is shared with the subscriber’s home network**.
## 4G EPS-AKA
The EPS-AKA is triggered after the UE completes the Radio Resource Control (RRC) procedure with eNodeB and sends an Attach Request message to the MME.

![[Pasted image 20241014173803.png]]

There are **two weaknesses in 4G EPS-AKA**.
1. The UE **identity is sent over radio networks without encryption**. Although a temporary identifier (e.g., Globally Unique Temporary Identity, GUTI) may be used to hide a subscriber’s long-term identity, researchers have shown that GUTI allocation is flawed: **GUTIs are not changed as frequently as necessary**, and GUTI allocation is predictable (e.g., with fixed bytes). More importantly, the **UE’s permanent identity may be sent in clear text** in an Identity Response message when responding to an Identity Request message from a network.
2. Second, a home network provides authentication vectors (AVs) when consulted by a serving network during UE authentication, but it is not a part of the authentication decision. Such a **decision is made solely by the serving network**.
# 5G Authentication
Service-based architecture (SBA) has been proposed for the 5G core network. Accordingly, new entities and new service requests have also been defined in 5G. Some of the new entities relevant to 5G authentication are listed below.
- The Security Anchor Function (SEAF) is in a serving network and is a “middleman” during the authentication process between a UE and its home network. It can reject an authentication from the UE, but it relies on the UE’s home network to accept the authentication.