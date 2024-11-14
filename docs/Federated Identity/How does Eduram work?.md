*eduroam* is based on 802.1X* and a linked hierarchy of RADIUS servers containing users’ data (usernames and passwords).

Participating institutions **must have operating RADIUS infrastructure** and agree to the terms of use. *eduroam* can be set up in three easy steps:
1. Set up a RADIUS server connected to your institutional identity server (LDAP).
2. Connect your access points to your RADIUS server.
3. Federate your RADIUS server.
![[Pasted image 20241111152310.png]]
In the context of eduroam, the institutions' RADIUS servers act as Identity Providers (IdPs). Each participating institution is responsible for authenticating its own users through their RADIUS server. The federation allows these institutions to trust each other’s authentication processes, enabling seamless roaming.

The RADIUS hierarchy forwards user credentials securely to the users’ home institutions, where they are verified and validated.

To protect the privacy of the traffic from the user’s device over the wireless network, the latest up-to-date data encryption standards are used.

The user’s home institution is responsible for maintaining and monitoring user information, even when the user is at a guest campus. Thus, this data is not shared with other connected institutions.
# IEEE 802.1X
IEEE 802.1X defines the encapsulation of the Extensible Authentication Protocol (EAP) over wired IEEE 802 networks  and over 802.11 wireless networks, which is known as "EAP over LAN" or EAPOL.


## Typical authentication progression
![[Pasted image 20241111175529.png]]
The typical authentication procedure consists of:
1. **Initialization**: On detection of a new supplicant, the port on the switch (authenticator) is enabled and set to the "unauthorized" state. In this state, only 802.1X traffic is allowed; other traffic, such as the Internet Protocol (and with that TCP and UDP), is dropped.
2. **Initiation**: To initiate authentication the authenticator will periodically transmit EAP-Request Identity frames to a special Layer 2 MAC address (01:80:C2:00:00:03) on the local network segment. The supplicant listens at this address, and on receipt of the EAP-Request Identity frame, it responds with an EAP-Response Identity frame containing an identifier for the supplicant such as a User ID. The authenticator then encapsulates this Identity response in a RADIUS Access-Request packet and forwards it on to the authentication server.
3. **Negotiation** _(Technically EAP negotiation)_
4. **Authentication**: If the authentication server and supplicant agree on an EAP Method, EAP Requests and Responses are sent between the supplicant and the authentication server.
# [RADIUS](https://www.techtarget.com/searchsecurity/definition/RADIUS)
RADIUS (Remote Authentication Dial-In User Service) is a client-server protocol and software that enables remote access servers to communicate with a central server to authenticate dial-in users and authorize their access to the requested system or service.

RADIUS enables a company to maintain user profiles in a central database that all remote servers can share. Having a central database provides better security, enabling a company to set up a policy that can be applied at a single administered network point.

RADIUS was originally designed to support large numbers of users connecting remotely to internet service providers (ISPs) or corporate networks via modem pools or other point-to-point serial line links.
## How does RADIUS authentication work?
In the RADIUS protocol, remote network users connect to their networks through a network access server (NAS).

Unlike other client-server applications, where the client is often an individual user, RADIUS clients are the NAS systems used to access a network and the authentication server is the RADIUS server.
![[Pasted image 20241112104926.png]]
Types of remote user access authentication servers can include:
- **Dial-in servers**, which mediate access to corporate or ISP networks through modem pools.
- **Virtual private network servers**, which accept requests from remote users to set up secure connections to a private network.
- **Wireless access points**, which accept requests from wireless clients to connect to a network.
- **Managed network access switches** that implement the 802.1x authenticated access protocol to mediate access to networks by remote users.
![[Pasted image 20241112105316.png]]When an end user opens a connection with a remote network, the NAS initiates a RADIUS exchange with the authentication server.

When a remote user initiates a connection through a NAS, the request can include the remote user ID, password and IP address. The NAS then sends a request for authentication to the RADIUS server.
## How are RADIUS servers used?
RADIUS authenticates using two approaches:
- **Password Authentication Protocol (PAP)**. The RADIUS client forwards the remote user's user ID and password to the RADIUS authentication server. If the credentials are correct, the server authenticates the user and the RADIUS client enables the remote user to connect to the network.
- **Challenge Handshake Authentication Protocol (CHAP)**. Also known as a three-way handshake, CHAP authentication relies on the client and server using an encrypted shared secret. Compared to PAP, CHAP authentication is considered more secure because it encrypts authentication exchanges and it can be configured to do repeated mid-session authentications.

A RADIUS proxy client can be configured to forward RADIUS authentication requests to other RADIUS servers. RADIUS proxies enable centralized authentication in large or geographically dispersed networks.