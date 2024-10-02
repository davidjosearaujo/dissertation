Private **cellular networks** are not new, but, as shown in Figure 1, from a current estimate of **about one thousand private networks based on 3rd Generation Partnership Project (3GPP)** specified technology, their numbers are expected to grow to 10’s of thousands due to:
- Guaranteed local coverage, enabled by the availability of localized private, unlicensed/ shared spectrum and cost-efficient, cloud-based 4G/5G core deployments;
- Reduced total cost of ownership (TCO) from the elimination of wired and other connectivity (e.g., Wi-Fi);
- Growing demand for enterprise information and data security, associated with localized data processing capabilities for ultra-high-performance applications;
- Search for increased productivity solutions through automation and digitalization of enterprise processes;
- Demand for low latency and reliable services.

![[Pasted image 20241002103905.png]]

# Private 5G networks
In the technical specification related with service requirements for the 5G system, 3GPP describes **non-public networks (NPN)**, 3GPP’s terminology for private networks, as being "*intended for the sole use of a private entity
such as an enterprise, and may be deployed in a variety of configurations, utilizing both virtual and physical elements. Specifically, they may be deployed as completely standalone networks, they may be hosted by a **public land mobile network (PLMN)**, or they may be offered as a slice of a PLMN*"

That way, a 5G NPN consists in the usage of a 5G system for private use, being deployed as:
- a **standalone NPN (SNPN**): operated by an NPN operator and not relying on network functions provided by a PLMN, or
- a **public network-integrated NPN (PNI-NPN)**: a non-public network deployed with the support of a PLMN.

As any 5G system, 5G NPN are composed of **user equipments (UE)** - terminals or end-systems -, 5G **accesses** - consisting of **next-generation NodeB (gNB)** units connecting UE via a 5G new radio (5G-NR) wireless interface, **and a 5G core (5GC)**.

# Convergence and 5G
With convergence, private industrial networking continues being a heterogeneous environment, but with common management and operation of all accesses as a single network, via a common 5G control plane and traffic aggregation entities. The 5GC control and data planes have the capability to serve other access technologies. In the scope of private deployments, potential 5GC shared services include:

**Common and consistent authentication/registration** and global assignment of security policies;
- **Unique IP** address management;
- **Consistent traffic management** (e.g., routing, forwarding, inspection, policy enforcement, QoS handling, and reporting) across all access types;
- Transversal slicing/virtual networking management;
- Exposure to external entities as a single network.

Security mechanisms for authentication and data encryption are key aspects in convergence since they must be present whenever a terminal attaches to the network. However, they are deeply dependent on the nature of the used access network. A unified authentication framework was defined for 5G, where 5G authentication and key agreement (5G-AKA) and extensible authentication protocol (EAP) for the 3rd generation authentication and key agreement (EAP-AKA’) are mandatory 5G primary authentication methods.

Besides gNB and next-generation e-NodeB (ng-eNB), for native 5G-NR and LTE accesses, respectively, the following four additional 5G access node types exist:

- **Non-3GPP interworking function (N3IWF)**, which allows 5G capable terminals, supporting non-access stratum (NAS) to connect from untrusted WLAN or other accesses deployed by third-party entities, out of the scope of 5G network owner control.
- Trusted non-3GPP gateway function (TNGF) and trusted WLAN interworking function (TWIF), aimed for trusted non-3GPP and WLAN accesses, but requiring the UE to have 3GPPcredentials and, for the first case, to support NAS. They are based on the tight coupling between a trusted access point and a gateway or interworking function.
- **Wireline access gateway function (W-AGF)**, which **connects a wireline 5G access network (W-5GAN) to the 5GC network**. It is similar to the TNGF for 5G residential gateways (5G-RG) and the TWIF for fixed-network residential gateways (FN-RG) but considering the specific characteristics of fixed access networks. **5G-RG units support NAS signaling and authenticate themselves**, while FN-RG do not support 5G capabilities and do not have 3GPP credentials in this specific context.

From the previous, it can be observed the **lack of standards to support pure WLAN devices connections to the 5GC**. This is identified in the 5G work group by the Wireless Broadband Alliance (WBA), stating that “most Wi-Fi-only devices, e.g., devices in enterprise deployments, would not have a USIM included,” recommending that “3GPP needs to define architecture and procedures for supporting Wi-Fi only UE with non-IMSI based identity and EAP-TLS/EAP-TTLS based authentication”. Altice Labs is aware of this limitation and is working towards an interim, proprietary solution, while this is not addressed by 3GPP.