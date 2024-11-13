# Universal and constant identity
The identity of a 5G device is designed to be both constant and universal in certain ways, although there are nuances to consider.
## Constant Identity
### SIM Credentials
The credentials stored in the SIM (or equivalent) are typically constant.
### Immutable Nature
The fundamental credentials in the SIM or eSIM are generally not changed during the lifecycle of the device unless re-provisioned by the carrier or modified for a specific reason (e.g., switching mobile network operators).
## Universal Identity
### Network Independence
The identity associated with a 5G device is universal in the sense that it allows the device to connect to different parts of the mobile network, such as various antennas (base stations) or residential gateways, without changing its credentials.
### Roaming and Handover
#### Home and Visited Networks
The 5G identity allows devices to maintain their identity when moving between different network domains, like when roaming internationally or transitioning between different network segments of the same operator.
#### Network Interfaces (e.g., gNBs and Residential Gateways)
The device's identity remains valid across different access points, such as antennas (gNBs) or fixed-wireless access points (residential gateways).
A 5G device’s identity ensures that it can authenticate and maintain secure communication anywhere within the network coverage of its operator or partner networks during roaming. This universal recognition allows seamless mobility, regardless of which base station or residential gateway the device uses for access.
## Limitation of current prototype
Exactly, you've identified a critical limitation in your current prototype. The reliance on **MAC addresses** as identifiers for legacy IoT devices introduces issues related to **consistency** and **universality**:
### Why MAC-Based Identity Is Not Constant or Universal
#### MAC Address Characteristics
- **Device-Specific**: It only serves as a unique address within a given network and doesn't inherently convey a constant, network-wide identity.
- **Changeability**: MAC addresses can be changed or spoofed. Some devices allow for MAC randomization for privacy reasons, which means that the identity presented by the device can vary between different connections or networks.
#### Gateway Dependence
- **Local Scope**: A MAC address is specific to the local network layer (Layer 2) and doesn't persist across different gateways or broader network segments.
- **Limited Universality**: Unlike credentials stored in a SIM or eSIM that are recognized across an entire 5G network, MAC-based identification does not provide a universal identity
