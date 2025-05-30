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

**6. Implemented Traffic Mapping:**
The `interceptor` system, through the `routing_handler.go` module [cite: routing_handler.go], implements concrete traffic mapping for each authenticated NAUN3 device. This is no longer conceptual.
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
## Adaptation of Network Functions
The successful implementation of the proposed framework did not require code-level modifications to the standard 5G Network Functions (NFs) or RAN/UE simulators. Instead, these components were adapted through specific configurations to support the gateway-centric authentication and identity management scheme for NAUN3 devices. The primary tools used for this were Open5GS for the 5GC, UERANSIM for the gNB and the 5G-RG's UE stack, all provisioned and configured via scripts within a Vagrant-managed virtualized environment.

**Open5GS (5GC) on `core` VM:** The Open5GS components (AMF, SMF, UPF, AUSF, UDM) were configured as follows:
- **Connectivity Configuration:** The AMF's NGAP interface and the UPF's GTP-U interface were bound to the `core` VM's specific IP address (`CORE_IP` variable from `Vagrantfile.txt`, e.g., `192.168.57.10`) within the private network connecting to the gNB. This was achieved by modifying `amf.yaml` and `upf.yaml` configuration files using `yq`.
- **DNN Configuration:** Two distinct DNNs were defined to segregate traffic:
    - `backhaul` DNN: Configured in `smf.yaml` and `upf.yaml`, associated with the default `ogstun` tunnel interface. This DNN provides the primary PDU session for the 5G-RG itself, used for operational traffic such as RADIUS communication with the EAP Authentication Server. The SMF configuration for this DNN typically assigns IPs from a specific subnet (e.g., `10.45.0.0/24` is a common Open5GS default, and the scripts use `AUTH_SERVER_IP="10.45.0.1"` and `CLIENT_EAP_IP="10.45.0.2"` which fall into this range).
    - `clients` DNN: Also configured in `smf.yaml` and `upf.yaml`. A dedicated tunnel interface (`clientun0`) was created on the `core` VM and assigned an IP address (e.g., `10.46.0.1/24`). The `clients` DNN was associated with this interface and subnet (`10.46.0.0/24`), enabling the SMF/UPF to allocate IPs from this range for the individual PDU sessions established for each NAUN3 device.
- **Subscriber Provisioning (for the 5G-RG):** The 5G-RG itself was provisioned as a standard UE in the Open5GS MongoDB database using the `open5gs-dbctl.sh` script. This involved adding an entry with the 5G-RG's IMSI (`UE_IMSI`), pre-shared key (`UE_KEY`), and OPC (`UE_OPC`). This subscription was configured to allow the 5G-RG to establish PDU sessions on both the `backhaul` and `clients` DNNs using the `update_apn` command of `open5gs-dbctl.sh`.
- **AUSF/UDM Operation:** These functions operate in a standard manner for the 5G-RG. The 5G-RG authenticates to the 5GC using its provisioned USIM credentials. The NAUN3 devices are not directly known to or authenticated by the 5GC's AUSF/UDM. The 5GC only interacts with the authenticated 5G-RG, which then requests PDU sessions on behalf of the NAUN3 devices.

**UERANSIM (gNB on `gnb` VM and UE stack on `ue` VM):** The UERANSIM components were configured without any code modifications, using their standard YAML configuration files manipulated by `yq`:
- **gNB Configuration (`gnb` VM):** The `open5gs-gnb.yaml` configuration file was modified to set:
    - The gNB's link IP address (`GNB_IP_UE` from `Vagrantfile.txt`, e.g., `192.168.58.10`) for communication towards the UE (5G-RG).
    - The gNB's NGAP and GTP IP addresses (`GNB_IP_CORE` from `Vagrantfile.txt`, e.g., `192.168.56.100`) for N2 and N3 interface communication with the AMF and UPF respectively.
    - The AMF's IP address (`CORE_IP` from `Vagrantfile.txt`, e.g., `192.168.57.10`) for establishing the N2 connection.
- **5G-RG UE Stack Configuration (`ue` VM):** The `open5gs-ue.yaml` file for the UERANSIM instance representing the 5G-RG was configured to:
    - Specify the gNB's IP address in its search list (`GNB_IP_UE`).
    - Use the provisioned USIM credentials (IMSI, Key, OPC) identical to those configured in the Open5GS database.
    - Support multiple PDU sessions. The `ue_install.sh` script explicitly configures the first session (index 0) to use the `backhaul` APN/DNN. It then prepares subsequent session configurations (though `DEFAULT_PDU_SESSIONS` is 0 in the `Vagrantfile.txt`, the script allows for more) to support IPv4 type PDU sessions on the `clients` APN/DNN with SST 1. This configuration enables the 5G-RG's UE stack to request the necessary PDU sessions as orchestrated by the `interceptor` application.

These configurations ensured that the standard 5G components could support the project's architecture, where the 5G-RG acts as a legitimate UE capable of managing multiple distinct data paths (via PDU Sessions on different DNNs) for its own operational needs and for proxying connectivity for the locally authenticated NAUN3 devices.
## System Integration and Configuration
The various implemented components were integrated into a cohesive test environment using Vagrant for VM orchestration and shell scripts for provisioning. This ensured consistent configurations and network connectivity across the simulated 5G system and attached devices.

**Virtual Machine Orchestration and Network Topology:** The `Vagrantfile.txt` defines four VMs (`core`, `gnb`, `ue`, `naun3`), each running Ubuntu 22.04 LTS. Vagrant establishes several private networks to facilitate communication:
- **5GC-gNB Network (`192.168.57.0/24`):** Connects the `core` VM (hosting Open5GS NFs) with the `gnb` VM. The `core` VM uses `CORE_IP` (e.g., `192.168.57.10`) and the `gnb` VM uses `GNB_IP_CORE` (e.g., `192.168.57.100`) on this network for N2 (AMF-gNB) and N3 (UPF-gNB) interface traffic.
- **gNB-UE (5G-RG) Network (`192.168.58.0/24`):** Connects the `gnb` VM with the `ue` VM (representing the 5G-RG). The `gnb` VM uses `GNB_IP_UE` (e.g., `192.168.58.10`) and the `ue` VM uses `UE_IP` (e.g., `192.168.58.100`) for the simulated Uu radio interface.
- **5G-RG LAN (`192.168.59.0/24`):** Connects the `ue` VM (interface `enp0s9` with IP `UE_LAN_IP` from `Vagrantfile.txt`, e.g., `192.168.59.10`) to the `naun3` VM (interface `enp0s8` with IP `NAUN3_IP` from `Vagrantfile.txt`, e.g., `192.168.59.100`). This segment handles local EAPoL for authentication and the NAUN3 device's data traffic before it's routed into a PDU session. Vagrant's synced folders facilitate sharing of EAP/RADIUS certificates (from `core` to `naun3` via `/certs`), runtime logs from all VMs to the host (`./build/runtime-logs`), and the compiled `interceptor` binary to the `ue` VM (`./build/interceptor` to `/home/vagrant/interceptor`).

**5G-RG (`ue` VM) as the Central Integration Hub:** The `ue` VM is the cornerstone of the integration, performing multiple functions simultaneously:
- **RAN Connectivity:** Its UERANSIM UE stack connects to the simulated gNB on the `gnb` VM.
- **Local Network Services:**
    - `hostapd` is configured on its LAN interface (`enp0s9`) to provide 802.1X/EAP-TLS authentication for devices on the `192.168.59.0/24` network, relaying authentication requests to the FreeRADIUS server.
    - `dnsmasq` acts as a DHCP server for this LAN, dynamically assigning IPs to authenticated NAUN3 devices based on the `/etc/allowed-macs.conf` file managed by the `interceptor`.
- **`backhaul` PDU Session:** Upon registration with Open5GS, the 5G-RG's UE stack establishes its primary PDU session on the `backhaul` DNN. This session receives an IP address (e.g., `CLIENT_EAP_IP = 10.45.0.2`) from the 5GC. This IP is then used as the source for RADIUS messages sent from `hostapd` (via the 5G-RG) to the FreeRADIUS server (`AUTH_SERVER_IP = 10.45.0.1`) on the `core` VM.
- **`clients` PDU Sessions:** Orchestrated by the `interceptor` application, for each successfully EAP-authenticated NAUN3 device, the 5G-RG's UE stack requests a new, dedicated PDU session on the `clients` DNN. These sessions are assigned IPs from the `10.46.0.0/24` subnet by the SMF/UPF.  
- **IP Forwarding and Traffic Routing:**
    - IP forwarding is enabled on the `ue` VM (`sudo sysctl -w net.ipv4.ip_forward=1`).
    - The `interceptor`, through its `routing_handler.go` module, dynamically configures policy-based routing using `ip rule` and `ip route` commands, and `iptables` rules for packet marking (`mangle` table), NAT (`nat` table, `MASQUERADE`), and forwarding (`filter` table). This ensures that traffic originating from an NAUN3 device on the local LAN is marked, routed through its dedicated `clients` PDU session tunnel interface (e.g., `uesimtunX`), and NATted with that PDU session's assigned 5GC IP address.

**NAUN3 Device (`naun3` VM) Configuration:** The `naun3` VM connects to the 5G-RG's LAN interface (`enp0s8` on `naun3` to `enp0s9` on `ue`). Its `wpa_supplicant` service is configured for EAP-TLS authentication using the client certificate and CA certificate shared via the Vagrant synced folder (`/certs`). Upon successful authentication, it obtains a local IP address via DHCP from the `dnsmasq` server on the 5G-RG.

**EAP Authentication Server (FreeRADIUS on `core` VM):** FreeRADIUS on the `core` VM is configured to listen for RADIUS requests. The `ue` VM (5G-RG) is defined as a RADIUS client, identified by its `backhaul` PDU session IP address (`CLIENT_EAP_IP`). It authenticates NAUN3 devices based on the client certificates issued by its internal CA.

**Parameter Consistency:** The `Vagrantfile.txt` serves as the central point for defining critical network parameters (IP addresses, IMSI, keys, shared secrets, certificate passwords). These parameters are passed as arguments to the respective provisioning scripts (`*.sh`), ensuring consistency across the configurations of Open5GS, UERANSIM, FreeRADIUS, `hostapd`, and the command-line flags for the `interceptor` application. This integrated setup allows for the end-to-end simulation and testing of the proposed NAUN3 device integration framework.
## Implementation Challenges

The development and implementation of this framework, while ultimately successful in the simulated environment, encountered several technical challenges. These ranged from complexities in orchestrating the virtualized 5G system to specific issues encountered when attempting to integrate physical hardware.

**1. Orchestration of Simulated 5G Components:**
- **Configuration Complexity:** Setting up a multi-VM environment with Open5GS, UERANSIM, FreeRADIUS, `hostapd`, and `dnsmasq` required careful management of numerous configuration files and network parameters. Ensuring IP address consistency, correct DNN definitions, subscriber provisioning, and proper inter-component communication (e.g., RADIUS, NGAP, GTP-U) across different virtual networks demanded meticulous scripting (as seen in the Vagrant provisioning scripts). Any misconfiguration in one component often had cascading effects, making debugging a time-consuming process.
- **Service Dependencies and Startup Order:** Ensuring that services started in the correct order and that dependencies were met (e.g., MongoDB before Open5GS NFs, Open5GS NFs before UERANSIM components could connect) was crucial and required careful scripting within the Vagrant provisioning process.
- **Dynamic PDU Session Management with UERANSIM:** While UERANSIM's `nr-cli` tool provides a command-line interface to manage PDU sessions, programmatically triggering and monitoring these from an external application (the custom orchestration logic) involved parsing command output and implementing polling mechanisms, which is less robust than a direct API-based interaction might be.

**2. Development of the Custom Orchestration Logic:**
- **Event Handling and State Management:** The custom orchestration application needed to reliably capture events from `hostapd` (EAP success), manage the state of multiple NAUN3 devices (authentication status, associated PDU session, applied routing rules), and react to network events (DHCP lease changes, device disconnections via ARP/neighbor cache monitoring). Coordinating these asynchronous events and maintaining a consistent internal state for each device was a key challenge.
- **Interfacing with System Utilities:** The custom logic interacts with system utilities like `nr-cli` (for PDU sessions), `iptables`, and `ip route`/`ip rule` (for traffic mapping). Ensuring these commands were executed correctly with the appropriate parameters for each device, and handling their output or potential errors, required careful implementation and robust error checking within the Go application.
- **Concurrency and Resource Management:** Managing multiple goroutines for listening to `hostapd`, `dnsmasq` leases, and network disconnects, while ensuring thread-safe access to shared data structures (like the map of allowed devices), required careful use of synchronization primitives.

**3. Challenges with Physical Modem Integration (Quectel RG500Q-GL RedCap Attempt):** An attempt was made to integrate a physical 5G modem, specifically a Quectel RG500Q-GL (RedCap) USB modem, to explore the feasibility of the solution with real hardware acting as the 5G-RG's connection to the 5G network. This presented significant challenges distinct from the simulated UERANSIM environment [cite: `Notes.md`]:
- **Proprietary Drivers and Kernel Dependencies:** The Quectel RG500Q-GL, being an experimental sample, relied on proprietary drivers provided by Quectel rather than standard Linux kernel drivers. These drivers had to be compiled from source and were highly sensitive to specific kernel versions. This severely restricted the choice of host operating system and often necessitated the use of a dedicated Single Board Computer (SBC) that met the kernel requirements, complicating the development workflow (requiring SSH access to the SBC for modem interaction).
- **Lack of Public Documentation:** Comprehensive public documentation for the modem's AT commands, QMI interface (`qmicli`), and particularly for advanced features like establishing multiple concurrent PDU sessions with QMAP (Quectel Multiplexing an APplication processor) mode, was scarce or non-existent. This made configuring the modem for the project's specific needs (e.g., one `backhaul` PDU session, multiple `clients` PDU sessions) a process of trial, error, and reliance on limited provided snippets.
- **Difficulties with Multiple PDU Sessions:**
    - While AT commands and `qmicli` could be used to define PDP contexts for different APNs/DNNs (e.g., `AT+CGDCONT` or `qmicli --wds-create-profile`), activating and managing multiple _simultaneous_ PDU sessions, especially binding them to distinct virtual network interfaces for independent routing by the custom orchestration logic, proved extremely challenging with the provided Quectel tools (`quectel-qmi-proxy`, `quectel-CM`).
    - The available examples and tools from Quectel primarily demonstrated setting up multiple connections to different APNs but did not clearly address the scenario of multiple active PDU sessions to the _same_ APN (our `clients` DNN) or robustly exposing these as distinct network interfaces to the Linux system in a way that the custom routing logic could easily manage.
    - The `qmicli` tool, while powerful, did not offer a straightforward or well-documented method for QMAP-based multiplexing of multiple PDU sessions that was confirmed to work with this specific modem model and firmware.
- **Contrast with UERANSIM:** The UERANSIM environment, by comparison, allowed for relatively straightforward programmatic control over PDU session establishment and release via `nr-cli`, making it a more tractable platform for developing and testing the core logic of the custom orchestration application and the overall framework. The complexities of the physical modem's driver and proprietary connection manager abstracted away much of the direct control needed for fine-grained, multi-session management.

These challenges highlight the gap that can exist between simulated environments and the intricacies of physical hardware, especially when dealing with proprietary drivers and limited documentation for specialized or pre-release components. While the core concepts of the project were validated in simulation, porting to such physical hardware would require significant additional effort in driver-level integration and modem-specific control.

# Validation and Results Evaluation
This chapter details the methodology employed to validate the proposed framework for integrating Wi-Fi-only/NAUN3 devices into the 5G network. It outlines the test scenarios, key performance indicators (KPIs), and the evaluation of the results obtained from the implemented simulation environment.
## Validation Methodology
To assess the feasibility, functionality, and effectiveness of the proposed solution, a series of tests were conducted within the simulated environment described in the "Development and Implementation" chapter. The overall validation approach was **simulation-based testing**, leveraging the orchestrated virtual machines and configured 5G components (Open5GS, UERANSIM) and local network services (`hostapd`, `dnsmasq`, and the custom `interceptor` application).

The validation focused on several key aspects of the system:

**I. Overall Validation Approach:**
- **Simulation-Based Testing:** All validation activities were performed within the virtualized environment created using Vagrant, Open5GS, UERANSIM, FreeRADIUS, and the custom `interceptor` application. This allowed for controlled and repeatable testing of the end-to-end solution.
- **Focus on Functional Correctness and Integration:** The primary goal was to verify that the proposed mechanisms for authentication, proxy identity creation (per-device PDU session), traffic mapping, and lifecycle management operate as designed.
- **Qualitative Security Assessment:** While not a formal security audit, the validation included observing whether the implemented security measures (EAP-TLS, traffic separation) were functioning as intended.

**II. Key Performance Indicators (KPIs) and Metrics for Evaluation:**
The evaluation of the framework centered on the following indicators and metrics, primarily assessed through functional testing and observation of system logs and behavior:

**1. Functional Correctness:**
The functional correctness of the core mechanisms was evaluated based on the following aspects: 
* **NAUN3 Device Authentication Success:** 
	* Metric: Successful completion of the EAP-TLS authentication process for an NAUN3 device with the FreeRADIUS server, relayed by the 5G-RG (`hostapd` and `interceptor`).
	* Verification: Logs from `wpa_supplicant` (NAUN3 VM), `hostapd` (5G-RG/`ue` VM), FreeRADIUS (`core` VM), and the custom `interceptor` application on the 5G-RG. 
* **Dedicated PDU Session Establishment:**
	* Metric: Successful establishment of a unique PDU session on the `clients` DNN by the 5G-RG for each successfully authenticated NAUN3 device.
	* Verification: Output of `nr-cli ps-list` command on the 5G-RG (`ue` VM); logs from Open5GS SMF and UPF on the `core` VM. 
* **IP Address Allocation:** 
	* Metric (Local): Successful assignment of a local IP address to the NAUN3 device by `dnsmasq` on the 5G-RG after EAP-TLS authentication. 
	* Metric (5GC): Successful assignment of a 5GC IP address by the 5GC (SMF/UPF) to the dedicated PDU session for the NAUN3 device. 
	* Verification: `ip addr` command output on the NAUN3 VM; `dnsmasq` logs on the 5G-RG; `nr-cli ps-list` output on the 5G-RG; Open5GS SMF/UPF logs. 
* **End-to-End Data Plane Connectivity:** 
	* Metric: Ability of an authenticated NAUN3 device to send and receive IP traffic to/from an external network via its dedicated PDU session. 
	* Verification: Ping tests and simple data transfer (e.g., HTTP GET) from the NAUN3 VM to a target beyond the UPF; packet captures (`tcpdump`) on NAUN3 LAN interface, 5G-RG's PDU session tunnel interface, and UPF interfaces. 
* **Traffic Isolation and Mapping:** 
	* Metric: Confirmation that traffic from a specific NAUN3 device is routed exclusively through its dedicated PDU session and associated routing rules. 
	* Verification: Packet captures on the 5G-RG; analysis of `iptables` counters, `ip rule` and `ip route` configurations on the 5G-RG during active traffic from one or more NAUN3 devices. 
* **Lifecycle Management Correctness:** 
	* Metric: Successful de-authentication of an NAUN3 device and termination of its associated PDU session upon simulated disconnection/unreachability. 
	* Verification: Logs from the `interceptor` application, `hostapd`, and `dnsmasq`; `nr-cli ps-list` output showing PDU session release; verification of removal of `iptables` rules and `dnsmasq` permissions.

**2. Security Aspects (Qualitative Observation):**
- **EAP-TLS Authentication Integrity:** Observation of the complete EAP-TLS handshake and successful mutual authentication through detailed logs from involved components (`wpa_supplicant`, `hostapd`, FreeRADIUS).
- **Traffic Segregation:** Confirmation via network monitoring and PDU session analysis that `backhaul` DNN traffic (e.g., RADIUS) remains logically separate from the `clients` DNN traffic (NAUN3 user plane data).
- **NAUN3 Identity Concealment from 5GC:** Verification that the NAUN3 device's local identifiers (e.g., MAC address) are not directly signaled to or stored by the core 5GC NFs (AMF, SMF, UDM), with the 5G-RG acting as the boundary.

**3. Resource Management (Observational):**
- **PDU Session Correlation:** The number of active PDU sessions on the `clients` DNN should directly correspond to the number of currently authenticated and connected NAUN3 devices, as tracked by the `interceptor` application.
- **Timeliness of Operations (Qualitative):** General observation of the time taken for the end-to-end process: NAUN3 device EAP-TLS authentication, subsequent PDU session establishment, and the teardown process upon device disconnection. Formal latency measurements were considered outside the primary scope of this functional validation.

**4. System Stability and Robustness (Qualitative):**
- **Handling Multiple Devices:** The ability of the `interceptor` application and the overall simulated system to manage sequential and concurrent connections and disconnections of multiple NAUN3 devices without instability.
- **Error Handling:** Observation of error logging and any recovery mechanisms within the `interceptor` application in scenarios such as a failed PDU session establishment attempt or unexpected disconnection.

This validation methodology aims to provide a comprehensive assessment of the implemented solution's ability to meet its design goals, focusing on correct functionality and integration within the simulated 5G environment. The subsequent sections will detail the specific test scenarios designed and the evaluation of the results obtained.
## Test Scenarios and Setup
To validate the different aspects of the proposed framework, a series of distinct test scenarios, or experiments, were designed and executed. These scenarios leveraged the fully configured simulation environment detailed in the "Development and Implementation" chapter, which includes the `core` VM (Open5GS, FreeRADIUS), `gnb` VM (UERANSIM gNB), `ue` VM (5G-RG with UERANSIM UE, `hostapd`, `dnsmasq`, and the custom `interceptor` application), and one or more `naun3` VMs (EAP supplicant).

Specific tools were employed for monitoring and verification in each scenario:
- **Log Analysis:** Reviewing logs from Open5GS NFs, UERANSIM components, FreeRADIUS, `hostapd`, `wpa_supplicant`, and the custom `interceptor` application was fundamental across all tests.
- **Packet Capture:** `tcpdump` and `tshark` (Wireshark CLI) were used on various interfaces (NAUN3 LAN, 5G-RG's `backhaul` and `clients` PDU session interfaces `uesimtunX`, gNB interfaces) to inspect signaling and data plane traffic.
- **Network Utilities:** Standard Linux utilities like `ping` (with the `-R` record route option), `ip addr`, `ip route`, `ip rule`, `iptables -L -v -n -t mangle -t nat -t filter`, and UERANSIM's `nr-cli` were used for connectivity testing and state verification.
- **Traffic Generation & Throughput Measurement:** `iperf3` was used for generating controlled TCP/UDP network traffic between NAUN3 devices and a server on the `core` VM (simulating an N6-connected server) to test data plane throughput and routing.
- **Timestamping/Scripting:** Basic shell scripting was used to capture timestamps before and after key events (e.g., `wpa_supplicant` start and `dhclient` IP acquisition) to measure onboarding delay.
### Experiment 1: Single NAUN3 Device Onboarding and Basic Connectivity
The objective was to verify the successful EAP-TLS authentication of a single NAUN3 device, the subsequent establishment of its dedicated PDU session on the `clients` DNN, local and 5GC IP address allocation, basic end-to-end data plane connectivity with path verification, and to measure the approximate onboarding delay.

Procedure:
1. Ensure all 5GC NFs, gNB, 5G-RG (including `hostapd` and `interceptor`), and FreeRADIUS are running. Start an `iperf3` server on the `core` VM listening on the IP address of its `clientun0` interface (e.g., `10.46.0.1`).
2. On the `naun3` VM, record a timestamp (`date +%s`). Immediately start the `wpa_supplicant` service.
3. Monitor the authentication process through logs on the `naun3` VM, `ue` VM (`hostapd`, `interceptor`), and `core` VM (FreeRADIUS).
4. Once `wpa_supplicant` indicates success and `dhclient` (run subsequently or as part of the script on `naun3`) obtains a local IP, record another timestamp. Calculate the difference to estimate onboarding delay.
5. Verify PDU session establishment for the `clients` DNN using `nr-cli ps-list` on the `ue` VM. Note the assigned 5GC IP address for this PDU session.
6. Verify local IP address assignment to the `naun3` VM via `ip addr` on the `naun3` VM.
7. Verify the application of `iptables` and `ip rule`/`ip route` rules on the `ue` VM specific to the authenticated NAUN3 device's MAC and PDU session ID.
8. From the `naun3` VM, initiate `ping -R <PDU_Session_Gateway_IP>` (e.g., `10.46.0.1`). Analyze the recorded route to confirm it passes through the NAUN3's local IP, then its assigned 5GC PDU session IP, and then to the target.
9. From the `naun3` VM, run an `iperf3` client connecting to the `iperf3` server on the `core` VM.
10. Capture traffic on relevant interfaces (`naun3` LAN, `ue` VM LAN, `uesimtunX` on `ue` VM, `clientun0` on `core` VM) to observe data flow and NAT.

Metrics/Verification Points:
- Successful EAP-TLS authentication (logs).
- Onboarding delay (timestamp difference).
- One new PDU session active on `clients` DNN for the 5G-RG, with a unique 5GC IP.
- NAUN3 device receives a local IP.
- Successful `ping` with recorded route showing NAT via the PDU session IP.
- Successful `iperf3` data transfer.
- Correct `iptables` and policy routing rules applied.
### Experiment 2: Multi-Device Connectivity, Traffic Isolation, and PDU Session Mapping
In this experiment the goal was to verify that when multiple NAUN3 devices connect simultaneously, each gets a unique dedicated PDU session, their traffic is correctly mapped and isolated, and NAT occurs via their respective PDU session IPs.

Procedure:
1. Start two (or more) `naun3` VMs (e.g., `naun3-1`, `naun3-2`), each configured with a client certificates for EAP-TLS.
2. Allow both devices to authenticate and establish their dedicated PDU sessions as per Experiment 1. Verify that two distinct `clients` PDU sessions are created by the 5G-RG, each with a unique 5GC IP address (e.g., `10.46.0.3` and `10.46.0.4`).
3. On the `ue` VM, verify that distinct sets of `iptables` and `ip rule`/`ip route` entries are created for each NAUN3 device, mapping each to its unique PDU session ID and tunnel interface.
4. From `naun3-1`, execute `ping -R <PDU_Session_Gateway_IP>`. Note the recorded route, particularly the 5GC IP address assigned to `naun3-1`'s PDU session.
5. From `naun3-2`, execute `ping -R <PDU_Session_Gateway_IP>`. Note the recorded route, verifying it uses a _different_ 5GC IP address assigned to `naun3-2`'s PDU session.
6. Start an `iperf3` server on the `core` VM.
7. Simultaneously (or sequentially) run `iperf3` clients from `naun3-1` and `naun3-2` to the server on the `core` VM.
8. Monitor traffic on the `ue` VM's `uesimtunX` interfaces and check `iptables` counters to confirm traffic from each NAUN3 device is routed through its distinct PDU session.
9. On the `core` VM (`iperf3` server logs), verify that connections are received from the distinct 5GC IP addresses assigned to each NAUN3's PDU session.

Metrics/Verification Points:
- - Each NAUN3 device establishes its own unique PDU session on the `clients` DNN with a distinct 5GC IP.
- `ping -R` from each NAUN3 shows a path NATted through its unique PDU session IP.
- `iperf3` server logs show connections from distinct PDU session IPs.
- `iptables` and routing rules correctly isolate traffic per device.
### Experiment 3: Lifecycle Management (Device Disconnection and Resource Cleanup)
In order to verify that when an NAUN3 device disconnects or becomes unreachable, its local authentication is revoked, its dedicated PDU session is terminated, and associated network resources (IP addresses, routing rules) are correctly cleaned up by the `interceptor` the following procedure was followed.

1. Successfully onboard a single NAUN3 device as per Experiment 1. Verify its PDU session is active and traffic flows.
2. Simulate device disconnection:
	1. Option A (Graceful): Stop the `wpa_supplicant` service on the `naun3` VM.
	2. Option B (Abrupt): Power off or disconnect the network interface of the `naun3` VM.
3. Monitor the `interceptor` logs on the `ue` VM for detection of device unreachability and initiation of cleanup procedures.
4. Verify the following cleanup actions occur:
	1. `hostapd` deauthenticates the client (logs).
	2. The `interceptor` removes the device's MAC from `/etc/allowed-macs.conf` and restarts `dnsmasq`.
	3. The `interceptor` removes the specific `iptables` and `ip rule`/`ip route` entries for the device.
	4. The `interceptor` initiates a PDU session release for the `clients` DNN using `nr-cli`.
	5. Verify using `nr-cli ps-list` that the PDU session is terminated.
	6. Verify in Open5GS SMF/UPF logs that the session and associated resources are released.

Metrics/Verification Points:
- Detection of device disconnection by the `interceptor`.
- Successful local deauthentication.
- Removal of DHCP permission.
- Correct removal of all associated `iptables` and routing rules.
- Successful PDU session termination in UERANSIM and Open5GS.
- Internal state of the `interceptor` (e.g., `allowedDevices` map) reflects the device removal.
### Experiment 4: Security Aspects Observation (Qualitative)
To qualitatively observe key security aspects of the implemented solution, such as the EAP-TLS handshake and traffic segregation, we performed the following:
1. During Experiment 1 (NAUN3 Device Onboarding):
	1. Use `tshark` or `tcpdump` on the `ue` VM's LAN interface to capture the EAPoL and EAP-TLS handshake.
	2. Use `tshark` or `tcpdump` on the `ue` VM's `backhaul` PDU session interface and on the `core` VM's interface connected to the `ue` VM's `backhaul` to capture RADIUS traffic.
2. During Experiment 2 (Multiple Devices):
	1. Observe the source and destination IPs of the RADIUS traffic on the `backhaul` PDU session.
	2. Observe the source and destination IPs of the NAUN3 user plane traffic on the respective `clients` PDU sessions.

Metrics/Verification Points:
- Observation of a complete and successful EAP-TLS handshake.
- Confirmation that RADIUS traffic is encapsulated and transported over the `backhaul` PDU session.
- Confirmation that user plane traffic for each NAUN3 device is transported over its distinct `clients` PDU session.
- Review of Open5GS logs (AMF, SMF) to confirm that only the 5G-RG's identity (IMSI) is used for 5GC signaling, and NAUN3 MAC addresses are not propagated to these core NFs.

These experiments are designed to provide a holistic view of the system's functionality, its ability to manage multiple devices correctly, handle their lifecycle, and maintain basic security and traffic segregation principles, incorporating the specific types of data you've collected.