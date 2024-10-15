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
## SEAF
The **Security Anchor Function (SEAF)** is in a serving network and is a **“middleman” during the authentication** process between a UE and its home network. It can reject an authentication from the UE, but it relies on the UE’s home network to accept the authentication.
## AUSF
The **Authentication Server Function (AUSF)** is in a home network and **performs authentication with a UE**. It makes the decision on UE authentication, but it **relies on backend service for** computing the authentication data and keying materials when **5G-AKA or EAP-AKA’** is used.
## UDM
**Unified data management (UDM)** is an entity that hosts functions related to data management, such as the **Authentication Credential Repository and Processing Function (ARPF)**, which **selects an authentication method based on subscriber identity and configured policy** and computes the authentication data and keying materials for the AUSF if needed.
## SIDF
The **Subscription Identifier De-concealing Function (SIDF) decrypts a Subscription Concealed Identifier (SUCI)** to obtain its long-term identity, namely the Subscription Permanent Identifier (SUPI), e.g., the IMSI. In 5G, a subscriber long-term identity is always transmitted over the radio interfaces in an encrypted form. More specifically, a **public key-based encryption is used to protect the SUPI**. Therefore, only the SIDF has access to the private key associated with a public key distributed to UEs for encrypting their SUPIs.
# 5G Authentication Framework
A unified authentication framework has been defined to make 5G authentication both open (e.g., with the support of EAP) and access-network agnostic (e.g., supporting both 3GGP access networks and non-3GPP access networks such as Wi-Fi and cable networks).

![[Pasted image 20241014203612.png]]

**When EAP** (Extensible Authentication Protocol) **is used** (e.g., EAP-AKA’ or EAP-TLS), EAP **authentication is between the UE** (an EAP peer) **and the AUSF** (an EAP server) **through the SEAF** (functioning as an EAP pass-through authenticator).

**When authentication is over untrusted, non-3GPP access networks**, a new entity, namely the **Non-3GPP Interworking Function (N3IWF), is required** to function **as a VPN server** to allow the UE to access the 5G core over untrusted, non-3GPP networks through IPsec (IP Security) tunnels.
# 5G-AKA
In 5G-AKA, the SEAF may start the authentication procedure after receiving any signaling message from the UE. Note that the **UE should send the SEAF a temporary identifier** (a 5G-GUTI) or an encrypted permanent identifier (a SUCI) if a 5G-GUTI has not been allocated by the serving network for the UE. **The SUCI is the encrypted form of the SUPI using the public key of the home network**. Thus, a UE’s permanent identifier, e.g., the IMSI, is never sent in clear text over the radio networks in 5G. This feature is considered a major security improvement over prior generations such as 4G.

The SEAF starts authentication by sending an authentication request to the AUSF, which first verifies that the serving network requesting the authentication service is authorized. Upon success, the AUSF sends an authentication request to UDM/ARPF. If a SUCI is provided by the AUSF, then the SIDF will be invoked to decrypt the SUCI to obtain the SUPI, which is further used to select the authentication method configured for the subscriber. In this case, it is 5G-AKA, which is selected and to be executed.

UDM/ARPF starts 5G-AKA by sending the authentication response to the AUSF with an authentication vector consisting of an AUTH token, an XRES token, the key KAUSF, and the SUPI if applicable (e.g., when a SUCI is included in the corresponding authentication request), among other data.

![[Pasted image 20241015111210.png]]

The AUSF computes a hash of the expected response token (HXRES), stores the KAUSF, and sends the authentication response to the SEAF, along with the AUTH token and the HXRES. Note that the SUPI is not sent to the SEAF in this authentication response. It is only sent to the SEAF after UE authentication succeeds.

The SEAF stores the HXRES and sends the AUTH token in an authentication request to the UE. **The UE validates the AUTH token by using the secret key it shares with the home network**. If validation succeeds, the UE considers the network to be authenticated. The UE continues the authentication by computing and sending the SEAF a RES token, which is validated by the SEAF. Upon success, the RES token is further sent by the SEAF to the AUSF for validation.  **Note that the AUSF, which is in a home network, makes the final decision on authentication.** If the RES token from the UE is valid, the AUSF computes an anchor key (KSEAF) and sends it to the SEAF, along with the SUPI if applicable. The AUSF also informs UDM/ARPF of the authentication results so they can log the events, e.g., for the purpose of auditing.

Upon receiving the KSEAF, the **SEAF derives the AMF key (KAMF)** (and then deletes the KSEAF immediately) **and sends the KAMF to the co-located Access and Mobility Management Function (AMF)**. The AMF will then derive from the KAMF:
- (a) - the confidentiality and integrity keys needed to protect signaling messages between the UE and the AMF and,
- (b) - another key, KgNB, which is sent to the Next Generation NodeB (gNB) base station for deriving the keys used to protect subsequent communication between the UE and the gNB.
Note that the **UE has the long-term key**, which is the **root of the key derivation hierarchy**. Thus, the **UE can derive all above keys, resulting a shared set of keys between the UE and the network**.

5G-AKA differs from 4G EPS-AKA in primarily the following areas:
- The **UE always uses the public key of the home network to encrypt the UE permanent identity before it is sent to a 5G network**. In 4G, the UE always sends its permanent identifier in clear text to the network, allowing it to be stolen by either a malicious network (e.g., a faked base station) or a passive adversary over the radio links (if communication over radio links is not protected).
- **The home network** (e.g., the AUSF) **makes the final decision on UE authentication** in 5G. In addition, results of UE authentication are also sent to UDM to be logged. In 4G, a home network is consulted during authentication only to generate authentication vectors; it does not make decisions on the authentication results.
- **Key hierarchy is longer** in 5G than in 4G because 5G introduces two intermediate keys, KAUSF and KAMF
# EAP-AKA’
EAP-AKA’ is another authentication method supported in 5G. It is also a challenge-and-response protocol based on a cryptographic key shared between a UE and its home network. It accomplishes the same level of security properties as 5G-AKA, e.g., mutual authentication between the UE and the network. Because it is based on EAP, its message flows differ from those of 5G-AKA. Note that EAP messages are encapsulated in NAS messages between the UE and the SEAF and in 5G service messages between the SEAF and the AUSF.  Other differences between 5G-AKA and EAP-AKA’ are as follows.
- **The role of the SEAF** in authentication differs slightly. In EAP-AKA’, EAP message exchanges are between the UE and the AUSF through the SEAF, which **transparently forwards the EAP messages without being involved in any authentication decision**. In 5G-AKA, the SEAF also verifies the authentication response from the UE and may take action if the verification fails, albeit such action has not yet been defined in 3GPP TS 33.501
- Key derivation differs slightly. In 5G-AKA, the KAUSF is computed by UDM/ARPF and sent to the AUSF. In EAP-AKA’, **the AUSF derives the KAUSF itself in part based on the keying materials received from UDM/ARPF**. More specifically, the AUSF derives an Extended Master Session Key (EMSK) based on the keying materials received from UDM according to EAP and then uses the first 256 bits of the EMSK as the KAUSF.
# EAP-TLS
EAP-TLS is defined in 5G for subscriber authentication in limited use cases such as private networks and IoT environments. When selected as the authentication method by UDM/ARPF, EAP-TLS is performed between the UE and the AUSF through the SEAF, which functions as a transparent EAP authenticator by forwarding EAP-TLS messages back and forth between the UE and the AUSF. To accomplish mutual authentication, both the UE and the AUSF can verify each other’s certificate or a pre-shared key (PSK) if it has been established in a prior Transport Layer Security (TLS) handshaking or out of band. At the end of EAP-TLS, an EMSK is derived, and the first 256 bits of the EMSK is used as the KAUSF. As in 5G-AKA and EAP-AKA’, the KAUSF is used to derive the KSEAF, which is further used to derive other keying materials (see Figure 5) needed to protect communication between the UE and the network.

![[Pasted image 20241015161404.png]]