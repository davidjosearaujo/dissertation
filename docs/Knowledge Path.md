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

## What is trusted and untrusted 3GPP?

## What types of devices can connect to 5G?

## 5GC vs. N5GC device. What is the difference?

## How does each device connect to a 5G network?