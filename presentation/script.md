<!--
 Copyright 2025 David Araújo
 
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
 
     https://www.apache.org/licenses/LICENSE-2.0
 
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
-->
### **Slide 5: State of the Art and The Specific Gap**

**(What to say):**

"To set the stage, let's briefly review the state of the art. 5G's Service Based Architecture includes key network functions for authentication and session management. It also defines mechanisms for Non-3GPP access.

When we look at devices behind a Residential Gateway, 3GPP distinguishes between two main types: N5GC devices, which can be authenticated by the 5G Core, and NAUN3 devices, which cannot be directly authenticated and are often managed in groups.

This leads us to the specific gap that my thesis addresses: While 3GPP has solutions for N5GC devices and for grouping NAUN3 traffic, a robust mechanism for the *individualized, secure authentication* of USIM-less Wi-Fi-only devices, and their subsequent *per-device management* within the 5GC, is not fully detailed. My work focuses on filling this gap."

---

### **Slide 6: Proposed Framework: Overview and Guiding Principles**

**(What to say):**

"To solve this problem, I propose a framework where an intelligent 5G Residential Gateway, or 5G-RG, acts as a mediator.

This solution is guided by three key design principles:

First, all the new adaptation logic is centralized within the 5G-RG itself.

Second, there is minimal impact on the end-devices; they only need a standard EAP supplicant, which is widely available.

And third, there is minimal impact on the 5G Core, as the 5G-RG interacts with it just like any standard User Equipment."

---

### **Slide 7: Proposed Framework: Overall Architecture**

**(What to say):**

"This is the high-level architecture of the proposed solution. There are three main communication paths.

First, on the local network, the NAUN3 device uses EAP over LAN for authentication signaling with the 5G-RG.

Second, for the authentication backhaul, the 5G-RG relays these EAP messages to an external Authentication Server using the RADIUS protocol. This is securely transported over a dedicated 'Backhaul' PDU Session.

Finally, for client traffic, once a device is authenticated, the 5G-RG establishes a unique and dedicated 'Client' PDU Session for that device, which carries all of its user data."

---

### **Slide 8: Proposed Framework: Authentication Mechanism (EAP-TLS)**

**(What to say):**

"For the authentication mechanism itself, I chose EAP-TLS. This method provides strong, mutual authentication based on certificates.

In this model, the NAUN3 device acts as the Supplicant, holding a client certificate. The 5G-RG is the Authenticator, relaying messages. And an external, ISP-operated RADIUS server is the Authentication Server, which validates the device's certificate.

This process ensures that only legitimate and verified devices are authenticated locally before any 5G network resources are committed to them."

---

### **Slide 9: Proposed Framework: Identity Management (PDU Session as Proxy)**

**(What to say):**

"The core innovation of this framework is how we manage device identity. Since these devices lack a 5G identity like a SUPI, we create a proxy for it.

After a successful local EAP-TLS authentication, the 5G-RG requests a new, dedicated PDU Session from the 5G Core, specifically for that authenticated device.

This PDU Session, with its unique IP address assigned by the 5G Core, effectively becomes the dynamic proxy identity of the NAUN3 device within the 5G system. To manage this, the 5G-RG maintains an internal mapping of the device's MAC Address to its assigned PDU Session."

---

### **Slide 10: Proposed Framework: Traffic Management and Routing**

**(What to say):**

"To ensure traffic from each device is handled correctly, the 5G-RG implements a policy-based routing scheme.

First, incoming packets from a device's MAC address are marked with a unique ID.

Second, a policy rule directs these marked packets to a dedicated routing table.

Third, this table routes all traffic through the device's specific PDU session tunnel interface.

Finally, as traffic exits, its source address is translated to the 5G-assigned IP of that PDU session, a process known as NAT."

---

### **Slide 11: Implementation: Testbed, Components, and `interceptor` Logic**

**(What to say):**

"The framework was implemented and validated using a virtualized testbed. This setup included standard open-source components like Open5GS for the 5G Core, UERANSIM for the RAN and UE simulators, and FreeRADIUS for authentication.

The brain of the solution is a custom application I developed, called the 'interceptor'. This Go application runs on the 5G-RG and orchestrates the entire process. It monitors authentication events, manages the PDU session lifecycle, controls local network access, and dynamically configures the routing rules."

---

### **Slide 12: Validation: Successful Onboarding and PDU Creation**

**(What to say):**

"The first validation test confirmed successful device onboarding. The logs showed that EAP-TLS authentication completed successfully for each device.

The key result here is that each authenticated device triggered the 5G-RG to establish a unique and dedicated PDU session on the 'clients' network. As you can see from the command-line output, in addition to the primary backhaul session, new PDU sessions are created, each with a unique IP address assigned by the 5G Core."

---

### **Slide 13: Validation: End-to-End Connectivity and Traffic Isolation**

**(What to say):**

"Next, I verified end-to-end connectivity and traffic isolation. Ping tests with the 'record route' option confirmed that traffic from each device was correctly routed through its dedicated PDU session IP address.

Furthermore, `iperf3` traffic tests were successful, and the server logs show that it received connections from the distinct IP addresses of each device's PDU session. This confirms that traffic from different devices is correctly isolated and managed."

---

### **Slide 14: Validation: Lifecycle Management and Onboarding Delay**

**(What to say):**

"The system's lifecycle management was also validated. When a device was disconnected, the interceptor correctly detected it, initiated local deauthentication, cleaned up all associated routing rules and DHCP permissions, and terminated the dedicated PDU session in the 5G Core.

I also measured the onboarding delay for the entire process, from the start of authentication to the device receiving a local IP. The average time in the testbed was approximately 33 seconds."

---

### **Slide 15: Conclusion**

**(What to say):**

"In conclusion, this dissertation successfully designed, implemented, and validated a novel gateway-centric framework.

This framework enables the secure integration of Wi-Fi-only devices that lack USIMs into 5G networks.

The core mechanism is the use of local EAP-TLS authentication to trigger the creation of per-device PDU sessions, which act as proxy identities within the 5G Core.

The proof-of-concept demonstrated that this approach is achievable with minimal impact on standard 5G components and on the end-user devices themselves."

---

### **Slide 16: Key Contributions**

**(What to say):**

"The key contributions of this work are:

First, a practical, end-to-end framework for integrating USIM-less Wi-Fi-only devices into 5G.

Second, the innovative use of per-device PDU Sessions as dynamic proxy identities, orchestrated by an intelligent 5G-RG.

Third, the tight coupling of strong, local EAP-TLS authentication with 5G session management at the network edge.

And finally, a working proof-of-concept that validates the architecture using open-source tools and custom logic."

---

### **Slide 17: Limitations**

**(What to say):**

"Of course, this research has some limitations. The onboarding delay of 33 seconds in the proof-of-concept has room for optimization. The system's scalability was not stress-tested, and the CLI-based orchestration could be a bottleneck. The use of NAT restricts inbound connections to the devices, and I encountered challenges with physical modem integration. Finally, the custom interceptor logic would require further security hardening for a production environment."

---

### **Slide 18: Future Work**

**(What to say):**

"Based on these limitations, there are several avenues for future work. This includes optimizing the onboarding delay, possibly with API-based controls. A rigorous performance and scalability analysis is needed. Further work could also involve enhancing security, deeper integration with the 5G Policy Control Function for per-device QoS, and investigating solutions to address the NAT limitations."

---

### **Slide 19: Thank You and Questions**

**(What to say):**

"Thank you for your time and attention. I am now happy to answer any questions you may have."
