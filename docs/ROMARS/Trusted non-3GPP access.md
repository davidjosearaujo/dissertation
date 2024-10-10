Starting from 3GPP release 16, additional non-3GPP access mechanisms are defined, and specifically “trusted” ones. **Trusted means that the operator of the 5G Core Network is also in charge of operating and managing the non-3GPP access** which, in addition to WiFi, is also extended to terrestrial fixed access. These updates in rel-16 open to 4 new scenarios:
## Scenario 1

Mobile 5G **UE devices that can register, access and exploit Core Network services both by 3GPP New Radio access and trusted WiFi**. This is **similar to** what already allowed by **N3IWF** with the difference that WiFi access points are coordinated and under the control of the same 5G network operator (also in terms of QoS). This scenario can be called Mobile Wireless Access (MWA).
## Scenario 2

Mobile **UE devices which do not support 5G capabilities** (i.e., NAS protocol for N1 and NR-Uu) **but have 3GPP credentials** (USIM/eUICC), accessing exclusively to the CN via a trusted WiFi.
## Scenario 3

**Fixed 5G UE devices**, called 5G Residential Gateway (5G-RG) which combine 3GPP NR access with fixed access (specifications for DOCSIS are available, and more are under definition) to the same CN. This scenario is called Fixed Wireless Access (FWA).
## Scenario 4

**Fixed UE devices**, called Fixed Network Residential Gateway (FN-RG), **which do not support 5G capabilities but have 3GPP credentials**, accessing exclusively to the CN via fixed network (i.e., DOCSIS or Broadband fixed modem).

In all the above cases, the **key difference** with regard to the untrusted WiFi access (i.e., N3IWF) is that **security can be enabled directly on non-3GPP access layers,** such as IEEE 802.11, **and not as part of IKEv2 establishment**. Nonetheless, it was decided to maintain the same UE protocol stack, with the use of IPSec “VPN” but without encryption, this allows to develop in the UE the same protocol stack regardless on trusted or non-trusted access.