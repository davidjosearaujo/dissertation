# [What is the Extensible Authentication Protocol?](https://www.techtarget.com/searchsecurity/definition/Extensible-Authentication-Protocol-EAP)
EAP is used on encrypted networks to provide a secure way to send identifying information to provide network authentication. It supports various authentication methods, including as token cards, smart cards, certificates, one-time passwords and public key encryption.

EAP uses the 802.1x standard as its authentication mechanism over a local area network or a wireless LAN (WLAN). There are three primary components of 802.1X authentication:
- the user's wireless device;
- the wireless access point (AP) or authenticator; and
- the authentication database or the authentication server.

**The organization or user must choose what type of EAP to use based on their requirements. EAP transfers authentication information between the user and authenticator database or server.**

The EAP process works as follows:
1. A user requests connection to a wireless network through an AP.
2. The AP requests identification data from the user and transmits that data to an authentication server. ***(What can be identity of a legacy IOT device in 5GC?)***
3. The authentication server asks the AP for proof of the validity of the identification information.
4. The AP obtains verification from the user and sends it back to the authentication server.
5. The user is connected to the network as requested.
![[Pasted image 20241030152521.png]]
# Tunneled EAP methods
There are upwards of 40 EAP methods, including several commonly used ones that are often called inner methods or tunneled EAP methods. These include the following.
## EAP-TLS (Transport Layer Security)
EAP-TLS provides certificate-based, mutual authentication of the network and the client. **Both the client and the server must have certificates to perform this authentication.** EAP-TLS randomly generates session-based, user-based Wired Equivalent Privacy (WEP) keys. These keys secure communications between the AP and the WLAN client.

One **disadvantage** of EAP-TLS is the **server and client side both must manage the certificates**. This can be challenging for organizations with an extensive WLAN.
## EAP-TTLS (Tunneled TLS)
Like EAP-TLS, EAP-TTLS offers an extended security method with certificate-based mutual authentication. **However, instead of both the client and the server requiring a certificate, only the server side does.** EAP-TTLS enables WLANs to securely reuse legacy user authentication databases, such as Active Directory.