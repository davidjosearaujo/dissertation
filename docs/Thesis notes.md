[Gemini conversation](https://gemini.google.com/app/75915c139d9aa2a2)
# Methodology and Proposed Framework
## Overall Research Approach

## Query
1.
It was constructive. I tried using the existing capabilities and features of 5G to find new solutions to the problem.

2.
The TS documents didn't help much, as they don't quite envision this limitation on the network. Except on the last part I mention, the Connectivity Group Identifiers. Which although they do not explain how they are supposed to be implemented or how do they work, the PDU Session separations were a huge inspiration to my solution

The problem was the fact that wifi-only devices do not have credentials for 5G and thus are incapable of authentication to it, and thus the Core can't recognize them. How could we reference such devices in the overall network? That was the problem

3.
For requirements, I wanted that devices and the core should suffer as little of disturebance and modification as possible. I didn't wnat to mess with the 5G core not require used wifi-only devices to have special software to be able to connect. Basically, if something was to be modded, it should be the gateways since that is the equipment the ISP has access to.

That is why CGI were interesting. They devided groups of devices in PDU Sessions managed by it, and the communication between the gateway and the devices is up to it.
But I though "Why use a PDU Session for a group of devices, when we could use a PDU Session per Device, and that PDU Session could become a 'proxy identity' for the device in the network"? What this is our goal, make a the gateway manage the mapping between PDU Session and devices, and this will be transparent to both Core and devices, requiring them to change nothing.

But how does the gateway authenticates a new device? Well, the machine were the Core is running, is also running a EAP-TLS authentication service, and the gateway is the authenticatior and the device the supplicante with a given certificate. If when connecting the device successsfully authenticates, this triggers the gateway to request a new PDU Session to then map the devices traffic to
## Requirements Analysis
