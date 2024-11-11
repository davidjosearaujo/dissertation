*eduroam* is based on 802.1X* and a linked hierarchy of RADIUS servers containing users’ data (usernames and passwords).

Participating institutions **must have operating RADIUS infrastructure** and agree to the terms of use. *eduroam* can be set up in three easy steps:
1. Set up a RADIUS server connected to your institutional identity server (LDAP).
2. Connect your access points to your RADIUS server.
3. Federate your RADIUS server.
![[Pasted image 20241111152310.png]]
In the context of eduroam, the institutions' RADIUS servers act as Identity Providers (IdPs). Each participating institution is responsible for authenticating its own users through their RADIUS server. The federation allows these institutions to trust each other’s authentication processes, enabling seamless roaming.

The RADIUS hierarchy forwards user credentials securely to the users’ home institutions, where they are verified and validated.

To protect the privacy of the traffic from the user’s device over the wireless network, the latest up-to-date data encryption standards are used.