- PDU stands for **Packet Data Unit**. PDU Session Establishment is the process of establishing a **data path between the UE and the 5G core network**.
- A PDU session **is a logical connection** between the UE and a data network, such as the internet or a private network. It is used to carry user data and can support different types of services, such as voice, video, and data.
- The **UE initiates the PDU Session Establishment process by sending a request to the 5G core network**. The request includes information about the type of service that the UE wants to use, and the type of traffic.
- Once the PDU session has been established, the UE can use it to send and receive data. **The 5G core network manages the resources used by the PDU session** to ensure that the network is used efficiently and that the UE receives the appropriate QoS.
- PDU Session Establishment is a key component of 5G networks, as it enables the efficient and secure transport of data between the UE and the network.
- **This is equivalent to PDN Setup process in LTE.** If you have a good understandins on [PDN setup process in LTE](https://www.sharetechnote.com/html/Handbook_LTE_IP_Allocation.html), it would be easy to get the picture of PDU Session Establishment.hahaha

In 5G, we use a “PDU Session” to provide end-to-end user plane connectivity between the UE and a specific Data Network (DN) through the User Plane Function (UPF). A PDU Session supports one or more QoS Flows. There is a one-to-one mapping between QoS Flow and QoS profile, i.e. all packets belonging to a specific QoS Flow have the
same “5G Quality of Service Identifier” (5Ql).

![[Pasted image 20250128101204.png]]
# PDU session type
"PDU session type" indicates the type of data traffic that will be transmitted over the PDU session between UE and the 5G Core Network (5GC). The PDU session type determines the characteristics of the PDU session, such as the IP version used for the data traffic and the type of data bearers required.
The PDU session type is negotiated between the UE and the 5GC during the PDU Session Establishment process. The UE indicates its supported PDU session types in the PDU Session Establishment Request message, and the network selects the appropriate PDU session type based on the UE's capabilities, network configuration, and specific use case. Once the PDU session type is agreed upon, the network assigns IP addresses and allocates resources accordingly, enabling the UE to send and receive data over the established PDU session.
These types can be:
- IPv4
- IPv6
- IPv4v6
- Unstructured
- Ethernet
- reserved