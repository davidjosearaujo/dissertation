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