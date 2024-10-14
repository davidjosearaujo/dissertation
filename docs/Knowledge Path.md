# What is the problem we are trying to solve?
>[[5G and Wi-Fi RAN Convergence#3.6 Support for Wi-Fi Only Devices|Current 3GPP standard does not define architecture to support Wi-Fi only devices without USIM connecting to 5G Core.]] *in [[5G and Wi-Fi RAN Convergence.pdf|5G and Wi-Fi RAN Convergence]]* 

> [[WP_Heading-to-a-successful-private-digital-convergence.pdf|It can be observed the lack of standards to support pure WLAN devices connections to the 5GC. This is identified in the 5G work group by the Wireless Broadband Alliance (WBA), stating that “most Wi-Fi-only devices, e.g., devices in enterprise deployments, would not have a USIM included” recommending that “3GPP needs to define architecture and procedures for supporting Wi-Fi only UE with non-IMSI based identity and EAP-TLS/EAP-TTLS based authentication”. Altice Labs is aware of this limitation and is working towards an interim, proprietary solution, while this is not addressed by 3GPP.]]
# 5G Overview
See [[Wireline Access in 5G - MPIRICAL|this]]
## How is it different from 4G?

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