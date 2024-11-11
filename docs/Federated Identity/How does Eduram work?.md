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
![[Pasted image 20241111154114.png]]
802.1X authentication involves three parties: a **supplicant**, an **authenticator**, and an **authentication server**.
- The supplicant is a client device (such as a laptop) that wishes to attach to the LAN/WLAN;
- The authenticator is a network device that provides a data link between the client and the network and can allow or block network traffic between the two, such as an Ethernet switch or wireless access point;
- The authentication server is typically a trusted server that can receive and respond to requests for network access, and can tell the authenticator if the connection is to be allowed, and various settings that should apply to that client's connection or setting.
## Typical authentication progression
![[Pasted image 20241111175529.png]]
The typical authentication procedure consists of:
1. **Initialization**: On detection of a new supplicant, the port on the switch (authenticator) is enabled and set to the "unauthorized" state. In this state, only 802.1X traffic is allowed; other traffic, such as the Internet Protocol (and with that TCP and UDP), is dropped.
2. **Initiation**: To initiate authentication the authenticator will periodically transmit EAP-Request Identity frames to a special Layer 2 MAC address (01:80:C2:00:00:03) on the local network segment. The supplicant listens at this address, and on receipt of the EAP-Request Identity frame, it responds with an EAP-Response Identity frame containing an identifier for the supplicant such as a User ID. The authenticator then encapsulates this Identity response in a RADIUS Access-Request packet and forwards it on to the authentication server.
3. **Negotiation** _(Technically EAP negotiation)_
# RADIUS