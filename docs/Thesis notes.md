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
![[general_topology.png]]
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
# Development and Implementation
## Development Environment and Tools
To construct and validate the proposed framework, a virtualized multi-VM environment was orchestrated using Vagrant with VirtualBox as the provider This approach allowed for the creation of a reproducible and isolated network testbed. The environment consists of four distinct Virtual Machines (VMs), each running **Ubuntu 22.04 LTS (Jammy Jellyfish)** as the base operating system. The roles and typical resource allocations for these VMs, as defined in the `Vagrantfile`, are:
1. **`core` VM:** Hosts the 5G Core Network (5GC) functions and the EAP Authentication Server. Allocated 2GB RAM and 1 CPU.
2. **`gnb` VM:** Runs the 5G RAN gNodeB simulator. Allocated 1GB RAM and 1 CPU.
3. **`ue` VM:** Represents the 5G Residential Gateway (5G-RG), acting as a UE towards the 5GC and as an EAP Authenticator/Gateway towards the NAUN3 device. Allocated 1GB RAM and 1 CPU.
4. **`naun3` VM:** Simulates the Wi-Fi-only/NAUN3 end device, acting as an EAP Supplicant. Allocated 1GB RAM and 1 CPU.

The following core software components and tools were utilized across these VMs, installed and configured via shell scripts executed during Vagrant provisioning:

1. **5G Network Simulation:**
	1. **Open5GS:** The open-source implementation of 5G Core Network functions (AMF, SMF, UPF, NRF, AUSF, UDM, UDR, PCF, NSSF). Installed from the official PPA (`ppa:open5gs/latest`) on the `core` VM.
	2. **MongoDB:** Used as the database backend for Open5GS, storing subscriber information and network function configurations. Installed from the official MongoDB repositories on the `core` VM.
	3. **UERANSIM:** An open-source 5G RAN (gNB) and UE simulator. Cloned from its GitHub repository and compiled from source on the `gnb` VM (for gNB functionality) and the `ue` VM (for UE/5G-RG functionality). The `nr-cli` utility from UERANSIM was also made available.
2. **Authentication Infrastructure:**
	1. **FreeRADIUS:** Employed as the EAP-TLS Authentication Server. Installed on the `core` VM and configured to handle EAP-TLS, manage client (UE/5G-RG) definitions, and generate/use X.509 certificates.
	2. **`hostapd`:** Utilized as the EAP Authenticator on the `ue` VM (5G-RG). Cloned from its official repository (`w1.fi/hostap.git`) and compiled from source with the `CONFIG_DRIVER_WIRED=y` option enabled to support EAP over wired interfaces for the NAUN3 device connection.
	3. **`wpa_supplicant`:** Used as the EAP Supplicant on the `naun3` VM. Installed via `apt` and configured to perform EAP-TLS authentication using client certificates.
3. **Networking and Utility Tools:**
	1. **`dnsmasq`:** Configured as a DHCP server on the `ue` VM to provide IP addresses to NAUN3 devices connecting to its local network interface (`enp0s9`).
	2. **`yq`:** A command-line YAML processor, installed via `snap`. Extensively used in provisioning scripts to modify Open5GS and UERANSIM configuration files (e.g., setting IP addresses, DNNs, APNs).
	3. **Build Tools:** `make`, `git`, `gcc`, `g++`, `cmake` (via `snap`), `libsctp-dev`, `lksctp-tools`, `pkgconf`, `libssl-dev`, `libnl-3-dev`, `libnl-genl-3-dev` were installed for compiling UERANSIM and `hostapd` from source.
	4. **System Utilities:** `iproute2`, `net-tools`, `curl`, `gnupg` were used for network configuration and repository management.
	5. **Node.js and Nginx:** Installed on the `core` VM to support and expose the Open5GS WebUI.
4. **Custom Tools and Scripts:**
	1. **`open5gs-dbctl`:** A shell script provided and used on the `core` VM to interact with the MongoDB database for managing Open5GS subscriber entries (adding UEs, defining APNs and slices).
	2. **`interceptor`:** A custom Go application (compiled from source located in an `interceptor` directory, as indicated in the `Vagrantfile`) deployed on the `ue` VM. This is the key tool developed to orchestrate the logic for monitoring `hostapd` events and managing PDU sessions. It's specific internal workings are detailed later.
	3. **Provisioning Scripts:** A set of shell scripts (`core_install`, `gnb_install`, `ue_install`, `naun3_install`, `auth_server_install`, `ueransim_install`) were used by Vagrant to automate the installation and configuration of all software components on their respective VMs.
### Network Topology and Configuration Management
The `Vagrantfile` defines several private networks to interconnect the VMs, establishing distinct network segments for communication between the 5GC and gNB (`192.168.56.0/24`), gNB and UE/5G-RG (`192.168.57.0/24`), and the UE/5G-RG's local network for NAUN3 devices (`192.168.60.0/24`). IP addresses for various interfaces and services (e.g., `CORE_IP`, `GNB_IP_CORE`, `UE_LAN_IP`, `AUTH_SERVER_IP` for RADIUS communication over the `backhaul` tunnel) are explicitly defined and passed as arguments to the provisioning scripts
.
Vagrant's synced folder feature was utilized to share:
- EAP/RADIUS certificates generated by FreeRADIUS on the `core` VM to the `naun3` VM (via `/certs` on the guest).
- Runtime logs from all VMs to a `./build/runtime-logs` directory on the host machine.
- The compiled `interceptor` binary to the `ue` VM.
- 
This comprehensive setup provides a fully functional, albeit simulated, environment for developing and testing the proposed solution for integrating NAUN3 devices into a 5G network.
### Implementation of Proposed Authentication Logic
This section details the practical implementation of the EAP-TLS authentication mechanism, which forms the first crucial step in integrating NAUN3 devices. The setup involves configuring an EAP Authentication Server (FreeRADIUS), an Authenticator (hostapd on the 5G-RG), and an EAP Supplicant (wpa_supplicant on the NAUN3 device), with a custom Go application (`interceptor`) orchestrating t
#### 1. EAP Authentication Server Setup (FreeRADIUS on `core` VM
The EAP Authentication Server was implemented using FreeRADIUS on the `core` VM, provisioned by the `auth_server_install.txt` script.
- **RADIUS Client Configuration:** The 5G-RG (`ue` VM) was registered as a RADIUS client in `/etc/freeradius/3.0/clients.conf`. This entry specified the `ue` VM's IP address designated for EAP traffic (`CLIENT_EAP_IP` from `Vagrantfile`, e.g., `10.45.0.2`) and a shared secret (`CLIENT_SECRET` from `Vagrantfile`) for securing RADIUS communication.
- **Certificate Infrastructure:** A Public Key Infrastructure (PKI) was established within FreeRADIUS. The script automates the generation of:
    - A Certificate Authority (CA) certificate (`ca.pem`).
    - A server certificate (`server.pem`) and private key (`server.key`) for FreeRADIUS itself.
    - A client certificate (client.p12, a PKCS#12 bundle containing the certificate and private key) for the NAUN3 device.
        Passwords for these certificates (CERT_CA_PASSWD, CERT_SERVER_PASSWD, CERT_CLIENT_PASSWD) are generated or read by the Vagrantfile and passed to the script.
- **EAP-TLS Module Configuration:** The EAP module in FreeRADIUS (typically `/etc/freeradius/3.0/mods-available/eap`) was configured to:
    - Set `default_eap_type = tls`.
    - Specify the `private_key_password` for the server's private key. 
    - Define the paths to `private_key_file` (`/etc/freeradius/3.0/certs/server.key`), `certificate_file` (`/etc/freeradius/3.0/certs/server.pem`), and `ca_file` (`/etc/freeradius/3.0/certs/ca.pem`).
- **Certificate Distribution:** The necessary certificates for the NAUN3 device (CA certificate `ca.pem` and the client PKCS#12 bundle `client.p12`) were copied from the `core` VM's FreeRADIUS certificate directory to a Vagrant synced folder (`./build/eap-radius-certs-sync`), making them accessible to the `naun3` VM at `/certs`.
#### 2. Authenticator Implementation (hostapd on 5G-RG/`ue` VM)
The 5G-RG (`ue` VM) acts as the EAP Authenticator, using `hostapd`. Its setup is managed by the `ue_install.txt` script.
- **Installation:** `hostapd` was cloned from `w1.fi/hostap.git` and compiled from source, ensuring the `CONFIG_DRIVER_WIRED=y` option was enabled in its `.config` file to support EAP authentication over the wired Ethernet interface connecting to the NAUN3 device.
- **Configuration (`hostapd.conf`):**
    - The `interface` was set to `enp0s9` (the LAN-facing interface of the 5G-RG).
    - `driver=wired` was specified.
    - `ieee8021x=1` enabled 802.1X/EAP authentication.
    - `ctrl_interface=/var/run/hostapd` created a control socket for `hostapd_cli` and, more importantly, for the custom `interceptor` application to receive events.
    - RADIUS server parameters were configured:
        - `auth_server_addr`: Set to `AUTH_SERVER_IP` (the IP of the `core` VM, e.g., `10.45.0.1`, reachable via the `backhaul` PDU session).
        - `auth_server_port=1812`.
        - `auth_server_shared_secret`: Set to `CLIENT_SECRET`.
        - `own_ip_addr`: Set to `CLIENT_EAP_IP` (e.g., `10.45.0.2`), ensuring RADIUS packets from `hostapd` originate from the correct IP address within the `backhaul` PDU session.
#### 3. Supplicant Implementation (wpa_supplicant on `naun3` VM
The NAUN3 device was configured to act as an EAP Supplicant using `wpa_supplicant`, as detailed in `naun3_install.txt`.
- **Installation and Configuration (`wpa_supplicant.conf`):**
    - A network block was defined specifically for EAP-TLS authentication:
        - `key_mgmt=IEEE8021X`.
        - `eap=TLS`.
        - `identity="user@example.org"` was used as the EAP identity.
        - `ca_cert="/certs/ca.pem"` specified the path to the CA certificate (obtained via the synced folder).
        - `private_key="/certs/client.p12"` specified the path to the client's PKCS#12 certificate/key bundle.
        - `private_key_passwd` was set to `CERT_CLIENT_PASSWD`.
    - `wpa_supplicant` was configured to use the `wired` driver (`-Dwired`) for its network interface (`-ienp0s8`).
    - `ap_scan=0` was set, appropriate for wired connections.
#### 4. Orchestration of Authentication Events (`interceptor.go` on 5G-RG/`ue` VM)
The custom Go application, `interceptor.go`, deployed on the `ue` VM, plays a pivotal role in bridging the EAP authentication outcome with the 5G session management logic.
- **Hostapd Interaction:** The `interceptor` establishes a connection to `hostapd`'s control interface socket (e.g., `/var/run/hostapd/enp0s9`) using a Unix domain socket. It sends an "ATTACH" command to subscribe to `hostapd` events.
- **Event Listening:** The `HostapdListener` goroutine continuously monitors messages from `hostapd`.
- **EAP Success Detection:** It specifically parses incoming messages for the string `"CTRL-EVENT-EAP-SUCCESS"`.
- **MAC Address Extraction:** Upon detecting a successful authentication, the `interceptor` extracts the MAC address of the authenticated NAUN3 device from the event message (e.g., `CTRL-EVENT-EAP-SUCCESS aa:bb:cc:dd:ee:ff`).
- **Trigger for Next Steps:** This successful authentication event, along with the identified MAC address, serves as the primary trigger. The `interceptor` then proceeds to manage the device in its `allowed_devices` map and initiates the logic for establishing a new PDU session for this device (detailed in the subsequent section on Identity Management Implementation).
#### 5. Implemented EAP-TLS Authentication Flow Su
The implemented authentication sequence is as follows:
1. The `naun3` VM (`wpa_supplicant`) attempts to authenticate over its `enp0s8` interface connected to the `ue` VM's (`5G-RG`) `enp0s9` interface.
2. `hostapd` on the `ue` VM detects the EAPOL-Start and initiates the EAP-TLS exchange.
3. EAP-TLS messages are relayed by `hostapd` to the FreeRADIUS server on the `core` VM. This communication occurs via RADIUS packets, which are transported over the 5G-RG's `backhaul` PDU session (using `CLIENT_EAP_IP` as the source and `AUTH_SERVER_IP` as the destination).
4. FreeRADIUS validates the NAUN3 device's client certificate against its CA and configuration.
5. Upon successful validation, FreeRADIUS sends a RADIUS Access-Accept containing an EAP-Success message back to `hostapd`.
6. `hostapd` relays the EAP-Success to `wpa_supplicant` on the NAUN3 device and concurrently emits the `CTRL-EVENT-EAP-SUCCESS <NAUN3_MAC_ADDRESS>` message to its control interface.
7. The `interceptor` application on the `ue` VM captures this event, confirming the NAUN3 device's successful local authentication and readiness for the next stage of network integration.    
This setup effectively implements the EAP-TLS authentication flow, using standard tools configured to interact in a specific way, with the custom `interceptor` acting as the crucial link to the 5G-specific actions.
### Implementation of Identity Management Mechanisms
With local EAP-TLS authentication successfully implemented, this section details how the framework manages a unique network presence for each authenticated NAUN3 device. Given that these devices lack native 5G identifiers (SUPI/SUCI), the core of the implemented solution is the dynamic establishment and management of a dedicated PDU Session by the 5G-RG for each NAUN3 device. This PDU Session effectively serves as its proxy identity within the 5G network. The custom Go application, whose entry point is `main.go` [cite: main.go] and which utilizes several handler modules, running on the 5G-RG (`ue` VM), is the central orchestrator of this mechanism.

**1. Triggering Proxy Identity Establishment:** The process of establishing a proxy identity for an NAUN3 device is initiated immediately after its successful local authentication. The `HostapdListener` goroutine (defined in `hostapd_interceptor.go` [cite: hostapd_interceptor.go] and launched by `main.go` [cite: main.go]) monitors `hostapd`'s control interface. Upon detecting the `CTRL-EVENT-EAP-SUCCESS <MAC_ADDRESS>` message (constant `hostapdEventEAPSuccess`), the listener extracts the MAC address of the successfully authenticated NAUN3 device. This event and the device's MAC address serve as the trigger for the subsequent identity management steps.

**2. PDU Session Creation by the Orchestration Logic:** Once an NAUN3 device is authenticated, the `HostapdListener` calls the `NewPDUSession` function (from `ueransim_pdu_handler.go` [cite: ueransim_pdu_handler.go]). This function is responsible for requesting a new PDU session from the 5GC via the 5G-RG's UE stack (UERANSIM).
- It constructs and executes the UERANSIM command-line interface tool: `nr-cli <5G-RG_IMSI> --exec "ps-establish IPv4 --sst 1 --dnn <DNN_NAME>"`.
    - The `<5G-RG_IMSI>` (e.g., `999700000000001`) is the pre-configured IMSI of the 5G-RG itself, passed as a command-line argument (`--imsi`) to the main application (`main.go` [cite: main.go]) and subsequently to the `HostapdListener` and `NewPDUSession`.
    - The `<DNN_NAME>` (e.g., `clients`) is also passed as a command-line argument (`--dnn`) [cite: main.go] and used to target the PDU session request. This DNN was specifically configured in Open5GS on the `core` VM (`core_install.txt` [cite: core_install.txt]) and made known to the UERANSIM UE stack on the `ue` VM (`ue_install.txt` [cite: ue_install.txt]).
- After requesting the session, `NewPDUSession` enters a polling loop, repeatedly calling `LastPDUSession` (which executes `nr-cli <5G-RG_IMSI> --exec "ps-list"` and parses its YAML output) until the newly requested session transitions to the "PS-ACTIVE" state (`pduSessionStateActive`) and has an IP address assigned by the 5GC. A timeout mechanism with retries (`pduSessionEstablishRetries`, `pduSessionEstablishInterval`) is implemented [cite: ueransim_pdu_handler.go].

**3. Gateway-Managed Internal Mapping:** The main application (`main.go` [cite: main.go]) maintains a global map: `allowedDevices := make(map[string]Device)`. The `Device` struct (defined in `network_handler.go` [cite: network_handler.go]) stores the state.
- When an NAUN3 device successfully authenticates and its dedicated PDU session (type `Session` from `ueransim_pdu_handler.go` [cite: ueransim_pdu_handler.go]) becomes active, an entry is added to this map by the `HostapdListener`. The MAC address of the NAUN3 device serves as the key.
- The `Device` struct stores crucial information: its current `state` (e.g., "AUTHENTICATED", "LEASED", "REACHABLE"), a pointer to the `Session` struct (containing PDU Session ID, state, APN/DNN, and the 5GC-allocated IP Address), `Lease` information (from `dnsmasq_handler.go` [cite: dnsmasq_handler.go]), and `AppliedIPTablesRules` (a slice of `AppliedRuleDetail` from `routing_handler.go` [cite: routing_handler.go]). This map is central to linking the local device to its 5G network representation and its specific traffic routing rules.

**4. NAUN3 Device Local IP Addressing:** For the NAUN3 device to communicate on the local network segment:
- Upon successful PDU session establishment, the `HostapdListener` calls `AllowMAC()` (from `dnsmasq_handler.go` [cite: dnsmasq_handler.go]). This function appends `dhcp-host=<NAUN3_MAC_ADDRESS>,<LEASE_TIME>,set:known` to `/etc/allowed-macs.conf` on the 5G-RG (`ue` VM). The `leaseTime` is passed as a command-line argument (`--lease-time`) to `main.go` [cite: main.go].
- `dnsmasq` on the 5G-RG (configured in `ue_install.txt` [cite: ue_install.txt]) serves DHCP only to MAC addresses listed as "known" in this file.
- `wpa_supplicant` on the NAUN3 VM then executes `sudo dhclient enp0s8` to obtain a local IP address (e.g., from `192.168.60.0/24`) [cite: naun3_install.txt].

**5. Proxy Identity Lifecycle Management (Termination):** The `HostDisconnectListener` goroutine (in `network_handler.go` [cite: network_handler.go], launched by `main.go` [cite: main.go]) and the `ForgetDevice` function (in `network_handler.go` [cite: network_handler.go]) manage the termination of the proxy identity.
- `HostDisconnectListener` periodically checks device reachability on the LAN interface (e.g., `enp0s9`, derived from the `--interface` flag passed to `main.go`) using `netlink.NeighList`. A device is considered for removal if its state becomes stale/failed and its DHCP lease is significantly into its expiry, or if the MAC is no longer in the ARP/neighbor list.
- This triggers the `ForgetDevice` function, which performs cleanup:
    - `DisallowMAC()` (from `dnsmasq_handler.go` [cite: dnsmasq_handler.go]): Removes the NAUN3's MAC from `/etc/allowed-macs.conf` and the leases file, followed by `RestartDnsmasq()`.
    - `Deauth()` (from `hostapd_interceptor.go` [cite: hostapd_interceptor.go]): Sends a "DEAUTHENTICATE <MAC_ADDRESS>" command to `hostapd`.
    - `ruleManager.RemoveRulesForDevice()` (from `routing_handler.go` [cite: routing_handler.go]): Removes the specific `iptables` and `ip route/rule` entries associated with the device, using the stored `AppliedIPTablesRules`.
    - `ReleasePDUSession()` (from `ueransim_pdu_handler.go` [cite: ueransim_pdu_handler.go]): Executes `nr-cli <5G-RG_IMSI> --exec "ps-release <PDU_SESSION_ID>"` to terminate the dedicated `clients` PDU session.
    - The device's entry is removed from the global `allowedDevices` map.
- The `DnsmasqListener` (in `dnsmasq_handler.go` [cite: dnsmasq_handler.go]) monitors the DHCP lease file, updating lease details and device state to "LEASED" within the `allowedDevices` map.

**6. Implemented Traffic Mapping:** The `interceptor` system, through the `routing_handler.go` module [cite: routing_handler.go], implements concrete traffic mapping for each authenticated NAUN3 device. This is no longer conceptual.
- After a PDU session is established for an NAUN3 device, the `HostapdListener` calls `ruleManager.ApplyMappingRules()`. The `ruleManager` is initialized in `main.go` [cite: main.go].
- `ApplyMappingRules` takes the LAN interface (e.g., `enp0s9`), NAUN3's MAC address, the PDU session's interface name (e.g., `uesimtun<ID-1>`), the PDU session's gateway IP (e.g., `10.46.0.1` for the `clients` DNN, passed as `--pdu-gw-ip` flag [cite: main.go]), and the PDU session ID.
- It then systematically configures policy-based routing:
    1. **Custom Routing Table:** An entry for a new routing table (e.g., `200+<PDU_ID> table_pdu_<PDU_ID>`) is added to `/etc/iproute2/rt_tables` using `manageRTTableEntry`.
    2. **Default Route in Custom Table:** `ip route add default via <PDU_GATEWAY_IP> dev <PDU_IF_NAME> table table_pdu_<PDU_ID>` directs traffic for this table out through the NAUN3's PDU session interface.
    3. **Policy Rule:** `ip rule add fwmark <PDU_ID> table table_pdu_<PDU_ID>` directs packets marked with the PDU session ID to use this custom routing table.
    4. **Packet Marking:** An `iptables` rule in the `mangle` table's `PREROUTING` chain (`-i <LAN_IF> -m mac --mac-source <NAUN3_MAC> -j MARK --set-mark <PDU_ID>`) marks incoming packets from the NAUN3 device.
    5. **Forwarding:** An `iptables` rule in the `filter` table's `FORWARD` chain allows marked packets from the NAUN3's MAC on the LAN interface to be forwarded to its PDU session interface (`-i <LAN_IF> -o <PDU_IF_NAME> -m mac --mac-source <NAUN3_MAC> -m mark --mark <PDU_ID> -j ACCEPT`).
    6. **NAT:** An `iptables` rule in the `nat` table's `POSTROUTING` chain (`-o <PDU_IF_NAME> -j MASQUERADE`) performs SNAT for traffic exiting via the PDU session interface.  
- These rules ensure that traffic from a specific NAUN3 device is marked, routed through its dedicated PDU session, and NATted appropriately, effectively isolating its traffic and using its PDU session IP for external communication. The applied rules are stored in the `Device` struct and removed by `RemoveRulesForDevice` upon termination..
