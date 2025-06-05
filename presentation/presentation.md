---
marp: true
class: lead
size: 4K
style: |
    .columns {
        display: grid;
        grid-template-columns: repeat(var(--columns), minmax(0, 1fr));
        gap: 1rem;
    }
---

![width:140px](./images/institutotelecomunicacoes.png) &nbsp; ![width:300px](./images/alticelabs.png) &nbsp; ![width:110px](./images/ua.png)

# Integration of Wi-Fi-Only Devices in 5G Core Networks: Addressing Authentication and Identity Management Challenges

<div class="columns2">
<div>


### Author

David Araújo, _DETI_, _IT_
_davidaraujo@ua.pt_

</div>
<div>

### Supervisors

Doctor Daniel Nunes Corujo, _DETI_, _IT_
Doctor Francisco Fontes, _Altice Labs_

</div>
</div>

<!-- header: Masters in Cybersecurity -->
<!-- footer: June 2025 &nbsp;—&nbsp; Aveiro, PT-->
---
<!-- paginate: true -->
<!-- header: Masters in Cybersecurity — Integration of Wi-Fi-Only Devices in 5G Core Networks: Addressing Authentication and Identity Management Challenges -->
<!-- footer: ![width:40px](./images/institutotelecomunicacoes.png) &ensp; ![width:90px](./images/alticelabs.png) &ensp; ![width:35px](./images/ua.png) &ensp; Instituto de Telecomunicações and Altice Labs-->

# Table of Contents


<div class="columns" style="--columns:2;">
<div>

1. The Core Problem and Its Significance
2. Research Objectives
3. State of the Art and The Specific Gap
4. Framework Concept and Architecture
5. Key Mechanisms: Authentication, Identity, Traffic

</div>
<div>

6. Implementation: Testbed and Orchestration Logic
7. Validation: Key Results
8. Conclusion and Contributions
9. Limitations and Future Work

</div>
</div>

---

# The Core Problem and Its Significance

<div class="columns" style="--columns:3;">
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