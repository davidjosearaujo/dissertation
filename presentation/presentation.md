---
marp: true
class: lead
size: 4K
style: |
    img[alt~="center"] {
        display: block;
        margin: 0 auto;
    }

    .columns2 {
        display: grid;
        grid-template-columns: repeat(2, minmax(0, 1fr));
        gap: 2rem;
    }

    .columns3 {
        display: grid;
        grid-template-columns: repeat(3, minmax(0, 1fr));
        gap: 2rem;
    }
---

![height:120px](./images/institutotelecomunicacoes.png) &nbsp; ![height:120px](./images/alticelabs.png) &nbsp; ![height:120px](./images/ua.png)

# Integration of Wi-Fi-Only Devices in 5G Core Networks: Addressing Authentication and Identity Management Challenges

<div class="columns2">
<div>

###### Author

David Araújo, _DETI_, _IT_
_davidaraujo@ua.pt_

</div>
<div>

###### Supervisors

Doctor Daniel Nunes Corujo, _DETI_, _IT_
Doctor Francisco Fontes, _Altice Labs_

</div>
</div>

<!-- header: Masters in Cybersecurity -->
<!-- footer: June 2025 &nbsp;—&nbsp; Aveiro, PT-->
---
<!-- paginate: true -->
<!-- header: Masters in Cybersecurity — Integration of Wi-Fi-Only Devices in 5G Core Networks: Addressing Authentication and Identity Management Challenges -->
<!-- footer: ![width:40px](./images/institutotelecomunicacoes.png) &ensp; ![width:90px](./images/alticelabs.png) &ensp; ![width:35px](./images/ua.png) &ensp; Instituto de Telecomunicações, Altice Labs and DETI-->

# The Core Problem and Its Significance

<div class="columns3">
<div>

## The Challenge

Current 3GPP standards don't fully address integrating **Wi-Fi-only devices lacking 5G credentials** into the 5G network, preventing standard 5G authentication.

</div>
<div>

## Impact

A significant hurdle for enterprise/residential environments with many such devices.

</div>
<div>

## Motivation

Solving this is crucial for 5G's success, enabling true **5G-Wi-Fi convergence** and extending 5G benefits (eMBB, mMTC, URLLC) to this vast device ecosystem.

</div>
</div>

---

# Research Objectives

To address this problem, this research aimed to:

1. **Investigate Secure Authentication:** Design a robust local authentication mechanism.
2. **Develop Device Identity Management:** Propose a method for 5GC to recognize and manage these device connections individually.
3. **Propose an Integrated Solution:** Develop a framework for seamless, secure integration with minimal impact.

---

# State of the Art

## The Gap

<div class="columns2">
<div>

**Device Types**

- **3GPP:** Have credentials and connect to the cellular network.
- **Non-3GPP:** Use technologies other than cellular and may or may not have 5G credentials.

</div>
<div>

_WiFi-only_ ➜ **Non-Authenticable Non-3GPP (NAUN3)** 

A robust mechanism for **individualized**, **secure authentication** of these devices and their subsequent per-device management within the 5GC is the focus of this project.

</div>
</div>

---

###### State of the Art

## Managing Device Groups (CGID)

<div class="columns2">
<div>

Connectivity Group ID (CGID) can manage **groups of devices behind** a 5G-RG with one PDU Session.

This does not provide per-device traffic management granularity.

Later developments envision a **network capable of distinguishing traffic from specific devices** behind an RG.

</div>
<div>

![height:400px center](images/cgid.png)

</div>
</div>

---

# Framework Concept and Architecture

## Overview and Guiding Principles

A _smart_ 5G Residential Gateway (5G-RG) capable of mediating the secure integration.

<div class="columns2">
<div>

### Key Design Principles

- Adaptation logic centralized at the 5G-RG.
- Minimal impact on end-devices and 5GC.

</div>
<div>

![center](images/conceptual-diagram.png)

</div>
</div>

---

###### Framework Concept and Architecture

## Overall Architecture

![center](images/general-envisioned-topology.png)

---

<div class="columns2">
<div>

###### Framework Concept and Architecture

## Authentication Mechanism

EAP-TLS is used for mutual, certificate-based local authentication.

- NAUN3 Device (**Supplicant**): Holds a client certificate.
- 5G-RG (**Authenticator**/Relay): Uses hostapd to relay EAP messages.
- RADIUS **Authentication Server**: ISP-operated, validates the device's certificate.

</div>
<div>

![height:500px center](images/eap-tls-topology.png)

</div>
</div>

---

###### Framework Concept and Architecture

## Identity Management (PDU Session as Proxy)

<div class="columns2">
<div>

After successful EAP-TLS authentication:

1. The 5G-RG requests a **new, dedicated** PDU Session.
2. This PDU Session becomes the **dynamic proxy identity** for the NAUN3.
3. The 5G-RG maintains a **mapping table** with NAUN3 MAC Addresses to PDU Session ID.

</div>
<div>

![width:500px center](images/identity-management-pdu-proxy-identity.png)

</div>
</div>

---

###### Framework Concept and Architecture

## Traffic Management and Policy-Based Routing

<div class="columns2">
<div>

1. **Packet Marking:** Incoming packets from the NAUN3's MAC are marked.
2. **Policy Routing:** Marked packets are directed to a specific table.
3. **Dedicated Route:** Traffic is routed via to a unique PDU interface.
4. **NAT:** Traffic is then masqueraded using the PDU session's 5GC-assigned IP address.

</div>
<div>

![height:400px center](images/policy-based-routing.png)

</div>
</div>

---

# Testbed and _Interceptor_: Central Orchestrator

![height:480px center](images/emulated-environment-topology.png)

---

###### Testbed and _Interceptor_: Central Orchestrator

**Virtualized testing environment** with Vagrant, Open5GS, UERANSIM, FreeRADIUS, `hostapd`, and `wpa_supplicant`.

<div class="columns2">
<div>

The custom logic developed, **_Interceptor_**, is the **brain of the solution**, responsible for:

- ☑ Monitor `hostapd`
- ☑ Trigger new PDU Sessions
- ☑ Configure DHCP and routing
- ☑ Clean up on disconnect

</div>
<div>

![height:320px center](images/emulated-environment-topology.png)

</div>
</div>

---

# Validation 

## Successful Onboarding and PDU Creation

<div class="columns2">
<div>

Local EAP-TLS authentication was consistently successful.

Each authenticated NAUN3 device triggered the 5G-RG to establish a unique, **dedicated PDU session in "clients" DNN**, and the 5GC assigned a **unique IP to each session**.

</div>
<div>


```
PDU Session2: 
 state: PS-ACTIVE
 session-type: IPv4
 apn: clients
 s-nssai: 
  sst: 0x01
  sd: null
 emergency: false
 address: 10.46.0.2
 ambr: up[1000000Kb/s] down[1000000Kb/s]
 data-pending: false
```

</div>
</div>

---

###### Validation

## End-to-End Connectivity and Traffic Isolation

<div class="columns2">
<div>

```
PING 10.46.0.1 (10.46.0.1) 56(124) bytes of data .
64 bytes from 10.46.0.1: icmp_seq =1 ttl =63 time =1.52 ms
RR :    192.168.59.100
    10.46.0.2
    10.46.0.1
    10.46.0.1
    192.168.59.10
    192.168.59.100
(...)
```

```
PING 10.46.0.1 (10.46.0.1) 56(124) bytes of data .
64 bytes from 10.46.0.1: icmp_seq =1 ttl =63 time =2.31 ms
RR :    192.168.59.15
    10.46.0.3
    10.46.0.1
    10.46.0.1
    192.168.59.10
    192.168.59.15
(...)
```

</div>
<div>

Using `ping -R` and `iperf3` we can confirm that traffic from different NAUN3 devices was **correctly and separately routed through their respective PDU session IPs**, confirming successful **traffic isolation** and NAT.

</div>
</div>

---

###### Validation

## Lifecycle Management and Onboarding Delay

**Onboarding Delay:** The average time for the full process (EAP auth, PDU setup, local IP) was approximately 33 (± 5) seconds in the testbed.

<div class="columns2">
<div>

**Lifecycle:** When a device disconnected, the system correctly deauthenticated it, cleaned up all routing rules and DHCP permissions, and terminated the dedicated PDU session.

</div>
<div>

1. ☑ Deauthenticate
2. ☑ Disallow DHCP lease
3. ☑ Release dedicated PDU Session
4. ☑ Remove routing table

</div>
</div>

---

# Key Contributions

1. A practical, end-to-end framework for integrating _5G-credential-less_ Wi-Fi-only devices into 5G.
2. The innovative use of **per-device PDU Sessions as dynamic proxy identities**, orchestrated by an intelligent 5G-RG.
3. The tight coupling of strong, local EAP-TLS authentication with 5G PDU session management at the network edge.
4. A working proof-of-concept validating the architecture with open-source tools and custom logic.

---

# Limitations

- **Physical Hardware Integration:** Physical 5G modem integration is an ongoing challenge. Proprietary drivers, kernel dependencies, and lack of documentation for multi-PDU session management.

- **Implementation Specifics:** The simulated PoC relies on CLI-based orchestration (`nr-cli`), which is not ideal for performance. The onboarding delay of ~33 (±5) seconds reflects this.

- **NAT Implications:** Inbound connection initiation to NAUN3 devices is restricted.

---

# Future Work

- **Modem Interface Adaptation:** The _Interceptor_ logic must be adapted to interface with modem-specific APIs, such as AT commands or QMI, replacing the UERANSIM CLI used in the simulation.
- **Performance and Scalability Analysis:** Rigorous testing and exploring alternatives like eBPF.
- **Enhanced Robustness:** Harden the _Interceptor_ and secure RADIUS transport (e.g., with IPSec).
- **Address NAT:** Investigate solutions like Framed-Routing.

---

# Thank You and Q&A

<div class="columns2">
<div>

## Author

David Araújo, _DETI_, _IT_
_davidaraujo@ua.pt_

</div>
<div>

## Supervisors

Doctor Daniel Nunes Corujo, _DETI_, _IT_
Doctor Francisco Fontes, _Altice Labs_

</div>
</div>

## In Colaboration With
![height:50px](./images/institutotelecomunicacoes.png) &nbsp; ![height:50px](./images/alticelabs.png) &nbsp; ![height:50px](./images/ua.png)
