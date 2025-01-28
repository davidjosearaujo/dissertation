![[Pasted image 20250128122327.png]]Network Slicing is considered as one of the key feature by 3GPP in 5G.  A network slice can be looked as a logical end-to-end network that can be dynamically created. A UE may access to multiple slices over the same gNB. Each slice may serve a particular service type with agreed upon Service-level Agreement (SLA).

Data between UE and Network (or another UE) go through various components on the data path. In most case, the resource allocation and the data path is configured statically or semi-statically. We cannot say all those components are optimized for each individual user or each individual use case (e.g, rush time traffic, regular hour traffic, eMBB, IoT etc). However, in ideal case where you can configure resource allocation and parameters of the components along the data path dynamically (by automation), we may define a set of parameters of all the components on the data path in most optimal way for specific UEs or specific use cases. The specific set of parameters assigned for the UEs or use cases is called a 'Slice' of the network. Network Slice is a logical concept of splitting all the resources along the data path into multiple sets, each of which is optimized for specific UEs or use cases.

# Who decides which slide to use ?
Of course, the final decision is done by network since network functions as a master in almost every communication, but the network slice selection can be triggered by either UE or Network or Both.  It depends on a variety of factors such as the type of service or application, the location and capabilities of the mobile phone, and the network conditions.

**The mobile phone will typically have a UICC that stores the credentials of the slice it is authorized to use.** When the phone connects to the network, it will use the credentials to authenticate with the network and request access to a specific slice.

**The network can also decide which slice to assign to the mobile phone based on the type of service or application the phone is trying to access.** For example, if the phone is trying to access a high-bandwidth video streaming service, the network may assign it to a slice that has a higher capacity for handling video traffic.

Additionally, Network can also use the mobile phone's location, device capabilities, and other information to decide which slice to assign to it. For example, if the mobile phone is located in an area with a high density of users, the network may assign it to a slice with a lower capacity to ensure fair usage of resources.
![[Pasted image 20250128152608.png]]
# Identification of Network Slice
To **uniquely identify a Network Slice**, the 5G system **defines the S-NSSAI** (Single – Network Slice Selection Assistance Information).

**S-NSSAI is made up of two field SST (Slice/Service Type) and SD (Service Differentiator). SD is an optional field**. SST has 8 bit field length implying that it can indicates a total of 255 different slice types.
![[Pasted image 20250128142956.png]]
1. _**SST (Slice/Service Type) –**_ this will define the expected behavior of the Network Slice in terms of specific features and services, such as [V2X](https://www.5gworldpro.com/blog/2022/02/24/audi-brings-5g-connectivity-to-its-vehicles-starting-from-2024/). The standard SST values are highlighted below :

| <br>Slice/Service Type | SST Value | Characteristics                                                                  |
| ---------------------- | --------- | -------------------------------------------------------------------------------- |
| eMBB                   | 1         | Slice suitable for the handling of 5G enhanced Mobile Broadband                  |
| URRLC                  | 2         | Slice suitable for the handling of ultra-reliable and low latency communications |
| MIoT                   | 3         | Slice suitable for the handling of massive IoT                                   |
| V2X                    | 4         | Slice suitable for the handling of V2X services                                  |
| HMTC                   | 5         | Slice suitable for the handling of High-Performance Machine-Type Communications  |
2. **SD (Slice Differentiator)** – this is optional information that complements the Slice/Service type and is used as an additional differentiator if multiple Network Slices carry the same SST value

Network slicing signaling process happenes at a few different stages like Initial Attach, PDU Establishment, Policy Change. Probably the most important step would be the process at Initial Attach. PDU Establishment would be mostly for defining various QoS Flow and Polich Change would be to associate the specific slices to a specific UE/Application as specified in the Policy Rules.

Following is the brief signaling flow at initial attach related to network slicing.
- When UE send Registration Request it specifies NSSAI. 3GPP defines a few large group of slice service/type as described [here](https://www.sharetechnote.com/html/5G/5G_NetworkSlicing.html#Identification_of_Network_Slice). By this NSSAI in Registration Request message UE tells Network saying that I want to get access to this and this type of slice.
- Then this requested Slice Information is tranferred to UDM. UDM check if the requested slice is allowed for the specific UE. If it is allowed, the UDP accept the request. If not, it rejects the request.
- Once the slice request is accepted, (with a few more steps of checkups on core network side) the acceptance is notified to UE with a few additional information like Configured NSSI, NSSI Inclusion Mode etc.
![[Pasted image 20250128144125.png]]
