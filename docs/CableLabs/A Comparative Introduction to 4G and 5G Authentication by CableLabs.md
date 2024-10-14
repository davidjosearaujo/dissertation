# [Introduction](https://www.cablelabs.com/insights/a-comparative-introduction-to-4g-and-5g-authentication)
3GPP defines an Authentication and Key Agreement (AKA) protocol and procedures that support entity authentication, message integrity, and message confidentiality, among other security properties.

The **3GPP AKA protocol is a challenge-and-response authentication protocol based on a symmetric key** shared between a subscriber and a home network. After the mutual authentication between a subscriber and a home network, cryptographic keying materials are derived to protect subsequent communication between a subscriber and a serving network, including both signaling messages and user plane data (e.g., over radio channels).
# 4G Authentication
From an authentication perspective, a cellular network consists of three main components: **UEs**, a **serving network (SN)**, and a **home network (HN)**
![[Pasted image 20241014163827.png]]
Each UE has a universal integrated circuit card (UICC) hosting at least a universal subscriber identity module (USIM) application, which **stores a cryptographic key that is shared with the subscriber’s home network**.

