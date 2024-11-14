# [What is the Extensible Authentication Protocol?](https://www.techtarget.com/searchsecurity/definition/Extensible-Authentication-Protocol-EAP)
EAP is used on encrypted networks to provide a secure way to send identifying information to provide network authentication. It supports various authentication methods, including as token cards, smart cards, certificates, one-time passwords and public key encryption.

EAP uses the 802.1x standard as its authentication mechanism over a local area network or a wireless LAN (WLAN). There are three primary components of 802.1X authentication:
- the user's wireless device;
- the wireless access point (AP) or authenticator; and
- the authentication database or the authentication server.

**The organization or user must choose what type of EAP to use based on their requirements. EAP transfers authentication information between the user and authenticator database or server.**

802.1X authentication involves three parties: a **supplicant**, an **authenticator**, and an **authentication server**.
- The **supplicant** is a client device (such as a laptop) that wishes to attach to the LAN/WLAN;
- The **authenticator** is a network device that provides a data link between the client and the network and can allow or block network traffic between the two, such as an Ethernet switch or wireless access point;
- The **authentication server** is typically a trusted server that can receive and respond to requests for network access, and can tell the authenticator if the connection is to be allowed, and various settings that should apply to that client's connection or setting.
![[Pasted image 20241111154114.png]]
## Typical authentication progression
![[Pasted image 20241111175529.png]]
The typical authentication procedure consists of:
1. **Initialization**: On detection of a new supplicant, the port on the switch (authenticator) is enabled and set to the "unauthorized" state. In this state, only 802.1X traffic is allowed; other traffic, such as the Internet Protocol (and with that TCP and UDP), is dropped.
2. **Initiation**: To initiate authentication the authenticator will periodically transmit EAP-Request Identity frames to a special Layer 2 MAC address (01:80:C2:00:00:03) on the local network segment. The supplicant listens at this address, and on receipt of the EAP-Request Identity frame, it responds with an EAP-Response Identity frame containing an identifier for the supplicant such as a User ID. The authenticator then encapsulates this Identity response in a RADIUS Access-Request packet and forwards it on to the authentication server.
3. **Negotiation** _(Technically EAP negotiation)_
4. **Authentication**: If the authentication server and supplicant agree on an EAP Method, EAP Requests and Responses are sent between the supplicant and the authentication server.
# Tunneled EAP methods
There are upwards of 40 EAP methods, including several commonly used ones that are often called inner methods or tunneled EAP methods. These include the following.
## EAP-TLS (Transport Layer Security)
EAP-TLS provides certificate-based, mutual authentication of the network and the client. **Both the client and the server must have certificates to perform this authentication.** EAP-TLS randomly generates session-based, user-based Wired Equivalent Privacy (WEP) keys. These keys secure communications between the AP and the WLAN client.

One **disadvantage** of EAP-TLS is the **server and client side both must manage the certificates**. This can be challenging for organizations with an extensive WLAN.
## EAP-TTLS (Tunneled TLS)
Like EAP-TLS, EAP-TTLS offers an extended security method with certificate-based mutual authentication. **However, instead of both the client and the server requiring a certificate, only the server side does.** EAP-TTLS enables WLANs to securely reuse legacy user authentication databases, such as Active Directory.