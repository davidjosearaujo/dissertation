# 4.2.8 Support of non-3GPP access
The following types of non-3GPP access networks are defined:
- Untrusted non-3GPP access networks;
- Trusted non-3GPP access networks; and
- Wireline access networks.
## General Concepts
The 5G Core Network supports connectivity of UEs via non-3GPP access networks, e.g. WLAN access networks.

Supports both untrusted non-3GPP access networks and trusted non-3GPP access networks (TNANs).

An **untrusted non-3GPP** access network **shall be connected** to the 5G Core Network **via a Non-3GPP InterWorking Function (N3IWF)**, via **N2** interface .

An untrusted non-3GPP access network shall be connected to the 5G Core Network via a Non-3GPP InterWorking Function (N3IWF), whereas a trusted non-3GPP access network shall be connected to the 5G Core Network via aTrusted Non-3GPP Gateway Function (TNGF).

A UE shall establish an IPsec tunnel with the N3IWF or with the TNGF in order to register with the 5G Core Network over non-3GPP access. Further details about the UE registration to 5G Core Network over untrusted non-3GPP access and over trusted non-3GPP access are described in clause 4.12.2 and in clause 4.12.2a of TS 23.502, respectively.
# 5.5.1 Registration Management
TODO
# 5.5.2 Connection Management
TODO
# 6.2.20 WAGF
[[Wireless and wireline convergence access support for the 5G System (5GS) - 23.316| The functionality of W-AGF is specified in TS 23.316]]

