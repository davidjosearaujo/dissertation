**Theme:** Integration of Wi-Fi-Only Devices in 5G Core Networks: Addressing Authentication and Identity Management Challenges **Total Time:** 20 Minutes (approx. 1 minute per content slide)

**Slide 1: Title Slide (~0.5 minute)**
- **Title:** Integração de Dispositivos Wi-Fi-Only em Redes 5G: Abordagem aos Desafios de Autenticação e Gestão de Identidade
- **Author:** David José Araújo Ferreira
- **Supervisors:** Doctor Daniel Nunes Corujo, Doctor Francisco Fontes
- **Institution:** Universidade de Aveiro
- **Date:** [Date of Defense]
- **(Optional: University/Department Logo)**

**Slide 2: Agenda / Presentation Outline (~0.5 minute)**
- **Content:**
    - The Core Problem and Its Significance
    - Research Objectives
    - State of the Art and The Specific Gap
    - Proposed Framework: Concept and Architecture
    - Proposed Framework: Key Mechanisms (Authentication, Identity, Traffic)
    - Implementation: Testbed and Orchestration Logic
    - Validation: Key Results
    - Conclusion and Contributions
    - Limitations and Future Work
- **Visual:** Simple bullet points.

**Slide 3: The Core Problem and Its Significance (1 minute)**
- **Content (Combines former Slides 3 and 5 - Chapter 1.2 and 1.1):**
    - **The Challenge:** Current 3GPP standards don't fully address integrating Wi-Fi-only devices _lacking USIMs_ into the 5G Core (5GC), preventing standard 5G authentication.
    - **Impact:** A significant hurdle for enterprise/residential environments with many such devices. WBA has noted this.
    - **Motivation:** Solving this is crucial for 5G's success, enabling true 5G-Wi-Fi convergence and extending 5G benefits (eMBB, mMTC, URLLC) to this vast device ecosystem.
- **Visual:** A clear diagram: Wi-Fi-only device ("No USIM") → Barrier → 5GC. Add a small icon/text indicating "Lost 5G Benefits."

**Slide 4: Research Objectives (1 minute)**
- **Content (Chapter 1.3 - Summarized):**
    - To address this problem, this research aimed to:
        1. **Investigate Secure Authentication:** Design a robust local authentication mechanism (focus: EAP-TLS).
        2. **Develop Device Identity Management:** Propose a method for 5GC to recognize and manage these devices individually.
        3. **Propose an Integrated Solution:** Develop a framework for seamless, secure integration with minimal impact.
- **Visual:** Clear, concise bullet points.

**Slide 5: State of the Art and The Specific Gap (1 minute)**
- **Content (Chapter 2.2, 2.4.1, 2.5):**
    - **5G Context:** Briefly mention SBA, key NFs, and standard Non-3GPP access.
    - **Device Types Behind RGs:** Explain the difference between N5GC devices (which can authenticate) and NAUN3 devices (which cannot directly authenticate and are often grouped by CGID).
    - **The Specific Gap Addressed:** While 3GPP handles N5GCs and NAUN3 _groups_, a robust mechanism for _individualized, secure authentication_ of USIM-less Wi-Fi-only devices (behaving like NAUN3s) and their subsequent _per-device management_ within the 5GC is the focus of this thesis.
- **Visual:** A diagram contrasting an N5GC device's path to 5GC (with an authentication checkmark) against multiple NAUN3-like devices behind an RG, highlighting the "Individual Authentication and Management Gap."

**Slide 6: Proposed Framework: Overview and Guiding Principles (1 minute)**
- **Content (Chapter 3.1, 3.5):**
    - **Solution Core Idea:** An intelligent 5G Residential Gateway (5G-RG) mediates the secure integration.
    - **Key Design Principles:**
        - All adaptation logic is centralized in the 5G-RG.
        - Minimal impact on end-devices (requires only a standard EAP supplicant).
        - Minimal impact on 5GC (the 5G-RG interacts as a standard UE).
- **Visual:** A high-level conceptual diagram showing: Wi-Fi Device ↔ **Intelligent 5G-RG (Mediation Logic)** ↔ 5G Core.

**Slide 7: Proposed Framework: Overall Architecture (1 minute)**
- **Content (Chapter 3.5 - Key interactions):**
    - Explain the main communication paths using the diagram.
    - **Local:** NAUN3 Device ↔ 5G-RG using EAPOL for authentication signaling.
    - **Backhaul:** 5G-RG ↔ External EAP Auth Server using RADIUS, transported over a dedicated "Backhaul" PDU Session.
    - **Client Traffic:** 5G-RG ↔ 5GC establishes a dedicated "Client" PDU Session for each authenticated NAUN3 device.
- **Visual:** Use your provided image `general-envisioned-topology.png` (Figure 3.3 from your dissertation).

**Slide 8: Proposed Framework: Authentication Mechanism (EAP-TLS) (1 minute)**
- **Content (Chapter 3.3):**
    - **Method:** EAP-TLS is used for strong, mutual, certificate-based local authentication. 
    - **Roles:**
        - **NAUN3 Device (Supplicant):** Holds a client certificate.
        - **5G-RG (Authenticator/Relay):** Uses `hostapd` to relay EAP messages.
        - **External RADIUS Server:** ISP-operated, validates the device's certificate.
    - This process ensures only legitimate devices are authenticated locally before any 5G network resources are allocated to them.
- **Visual:** Use Figure 3.1 from your dissertation (EAP-TLS Topology) to illustrate the roles.

**Slide 9: Proposed Framework: Identity Management (PDU Session as Proxy) (1 minute)**
- **Content (Chapter 3.4, esp. 3.4.2):**
    - **The Innovation:** After a successful local EAP-TLS authentication...
    - The 5G-RG requests a _new, dedicated PDU Session_ from the 5GC for _that specific NAUN3 device_.
    - This PDU Session (and its unique 5GC-assigned IP) becomes the **dynamic proxy identity** of the NAUN3 device within the 5G system.
    - **Internal Mapping:** The 5G-RG maintains a table: (NAUN3 MAC Address ↔ PDU Session ID and 5GC IP).
- **Visual:** A clear conceptual diagram showing: NAUN3 icon → "EAP Auth Success" checkmark → 5G-RG icon → arrow labeled "Requests and Establishes" → a box labeled "Unique PDU Session (Proxy Identity)" connecting to a 5GC icon.

**Slide 10: Proposed Framework: Traffic Management and Routing (1 minute)**
- **Content (Chapter 4.3.6):**
    - The 5G-RG implements policy-based routing for each NAUN3 device's traffic:
        1. **Packet Marking:** Incoming packets from the NAUN3's MAC are marked using `iptables`.
        2. **Policy Routing:** `ip rule` directs marked packets to a specific routing table.
        3. **Dedicated Route:** That table routes all traffic via the device's unique PDU session interface (e.g., `uesimtunX`).
        4. **NAT:** Traffic is then masqueraded using the PDU session's 5GC-assigned IP address.
- **Visual:** Use your provided image `policy-based-routing.png` (Figure 4.2 from your dissertation).

**Slide 11: Implementation: Testbed, Components, and `interceptor` Logic (1 minute)**
- **Content (Chapter 4.1, 4.5.1, 4.5.2):**
    - **Testbed:** Built using a virtualized environment with Vagrant, Open5GS, UERANSIM, FreeRADIUS, `hostapd`, and `wpa_supplicant`.
    - **Core Custom Logic:** The `interceptor` (a Go application on the 5G-RG) is the brain of the solution.
    - **Interceptor's Role:** It monitors `hostapd` for successful authentication, orchestrates PDU sessions using `nr-cli`, manages local DHCP permissions via `dnsmasq`, and controls the routing rules.
- **Visual:** Use your provided image `emulated-environment-topology.jpg` (Figure 5.1 from your dissertation). Add a callout box highlighting the `interceptor` and its key interactions within the `ue VM`.

**Slide 12: Validation: Successful Onboarding and PDU Creation (1 minute)**
- **Content (Chapter 5.3.1):**
    - **Finding:** Local EAP-TLS authentication was consistently successful.
    - **Key Result:** Each authenticated NAUN3 device triggered the 5G-RG to establish a unique, dedicated "clients" PDU session, and the 5GC assigned a unique IP to each session.
- **Visual:** Use Figure 5.3 from your dissertation, showing the output of `nr-cli ps-list` with the backhaul session and multiple "clients" PDU sessions active.

**Slide 13: Validation: End-to-End Connectivity and Traffic Isolation (1 minute)**
- **Content (Chapter 5.3.2, 5.3.3):**
    - **Finding:** `ping -R` and `iperf3` tests confirmed end-to-end connectivity.
    - **Key Result:** Traffic from different NAUN3 devices was correctly and separately routed through their respective PDU session IPs, confirming successful traffic isolation and NAT.
- **Visuals:** Show small, clear snippets from your dissertation: Figure 5.7/5.8 (`ping -R` output) and Figure 5.10 (`iperf3` server logs showing connections from different PDU session IPs).

**Slide 14: Validation: Lifecycle Management and Onboarding Delay (1 minute)**
- **Content (Chapter 5.3.1, 5.3.4):**
    - **Lifecycle:** When a device disconnected, the system correctly deauthenticated it, cleaned up all routing rules and DHCP permissions, and terminated the dedicated PDU session.
    - **Onboarding Delay:** The average time for the full process (EAP auth, PDU setup, local IP) was approximately 33 seconds in the testbed.
- **Visuals:** Mention key findings from Figures 5.12-5.15. You can use a small icon representing "Resource Cleanup" or a simplified log snippet showing the `interceptor` taking action.

**Slide 15: Conclusion (1 minute)**
- **Content (Chapter 6.1):**
    - Successfully designed, implemented, and validated a novel gateway-centric framework.
    - This framework enables the secure integration of Wi-Fi-only/NAUN3-type devices (that lack USIMs) into 5G networks.
    - The core mechanism is using local EAP-TLS authentication to trigger the creation of per-device PDU sessions, which act as proxy identities.
    - The proof-of-concept demonstrated this is achievable with minimal impact on standard 5GC components and end-user devices.
- **Visual:** Concise summary bullet points.

**Slide 16: Key Contributions (1 minute)**
- **Content (Chapter 6.2):**
    1. A practical, end-to-end framework for integrating USIM-less Wi-Fi-only devices into 5G.
    2. The innovative use of **per-device PDU Sessions as dynamic proxy identities**, orchestrated by an intelligent 5G-RG.
    3. The tight coupling of strong, local EAP-TLS authentication with 5G PDU session management at the network edge.
    4. A working proof-of-concept validating the architecture with open-source tools and custom logic.
- **Visual:** A numbered list of your main contributions.

**Slide 17: Limitations (1 minute)**
- **Content (Chapter 6.3):**
    - **Onboarding Delay:** Approximately 33s in the PoC.
    - **Scalability:** Not stress-tested; CLI-based orchestration is a potential bottleneck.
    - **NAT Implications:** Restricts inbound connection initiation to NAUN3 devices.
    - **Physical Hardware:** Challenges encountered with physical modem integration.
    - **Security:** The custom `interceptor` logic requires further hardening for prouction environments.
- **Visual:** Bullet points.

**Slide 18: Future Work (1 minute)**
- **Content (Chapter 6.4):**
    - **Optimize Onboarding Delay:** Explore API-based PDU control or pre-established session pools.
    - **Performance and Scalability Analysis:** Rigorous testing and exploring alternatives like eBPF.
    - **Enhanced Security:** Harden the `interceptor` and secure RADIUS transport (e.g., with IPSec).
    - **Deeper 5G Integration:** Link to PCF for per-device QoS/slicing.
    - **Address NAT:** Investigate solutions like Framed-Route or UPF port forwarding.
- **Visual:** Bullet points.

**Slide 19: Thank You and Questions (~1 minute)**
- **Content:** "Thank You", "Questions?"
- **Visual:** Your name, contact (optional). University logo.