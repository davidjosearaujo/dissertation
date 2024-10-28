# 3.1 Security Architecture and Procedures for 5G System
3GPP TS 33.501 is the key document providing a detailed description of ‘security architecture and procedures for 5G system’.

The specification defines a model of a security architecture, consisting of six security domains.
![[2024-10-22_11-38.png]]
- **Network access security (I)** – security features that enable a user terminal to authenticate and access the network by providing protection on the radio interfaces.
- **Network domain security (II)** - security features that enable network nodes to exchange signalling and user data securely.
- **User domain security (III)** - security features that enable the secure user access to mobile devices.
- **Application domain security (IV)** - security features that enable user and provider domain applications to exchange messages securely. 33.501 specifications do not cover application domain security.
- **Service Based Architecture (SBA) domain security (V)** - a new set of security features that enable network functions of the SBA to communicate securely within serving and other network domains.
- **Visibility and configurability of security (VI)** - security features that enable the user to be informed regarding which security features are in operation or not.

The acronyms used on the above image are as follows: ME=Mobile Equipment, SN=Serving Network, HE=Home environment
# 3.2.6 Authentication Framework
## 3.2.6.2 Secondary Authentication
EAP supports both primary (typically implemented during initial registration for example when a device is turned on for the first time) and secondary (executed for authorization during the set-up of user plane connections, for example, to surf the web or to establish a call) authentication. **The secondary authentication allows the operator to delegate the authorization to a third party.** It is meant for authentication between UE and external data networks (DN), residing outside the operator’s domain.

n the 3GPP security specifications (TS 33.501), provisions for EAP based secondary authentication by an external DN authentication server are specified in the clause 11.1, within the section 11, defining general security procedures between UE and external data networks. Secondary authentication as defined in this clause is optional to use.

