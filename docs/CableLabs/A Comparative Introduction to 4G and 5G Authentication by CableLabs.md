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
- The **Security Anchor Function (SEAF)** is in a serving network and is a **“middleman” during the authentication** process between a UE and its home network. It can reject an authentication from the UE, but it relies on the UE’s home network to accept the authentication.
- The A**uthentication Server Function (AUSF)** is in a home network and **performs authentication with a UE**. It makes the decision on UE authentication, but it **relies on backend service for** computing the authentication data and keying materials when **5G-AKA or EAP-AKA’** is used.
- **Unified data management (UDM)** is an entity that hosts functions related to data management, such as the **Authentication Credential Repository and Processing Function (ARPF)**, which **selects an authentication method based on subscriber identity and configured policy** and computes the authentication data and keying materials for the AUSF if needed.
- The **Subscription Identifier De-concealing Function (SIDF) decrypts a Subscription Concealed Identifier (SUCI)** to obtain its long-term identity, namely the Subscription Permanent Identifier (SUPI), e.g., the IMSI. In 5G, a subscriber long-term identity is always transmitted over the radio interfaces in an encrypted form. More specifically, a **public key-based encryption is used to protect the SUPI**. Therefore, only the SIDF has access to the private key associated with a public key distributed to UEs for encrypting their SUPIs.
# 5G Authentication Framework
A unified authentication framework has been defined to make 5G authentication both open (e.g., with the support of EAP) and access-network agnostic (e.g., supporting both 3GGP access networks and non-3GPP access networks such as Wi-Fi and cable networks).

![[Pasted image 20241014203612.png]]

**When EAP** (Extensible Authentication Protocol) **is used** (e.g., EAP-AKA’ or EAP-TLS), EAP **authentication is between the UE** (an EAP peer) **and the AUSF** (an EAP server) **through the SEAF** (functioning as an EAP pass-through authenticator).

**When authentication is over untrusted, non-3GPP access networks**, a new entity, namely the **Non-3GPP Interworking Function (N3IWF), is required** to function **as a VPN server** to allow the UE to access the 5G core over untrusted, non-3GPP networks through IPsec (IP Security) tunnels.
# 5G-AKA
https://www.cablelabs.com/insights/a-comparative-introduction-to-4g-and-5g-authentication