# What is the problem we are trying to solve?
>[[5G and Wi-Fi RAN Convergence#3.6 Support for Wi-Fi Only Devices|Current 3GPP standard does not define architecture to support Wi-Fi only devices without USIM connecting to 5G Core.]] *in [[5G and Wi-Fi RAN Convergence.pdf|5G and Wi-Fi RAN Convergence]]* 

> [[WP_Heading-to-a-successful-private-digital-convergence.pdf|It can be observed the lack of standards to support pure WLAN devices connections to the 5GC. This is identified in the 5G work group by the Wireless Broadband Alliance (WBA), stating that “most Wi-Fi-only devices, e.g., devices in enterprise deployments, would not have a USIM included” recommending that “3GPP needs to define architecture and procedures for supporting Wi-Fi only UE with non-IMSI based identity and EAP-TLS/EAP-TTLS based authentication”. Altice Labs is aware of this limitation and is working towards an interim, proprietary solution, while this is not addressed by 3GPP.]]
# 5G Overview
5G is based upon a Service Based Architecture incorporating **NFV** (Network Functions Virtualization) and **SDN** (Software Defined Network) technology.

**These benefits** are not solely confined to mobile networks and as such, **could also be incorporated by the fixed network or wireline service providers**.

![[Pasted image 20241014162545.png]]
The 5G-RG could support both a wireless interface (Uu), similar to a FWA deployment and the wireline interface (Y4) towards the 5G core via the W-AGF (Wireline – Access Gateway Function)

![[Pasted image 20241014163143.png]]
In addition to the 5G RG, an alternative endpoint has also been defined termed the FN-RG (Fixed Network Residential Gateway). This operates in a similar way to the 5G RG however it does not support the 3GPP N1 reference point. As such, the W-AGF would need to deal with the NAS signalling on its behalf.
## How is it different from 4G?
A cellular network consists of three main components: **UEs**, a **serving network (SN)**, and a **home network (HN)**

![[Pasted image 20241014163827.png]]

5G differs from prior generations primarily in that it will not only provide faster speed, higher bandwidth, and lower delays, but also support more use cases such as **enhanced mobile broadband (eMBB)**, **massive machine-type communications (mMTC)**, and **ultra-reliable low-latency communications (uRLLC)**.

The **3GPP AKA protocol is a challenge-and-response authentication protocol based on a symmetric key** shared between a subscriber and a home network.
### Two weakness in 4G EPS-AKA
1. [[A Comparative Introduction to 4G and 5G Authentication by CableLabs#^5dd100|The UE identity is sent over radio networks without encryption.]]
2. [[A Comparative Introduction to 4G and 5G Authentication by CableLabs#^a72d4c|The authentication decision is solely made by the serving network]] 
### [[A Comparative Introduction to 4G and 5G Authentication by CableLabs#5G Authentication|New entities and new service requests have also been defined in 5G]]
#### SEAF
The **Security Anchor Function (SEAF)** is in a serving network and is a **“middleman” during the authentication** process between a UE and its home network.
#### AUSF
The **Authentication Server Function (AUSF)** is in a home network and **performs authentication with a UE**. It **relies on backend service for** computing the authentication data.
#### UDM
**Unified data management (UDM)** is an entity that hosts functions related to data management. It **selects an authentication method based on subscriber identity and configured policy**.
#### SIDF
The **Subscription Identifier De-concealing Function (SIDF) decrypts a Subscription Concealed Identifier (SUCI)** to obtain its long-term identity, namely the Subscription Permanent Identifier (SUPI), e.g., the IMSI. A **public key-based encryption is used to protect the SUPI**.
# Private digital convergence

## Private 5G networks
3GPP describes **non-public networks (NPN)** as being "*intended for the sole use of a private entity such as an enterprise, and may be deployed in a variety of configurations, utilizing both virtual and physical element completely standalone networks, they may be hosted by a **public land mobile network (PLMN)**, or they may be offered as a slice of a PLMN*"

As any 5G system, 5G NPN are composed of:
- **User equipment (UE)**
- 5G **accesses**
	- Consisting of **next-generation NodeB (gNB)** units **connecting UE via a 5G new radio** (5G-NR) wireless interface
- **5G core (5GC)**
## Convergence and 5G
With convergence, private industrial networking continues being a heterogeneous environment, but with common management and operation of all accesses as a single network, via a common 5G control plane and traffic aggregation entities. The 5GC control and data planes have the capability to serve other access technologies. In the scope of private deployments, potential 5GC shared services include:
- **Unique IP** address management;
- **Consistent traffic management** (e.g., routing, forwarding, inspection, policy enforcement, QoS handling, and reporting) across all access types;
- Transversal slicing/virtual networking management;
- Exposure to external entities as a single network.

**Convergence in 5G** is achieved at the core **via functional entities placed at the 5G domain entrance**, which adapt access specific protocols to standard N2 interface control plane (CP) and N3 interface data plane (DP). The N1 interface, used to convey non-radio signaling between the UE and the 5GC, may not be supported by the terminal equipment, forcing the adaptation entities to handle it on behalf of the terminal

A unified authentication framework was defined for 5G, where **5G authentication and key agreement** (5G-AKA) and **extensible authentication protocol (EAP) for the 3rd generation authentication and key agreemen**t (EAP-AKA’) are mandatory 5G primary authentication methods.

That framework makes 5G-AKA procedure **suitable for both open and access-network agnostic scenarios**, relying on three authentication methods: 5G-AKA, EAP-AKA’, and EAP transport layer security (EAP-TLS)
## Additional 5G Access nodes
Besides gNB and next-generation e-NodeB (ng-eNB), for native 5G-NR and LTE accesses, respectively, the following four additional 5G access node types exist:
- **Non-3GPP interworking function (N3IWF)**, which allows 5G capable terminals, supporting non-access stratum (NAS) to connect from untrusted WLAN or other accesses deployed by third-party entities, out of the scope of 5G network owner control.
- Trusted non-3GPP gateway function (TNGF) and trusted WLAN interworking function (TWIF), aimed for trusted non-3GPP and WLAN accesses, but requiring the UE to have 3GPPcredentials and, for the first case, to support NAS. They are based on the tight coupling between a trusted access point and a gateway or interworking function.
- **Wireline access gateway function (W-AGF)**, which **connects a wireline 5G access network (W-5GAN) to the 5GC network**. It is similar to the TNGF for 5G residential gateways (5G-RG) and the TWIF for fixed-network residential gateways (FN-RG) but considering the specific characteristics of fixed access networks. **5G-RG units support NAS signaling and authenticate themselves**, while FN-RG do not support 5G capabilities and do not have 3GPP credentials in this specific context.
# Connecting to 5G
## What is 3GPP and non-3GPP?
3GPP (3rd Generation Partnership Project) refers to the standards developed for mobile networks, including 3G, 4G (LTE), and 5G. These are cellular technologies used by mobile carriers to provide network services.

**Non-3GPP** refers to other access technologies not standardized by 3GPP but still capable of integrating with 3GPP networks, such as Wi-Fi or satellite networks. These networks can provide connectivity, typically offloading traffic from 3GPP networks, but they follow different standards (e.g., IEEE for Wi-Fi).

In short:
- **3GPP**: Cellular (e.g., 4G, 5G)
- **Non-3GPP**: Other networks (e.g., Wi-Fi) integrated with cellular networks.
### Access to the 3GPP 5G Core Network (5GCN) via Non-3GPP Access Networks (N3AN)
>For an untrusted non-3GPP access network, to secure communication between the UE and the 5GCN, **a UE establishes secure connection to the 5G core network over untrusted non-3GPP access via the N3IWF**. 
### [[Access to the 3GPP 5G Core Network (5GCN) via Non-3GPP Access Networks (N3AN) - 24.502|Authentication and authorization for accessing 5GS via non-3GPP access network]]
>In order to register to the 5G core network (5GCN) via untrusted non-3GPP IP access, the UE first needs to be configured with a local IP address from the untrusted non-3GPP access network (N3AN).

>During EAP authentication, authentication and authorization for access to 5GCN is performed by exchange of EAP-5G message encapsulated in the link layer protocol between the UE and the TNAN.
## What is trusted and untrusted non-3GPP?
In 3GPP architecture, trusted and untrusted refer to the way a non-3GPP network (like Wi-Fi) connects to the mobile core network.
### Trusted 3GPP access
This is a non-3GPP network that has been verified and trusted by the mobile operator. It connects directly to the core network using secure protocols and behaves similarly to 3GPP networks. For example, a mobile operator’s managed Wi-Fi network might be treated as trusted.
### Untrusted 3GPP access
This is when a non-3GPP network (e.g., public Wi-Fi) isn’t under the mobile operator’s control or doesn’t meet their security standards. To access the core network, traffic must go through an additional security layer called the evolved Packet Data Gateway (ePDG), which provides encryption and authentication.
#### [[System architecture for the 5G System (5GS) - 23.501#General Concepts|General concepts]]
An untrusted non-3GPP access network shall be connected to the 5G Core Network via a **Non-3GPP InterWorking Function (N3IWF)**, whereas a trusted non-3GPP access network shall be connected to the 5G Core Network via a **Trusted Non-3GPP Gateway Function (TNGF)**. Both the N3IWF and the TNGF interface with the 5G Core Network **CP and UP functions via the N2 and N3 interfaces**, respectively.

A UE shall establish an IPsec tunnel with the N3IWF or with the TNGF in order to register with the 5G Core Network over non-3GPP access. Further details about the UE registration to 5G Core Network over untrusted non-3GPP access and over trusted non-3GPP access are described in clause 4.12.2 and in clause 4.12.2a of TS 23.502, respectively.
## [[Procedures for the 5G System (5GS) - 23.502#4.2.2 Registration Management procedures|Procedures for the 5G System (5GS)]]
### [[Procedures for the 5G System (5GS) - 23.502#4.12.2 Registration via Untrusted non-3GPP Access|Registration via Untrusted non-3GPP Access]]
>Specifies how a UE can register to 5GC via an untrusted non-3GPP Access Network. It is based on the Registration procedure specified in clause [[#4.2.2.2.2 General Registration|4.2.2.2.2]] and it uses a vendor-specific EAP method called "EAP-5G".

>The "EAP-5G" method is used between the UE and the N3IWF and is utilized only for encapsulating NAS messages (not for authentication).

>![[Pasted image 20241017151403.png]]
### [[Procedures for the 5G System (5GS) - 23.502#4.12a.2 Registration via Trusted non-3GPP Access|Registration via Trusted non-3GPP Access]]
>![[Pasted image 20241017151444.png]]

>In this case, the "EAP-5G" method is used between the UE and the TNGF and is utilized for encapsulating NAS messages.
### [[Procedures for the 5G System (5GS) - 23.502#4.12b Procedures for devices that do not support 5GC NAS over WLAN access|Procedures for devices that do not support 5GC NAS over WLAN access]]
>![[Pasted image 20241017151526.png]]

>**Devices that do not support 5GC NAS** signalling over WLAN access (referred to as "Non-5G-Capable over WLAN" devices, or N5CW devices for short), **may access 5GC** in a PLMN or an SNPN **via a trusted WLAN Access Network** that supports a Trusted WLAN Interworking Function (TWIF).
## [[SECURITY IN 5GSPECIFICATIONS - Controls in 3GPP Security Specifications (5G SA)|Security Architecture and Procedures for 5G System]]
UE requirements:
- **The UE shall support 5G-GUTI**
- The Home Network Public Key shall be stored in the USIM.
- ...
AMF requirements:
- The AMF shall support assigning 5G-GUTI to the UE.
- ...
### [[Security architecture and procedures for 5G system - 33.501#7.2.1 Authentication for Untrusted non-3GPP Access|Authentication for Untrusted non-3GPP Access]]
>It uses a vendor-specific EAP method called "EAP-5G", utilizing the "Expanded" EAP type and the existing 3GPP Vendor-Id, registered with IANA under the SMI Private Enterprise Code registry.
### [[Security architecture and procedures for 5G system - 33.501#7A.2.1 Authentication for Trusted non-3GPP access|Authentication for Trusted non-3GPP access]]
>This is based on the specified procedure in TS 23.502 clause 4.12a.2.2 "Registration procedure for trusted non-3GPP access". 
### [[5G Identifiers|5G Identifiers]]
SUCI and SUPI
### [[A Comparative Introduction to 4G and 5G Authentication by CableLabs|4G vs. 5G Authentication]]
# [[Extensible Authentication Protocol|Understanding EAP Framework]]
The EAP process works as follows:
1. A user requests connection to a wireless network through an AP.
2. The AP requests identification data from the user and transmits that data to an authentication server. ***(What can be identity of a legacy IOT device in 5GC?)***
3. The authentication server asks the AP for proof of the validity of the identification information.
4. The AP obtains verification from the user and sends it back to the authentication server.
5. The user is connected to the network as requested.
# Challenges
The [[Current Prototype Limitations|current prototype has a major limitation]] in terms of device **identity consistency and universality**.

We need to devise a solution that **enables a device to be identified universally in the network** and where its **identity is independent of its gateway**.

A few solution may be taken into account.
## Possible solutions
### Federated Identity
Implement a system where devices have centralized credentials managed by a federated identity provider. This would allow them to maintain a consistent identity recognized across various network access points.
