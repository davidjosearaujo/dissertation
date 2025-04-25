 [Gemini conversation](https://gemini.google.com/app/75915c139d9aa2a2)
# Methodology and Proposed Framework
## Overall Research Approach
This work took a constructivist route to come up with an innovative solution to bring Wi-Fi-only devices under \ac{5G} networks through the use of existing \ac{5G} system facilities and features, predominantly through repurposing.

The first step involved examination of the \ac{5G} architecture and pertinent \ac{3GPP} standards. Although such standards prescribe mechanisms for non-\ac{3GPP} access along with the handling of devices behind residential gateways, it was realized that they cannot address in total the said specific problem of integrating devices that are not native to \ac{5G} credentials in terms of, e.g., \ac{USIM}, nor can they carry out standard authentication procedures of \ac{5G}. The root problem realized was how to recognize \ac{ 5GC} acknowledgment, authenticate indirectly, and manage individual Wi-Fi-only devices in such a way that no modifications are needed to their built-in facilities or the core facilities of \ac{5GC}. Ideas such as \ac{CGIDs} [cite: 38, 39] along with the involved \ac{PDU} session separation were tapped for inspiration, although such detailed implementation prescriptions are missing from such scenarios in the standards.

The primary necessity to minimize modifications to both the core facilities of \ac{5GC} along with those of the end-user Wi-Fi equipment influenced this work to be centered upon bridging in the network gateway (such as \ac{RG} or similar function), all needed smarts along with adaptation logic. In this way, the solution can be managed by network operators and is transparent to the core network as well as to those of the end devices.

The main methodology involved developing a framework in which the gateway plays the acting mediator role. This framework makes use of local authentication methods (that is, \ac{EAP-TLS}, with the gateway serving as an authenticator forwarding requests to a compulsory, network-operator-controlled, external \ac{EAP} server. This network operator-controlled server is necessary, as it is used to provide the means of authenticating these devices which have no native \ac{5G} identities) to authenticate Wi-Fi devices that wish to connect. Following successful local authentication, the gateway triggers the creation of an exclusive \ac{PDU} Session in the \ac{5GC} specific to that device. The per-device-specific \ac{PDU} Session essentially acts as a 'proxy identity' in the system-wide view of the \ac{5GC}, enabling traffic flows associated with the device to be managed by the \ac{5GC} without requiring direct knowledge of non-\ac{5G} credentials. The gateway takes over responsibility for forwarding traffic from the locally authenticated Wi-Fi device to its associated PDU Session.

This conceptual framework, defined to satisfy the analytic-defined requirements with minimal disruption, was then elaborated in detail, prototyped in a test environment, and then verified through functional and security testing, as detailed in subsequent chapters. This solution was selected for its likely ability to deliver a workable, minimally invasive solution to the defined integration gap.
## Requirements Analysis
Based on the identified problem of integrating Wi-Fi-only devices lacking native 5G credentials, and aiming for a practical and minimally disruptive solution, the following key requirements were established to guide the design of the proposed framework:
1. **Minimal Impact on Existing Infrastructure:**
    - **Core Network:** Modifications to standard 5G Core Network Functions (specifically Access and Mobility Management Function (AMF), Session Management Function (SMF), and User Plane Function (UPF)) were to be strictly limited to configuration changes (e.g., IP interface binding, Data Network Name (DNN) definitions). No code-level alterations to these core components were permissible.
    - **End Devices:** The solution must support standard Wi-Fi-only devices without requiring any specialized software, hardware modifications, or complex configurations beyond standard Wi-Fi connectivity and EAP supplicant capabilities.
    - **Radio Access Network (RAN):** Standard 5G RAN components (gNBs) should operate without modification, interacting with the core network via standard interfaces.
2. **Gateway-Centric Intelligence:**
    - To meet the minimal impact requirements above, the core adaptation logic must be concentrated within the network gateway function (e.g., 5G Residential Gateway (5G-RG) or equivalent). This gateway is responsible for mediating between the Wi-Fi device and the 5GC.
3. **Functional Requirements:**
    - **Secure Device Onboarding:** A robust mechanism for authenticating Wi-Fi devices locally before granting network access was required. This involved leveraging EAP-TLS, relayed by the gateway to an external, operator-controlled EAP server, which manages the credentials for these devices.
    - **Individual Device Representation:** Each successfully authenticated Wi-Fi device must be uniquely represented within the 5GC. This led to the requirement of establishing a dedicated PDU Session per device to act as its proxy identity.
    - **Traffic Separation:** A clear separation between internal network service traffic (e.g., RADIUS communication between gateway and EAP server) and end-user device traffic within the 5G transport network was necessary for security and management.
    - **Lifecycle Management:** The gateway must manage the complete lifecycle for each connected device, including handling initial authentication, triggering PDU Session establishment, managing potential re-authentications or disconnections, and ensuring corresponding PDU Session termination.
    - **Traffic Mapping and Isolation:** _(Note: Implementation pending)_ The gateway must reliably map upstream and downstream traffic between a specific Wi-Fi device and its dedicated PDU Session, ensuring traffic isolation between different devices connected through the same gateway.
4. **Operational Requirements:**
    - **Transparency:** The underlying mechanism (local authentication, PDU session mapping) should be transparent. The 5GC should primarily see standard PDU Session management procedures initiated by the gateway. The Wi-Fi device should only perceive a standard Wi-Fi connection and EAP authentication process.
    - **Operator Manageability:** The solution, including the gateway logic and the associated EAP infrastructure, must be deployable and manageable by the network operator.
These requirements collectively define the criteria for a successful solution capable of integrating unmodified Wi-Fi-only devices into a 5G network ecosystem with minimal disruption to existing components and processes.
## Proposed Identity Management Solution
### The Identity Management Challenge
The fundamental challenge in integrating Wi-Fi-only or NAUN3 devices into the 5G ecosystem lies in identity management. Standard 5G identification relies on the SUPI, typically derived from credentials stored securely on a USIM, such as an IMSI or a NAI. For transmission over the air, the SUPI is concealed within a SUCI. Devices lacking a USIM and the associated 5G credentials cannot generate a SUPI or SUCI, rendering them incapable of direct identification, authentication, and management by standard 5GC procedures involving the UDM and AUSF.
### Inspiration and Core Concept: PDU Session as Proxy Identity
The proposed solution leverages the inherent capabilities of the 5G-RG and the flexibility of 5G session management. A 5G-RG, from the 5GC's perspective, functions essentially as a UE, possessing its own USIM, credentials (SUPI/SUCI), and the ability to establish and maintain multiple concurrent PDU Sessions. This capability allows a single UE (the 5G-RG) to segregate its traffic based on different service requirements or destinations.

Inspiration was drawn from the concept of CGIDs described in 3GPP specifications [cite: 38, 39], where PDU Sessions can be used to manage traffic for groups of devices connected behind a gateway, potentially mapping groups to specific network interfaces (e.g., SSIDs, Ethernet ports). However, the specifications often imply a group-based segregation. We observed that since the 5G-RG incorporates routing logic and can dynamically request PDU Sessions, there was no fundamental limitation preventing a more granular approach.

Therefore, the core concept of this proposed solution is to utilize a dedicated PDU Session, established and managed by the 5G-RG, as a one-to-one proxy identity for each individually authenticated NAUN3 device. Instead of grouping devices, each device is mapped to its own PDU Session. This not only provides connectivity via the 5GC (similar to the goal of CGIDs) but also allows the PDU Session itself to serve as the handle or identifier for managing that specific device's connection within the 5G system.
### Establishing the Proxy Identity
The creation of this proxy identity is triggered by the successful local authentication of an NAUN3 device via the EAP-TLS mechanism described previously. Upon receiving the EAP-Success indication for a device, the 5G-RG initiates a PDU Session Establishment procedure towards the 5GC (SMF via AMF). Crucially, this request is made using the 5G-RG's own registered 5G identity (SUPI/IMSI). The request specifies the dedicated `clients` DNN, indicating the intended purpose of this session. Once the 5GC establishes the PDU Session and assigns it resources (like an IP address) and an identifier, the 5G-RG internally associates this specific PDU Session with the locally authenticated NAUN3 device.
### Gateway's Role in Identity Mapping
The 5G-RG is central to this identity management scheme. It maintains a dynamic internal mapping table that links the local identifier of the NAUN3 device (e.g., its MAC address or the identity used in the EAP exchange) to the specific PDU Session ID assigned by the 5GC. This mapping is essential for correctly routing traffic between the device on the local network and its representation within the 5G network.
### 5GC Perspective and Management
From the 5GC's viewpoint, the process appears relatively standard. It interacts with a registered UE (the 5G-RG, identified by its SUPI) that requests the establishment and termination of multiple PDU Sessions associated with the `clients` DNN. The 5GC's Session Management Function (SMF) manages these sessions, allocating resources and applying policies (like QoS, charging rules) on a per-PDU session basis. The 5GC remains unaware of the specific local identities or the EAP-TLS authentication details of the individual NAUN3 devices connected behind the 5G-RG; it simply manages the sessions requested by the gateway.
### Dynamic Mapping Lifecycle
The gateway actively manages the lifecycle of this identity mapping:
- **Creation:** A mapping entry is created when an NAUN3 device successfully authenticates via EAP-TLS and its corresponding PDU Session is established by the 5GC.
- **Maintenance:** The gateway routes traffic for the device via the mapped PDU Session. It also periodically checks the connectivity status of the NAUN3 device on the local network (e.g., via keep-alives or monitoring association status).
- **Termination:** If the gateway detects that a device's local connection has become stale or the device explicitly disconnects, it performs two actions:
    1. It issues a deauthentication command locally (e.g., instructing `hostapd` to disassociate the device).    
    2. It initiates a PDU Session Release procedure with the 5GC to terminate the corresponding proxy identity session.    
        The internal mapping entry is then removed.
### Advantages of this Approach
This session-based proxy identity approach offers several advantages:
- **Transparency:** It is largely transparent to both the 5GC, which manages standard PDU sessions, and the NAUN3 device, which undergoes standard local network authentication.
- **Leverages Existing Mechanisms:** It builds upon the standard 5G PDU session management framework rather than requiring fundamental changes to core identity protocols.
- **Gateway-Centric Complexity:** It concentrates the necessary adaptation logic within the 5G-RG, an entity typically managed by the network operator.
- **No SUPI Required for End Devices:** It successfully integrates devices lacking 5G credentials without needing to provision them with SUPIs or USIMs.
- **Individual Device Management:** By providing a per-device PDU session, it allows for potentially granular policy application (QoS, security) at the 5GC level based on the session, indirectly controlling individual device flows.
## Framework Architecture and Integration
### Overall Architecture Overview
![[general-topology.png]]
### Component Integration and Interactions
The proposed framework integrates several distinct components, coordinating their standard functionalities to achieve the goal of connecting NAUN3 devices to the 5G network. The interactions are orchestrated primarily by the 5G-RG, as illustrated in Figure [X.Z - Replace with actual figure number from the architecture diagram]:

- **NAUN3 Device:** Its interaction is confined to the local network segment managed by the 5G-RG. It initiates connection via standard link-layer protocols (Wi-Fi/Ethernet) and participates in the EAP-TLS authentication process as a supplicant, communicating only with the 5G-RG's authenticator function. It remains unaware of the 5GC, PDU Sessions, or the underlying mechanisms used for its connectivity beyond the local authentication step.
- **5G-RG (Gateway):** This component acts as the central integration point and performs multiple roles simultaneously:
    - **Towards 5GC:** It registers and authenticates with the 5GC using its own 5G credentials (SUPI/IMSI) like a standard UE, utilizing the N1 interface for NAS signaling with the AMF. It establishes an initial PDU session (linked to the `backhaul` DNN) via the SMF for its own operational traffic, including communication with the EAP Server.
    - **Towards NAUN3 Device:** It functions as a Layer 2 access point and an EAP authenticator (relay), managing the local connection and the EAPoL exchange.
    - **Towards EAP Server:** It securely relays EAP messages encapsulated within RADIUS packets to the external EAP Authentication Server. This communication is tunneled through its established `backhaul` PDU session.
    - **Orchestration Logic:** Upon receiving an EAP-Success via RADIUS for an NAUN3 device, its internal logic triggers a request to the SMF (via AMF) for a _new_ PDU session establishment, specifying the `clients` DNN. It then uses the internal mapping table to associate this new PDU session with the authenticated NAUN3 device and routes the device's user plane traffic accordingly between the local interface and the N3 user plane tunnel (GTP-U) established for that specific `clients` PDU session towards the UPF.
- **EAP Authentication Server:** Its interaction is solely with the 5G-RG via the RADIUS protocol. It receives authentication requests, performs the EAP-TLS server-side operations, validates client certificates, and returns an Access-Accept (with EAP-Success) or Access-Reject (with EAP-Failure) message. All this communication is secured and transported over the 5G-RG's `backhaul` PDU session.
- **5GC Network Functions (AMF, SMF, UPF):** These components perform their standard functions, treating the 5G-RG as the registered UE.
    - **AMF:** Manages the registration, authentication (of the 5G-RG itself), and mobility of the 5G-RG. It forwards session management requests from the 5G-RG to the appropriate SMF.
    - **SMF:** Handles all PDU session management procedures initiated by the 5G-RG for both the `backhaul` and the multiple `clients` DNNs. It selects the UPF, allocates IP addresses for each PDU session, and interacts with the UPF to establish the necessary user plane tunnels.    
    - **UPF:** Establishes GTP-U tunnels as instructed by the SMF for each PDU session and forwards user plane traffic between the 5G RAN (via N3 interface towards the 5G-RG) and the respective Data Network (DN) associated with the DNN.    
This integration ensures that while the NAUN3 device only undergoes local authentication, its traffic is securely tunneled through the 5GC via a dedicated, dynamically established PDU session managed by the mediating 5G-RG.
### Key Communication Flows in the Integrated System

!! CREATE A FLOW CHART !!

Illustrate the end-to-end flow for onboarding an NAUN3 device:
- Initial 5G-RG registration with the 5GC.
- Establishment of the `backhaul` PDU session for the 5G-RG.
- NAUN3 device connects locally -> EAP-TLS authentication flow occurs (Device <-> RG <-> EAP Server via `backhaul` PDU session).
- EAP-Success -> RG requests a new PDU session establishment for the `clients` DNN via AMF/SMF.
- SMF interacts with UPF to set up the user plane path for the new `clients` PDU session.
- 5G-RG receives confirmation, completes local setup, and starts mapping device traffic to the new UPF tunnel.
### Interface and Protocol Integration
A key aspect of this framework's design is its reliance on standard, well-defined interfaces and protocols, minimizing the need for proprietary extensions. The integration is achieved by orchestrating these standard elements in a specific manner, primarily through the logic implemented within the 5G-RG.

The main interfaces and protocols involved are:
- **Local Network Interface (NAUN3 Device <-> 5G-RG):**
    - **Link Layer:** Standard Ethernet (IEEE 802.3) or Wi-Fi (IEEE 802.11).    
    - **Authentication:** Extensible Authentication Protocol over LAN (EAPoL - IEEE 802.1X) is used to transport EAP messages over the local link.    
    - **EAP Method:** EAP-TLS is used for mutual authentication between the NAUN3 device (supplicant) and the EAP infrastructure (via the 5G-RG relay).    
- **5G Interfaces (5G-RG <-> 5GC):**
    - **N1 Interface:** Carries Non-Access Stratum (NAS) signaling between the 5G-RG (acting as UE) and the AMF for registration, authentication (of the RG itself), and session management procedures.
    - **N2 Interface:** Carries Next Generation Application Protocol (NGAP) signaling between the 5G RAN (gNB, which connects the 5G-RG) and the AMF, primarily for UE context management and PDU session resource setup requests related to the 5G-RG.
    - **N3 Interface:** Carries the user plane traffic encapsulated in GPRS Tunneling Protocol - User Plane (GTP-U) tunnels between the 5G RAN (gNB) and the UPF. This includes traffic for both the `backhaul` PDU session and all the individual `clients` PDU sessions.    
- **Authentication Interface (5G-RG <-> EAP Server):**
    - **Application Layer:** Remote Authentication Dial-In User Service (RADIUS) protocol is used to carry EAP messages between the 5G-RG (acting as a RADIUS client and EAP relay) and the external EAP Authentication Server (RADIUS server).    
    - **Transport:** RADIUS messages are transported over IP, typically using UDP. This IP traffic is securely tunneled through the 5G-RG's dedicated `backhaul` PDU session via the N3 interface and UPF.    

The framework integrates these components by ensuring the 5G-RG correctly handles protocol relay (EAPoL to RADIUS/EAP) and coordinates actions based on protocol outcomes (EAP-Success triggering NAS session management requests). The novelty lies not in modifying these protocols but in configuring the system components (5GC NFs, 5G-RG, EAP Server) and implementing the orchestration logic within the 5G-RG to manage the per-device proxy identity using standard 5G session management procedures.
### Summary of Integration
Conclude by summarizing how the architecture effectively integrates unmodified NAUN3 devices by leveraging the 5G-RG as a mediating entity, utilizing the PDU session framework for proxy identification and traffic management, and interfacing with standard 5GC and authentication server components.