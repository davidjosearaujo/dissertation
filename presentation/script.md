##  **Slide 1: Title Slide**

**(30 seconds)**
“Good \[morning/afternoon], everyone. My name is David Araújo, and today I’m presenting my Master’s dissertation: *Integration of Wi-Fi-Only Devices in 5G Core Networks: Addressing Authentication and Identity Management Challenges*. This work was developed at the University of Aveiro with the collaboration of Altice Labs and Instituto de Telecomunicações, under the supervision of Dr. Daniel Corujo and Dr. Francisco Fontes.”

---

##  **Slide 2: The Core Problem and Its Significance**

**(1.5 minutes)**
“The core challenge tackled in this work is that, as it stands, Wi-Fi-only devices—those that do not have 5G credentials and thus can't be directly integrated into the 5G Core network using standard 3GPP methods. This is a real-world limitation in enterprise and residential environments where many devices rely solely on Wi-Fi.

As 5G expands, the lack of a seamless integration path for these devices becomes a bottleneck. Solving this gap is crucial to unlock true Wi-Fi–5G convergence and extend key 5G benefits—such as enhanced mobile broadband, massive IoT support, and low-latency applications—to legacy or resource-constrained Wi-Fi-only devices.”

---

##  **Slide 3: Research Objectives**

**(1.5 minutes)**
“To solve this integration challenge, this work focused on three main research goals:

1. First, to design a secure local authentication mechanism that does not rely on 5G credentials.
2. Second, to develop a way for the 5G Core to recognize and individually manage each Wi-Fi-only device connection.
3. And third, to combine both into an integrated solution that introduces minimal changes to the existing 5G architecture or the end devices.

The approach needed to be practical, scalable, and transparent from the network’s point of view.”

---

##  **Slide 4: State of the Art – Device Gap**

**(1.5 minutes)**
“In the current landscape, there are two categories of non-3GPP devices.

* N5GC devices which can authenticate, but have limited capabilities.
* NAUN3 devices — Non-Authenticable Non-3GPP — which have no native support for 5G authentication and cannot be directly onboarded.

These NAUN3 devices are typically treated as a group behind a Residential Gateway. While manageable in small-scale scenarios, this approach lacks identity and traffic granularity. My work focuses on enabling **per-device authentication and management** for NAUN3 devices.”

---

##  **Slide 5: State of the Art – CGID Limitations**

**(1.5 minutes)**

“One mechanism proposed by 3GPP is CGID — Connectivity Group ID — which allows a group of devices behind a 5G-RG to share a single PDU session. But this approach fails when individual device identity and traffic isolation are needed.

Later developments in 3GPP, specifically in Release 19 in January this year, introduced the idea of per-device traffic distinction in QoS flows. However, these are mostly in early stages. This research anticipates those directions and provides a working, validated proof-of-concept that demonstrates per-device authentication and policy control using current tools.”

---

##  **Slide 6: Framework Overview and Principles**

**(2 minutes)**
“My framework is centered around a smart 5G Residential Gateway. The key principle here is **local intelligence** — all adaptation logic is placed at the 5G-RG. This allows:

* Minimal impact on NAUN3 devices — no firmware or configuration changes required.
* Minimal disruption to the 5GC — no additional interfaces or extensions were added to core components.

This edge-based approach is operator-friendly, scalable, and leverages familiar components like RADIUS, hostapd, and EAP.”

---

##  **Slide 7: Overall Architecture**

**(1.5 minutes)**
“This is the high-level view of the system.

* NAUN3 devices connect over Wi-Fi to the 5G-RG.
* The 5G-RG handles local certificate-based authentication along with an Authentication Server.
* Once authenticated, the RG requests a new PDU session and allocate it to handle this new device's traffic.
* Traffic from the device is then routed through this assigned PDU interface, achieving per-device IP allocation and routing.”

---

##  **Slide 8: Authentication Mechanism – EAP-TLS**

**(2 minutes)**
“For authentication, we use EAP-TLS, a widely accepted and secure method based on mutual certificate validation. Here’s how it works:

* The NAUN3 device holds a client certificate and acts as a supplicant.
* The 5G-RG uses `hostapd` to manage the Wi-Fi access point and relay EAP messages.
* The authentication server — FreeRADIUS — is operated by the ISP and validates the device certificate.

This setup achieves zero-trust-style security without relying on 5G credentials. It also reuses proven standards, making the solution practical and extensible.”

---

##  **Slide 9: Identity Management – Proxy PDU Sessions**

**(2 minutes)**
“Once a device is authenticated, it doesn’t yet have an identity in the 5G Core. So the RG creates a new PDU session and uses it as a **proxy identity** for the NAUN3 device.

A local mapping is maintained between each device’s MAC address and its PDU session. This enables:

* Per-device routing,
* Accurate lifecycle control,
* Individual policy application.

This method effectively gives a Wi-Fi-only device a unique 5G identity using standard mechanisms.”

---

##  **Slide 10: Traffic Management – Policy-Based Routing**

**(2 minutes)**
“Routing is handled using a policy-based approach:

1. Packets from each device are marked using their MAC address.
2. Marked packets are directed to a dedicated routing table.
3. That table routes traffic through the appropriate PDU session interface.
4. Finally, NAT ensures the traffic uses the IP assigned by the 5GC.

This gives complete traffic segregation between devices — even though they share the same physical Wi-Fi link and gateway.”

---

##  **Slide 11 and 12: Testbed Overview**

**(1 minute)**
“The testbed was built using Vagrant to orchestrate several virtual machines:

* Open5GS for the 5G Core,
* UERANSIM for gNB and UE emulation,
* FreeRADIUS, hostapd, and `wpa_supplicant`.

At the heart of the framework is a custom component which I developed: the *Interceptor*. It monitors authentication events and coordinates PDU session creation, DHCP permissions, routing table management, and teardown on disconnect.”

---

##  **Slide 13: Validation – Onboarding and PDU Creation**

**(1.5 minutes)**
“For each authenticated device:

* A dedicated PDU session was created.
* The 5GC assigned a unique IP address.
* The session appeared active and independently manageable.

This confirmed the viability of using proxy PDU sessions as per-device identities.”

---

##  **Slide 14: Validation – Connectivity and Isolation**

**(1.5 minutes)**
“To confirm routing and traffic isolation, we used `ping -R` and `iperf3`.
Each NAUN3 device could independently reach the 5GC using a different path and IP, and no cross-device traffic leakage occurred.

The system correctly enforced per-device NAT, routing, and IP attribution.”

---

##  **Slide 15: Validation – Lifecycle Management**

**(1.5 minutes)**
“When a device disconnected:

* It was deauthenticated.
* Its DHCP lease was revoked.
* Its routing rules were purged.
* The PDU session was released.

This clean lifecycle management demonstrates the framework’s operational viability in dynamic environments.”

---

##  **Slide 16: Key Contributions**

**(1 minute)**
“To summarize, this work contributes:

1. A full framework for onboarding Wi-Fi-only devices into 5GC.
2. A novel use of proxy PDU sessions for identity.
3. A practical integration of EAP-TLS and 5G session control.
4. A validated, open-source proof-of-concept with real 5G components.”

---

##  **Slide 17: Limitations**

**(1 minute)**
“There are a few limitations:

* Onboarding time averages 33 seconds.
* The Interceptor logic is CLI-driven which increased overhead.
* NAT prevents inbound connections to NAUN3 devices.
* Hardware integration was challenging do to experimental hardware.
* The robustness of the orchestration logic needs reinforcement.”

---

##  **Slide 18: Future Work**

**(1.5 minutes)**
“To address these issues, future work could:

* Reduce delay by using APIs or pre-pooling sessions.
* Explore high-performance mechanisms like eBPF.
* Add security layers to Interceptor and RADIUS transport.
* Address NAT using Framed-Route or UPF port forwarding.

These steps would bring the prototype closer to production-readiness.”

---

##  **Slide 18: Q\&A Slide**

**(30 seconds)**
“Thank you for your attention. I’m now happy to take any questions you may have.”

---

Would you like a version of this script formatted for cue cards or speaker notes, with time stamps for rehearsing?
