For an untrusted non-3GPP access network, the communication between the UE and the 5GCN is not trusted to be secure.

For an untrusted non-3GPP access network, to secure communication between the UE and the 5GCN, **a UE establishes secure connection to the 5G core network over untrusted non-3GPP access via the N3IWF**. The UE performs registration to the 5G core network during the IKEv2 SA establishment procedure as specified in TS 24.501 and RFC 7296. After the registration, the UE supports NAS signalling with 5GCN using the N1 reference point as specified in TS 24.501. The N3IWF interfaces the 5GCN CP function via the N2 interface to the AMF and the 5GCN UP functions via N3 interface to the UPF as described in TS 23.501.
# 4.2 Untrusted access
For an untrusted non-3GPP access network, the communication between the UE and the 5GCN is not trusted to be secure.

For an untrusted non-3GPP access network, to secure communication between the UE and the 5GCN, a UE establishes secure connection to the 5G core network over untrusted non-3GPP access via the N3IWF. The UE performs registration to the 5G core network during the IKEv2 SA establishment procedure as specified in TS 24.501 and RFC 7296. After the registration, the UE supports NAS signalling with 5GCN using the N1 reference point as specified in TS 24.501. The N3IWF interfaces the 5GCN CP function via the N2 interface to the AMF and the 5GCN UP functions via N3 interface to the UPF as described in TS 23.501.
# 4.3 Identities
## 4.3.1 User identities
When the UE accesses the 5GCN over non-3GPP access networks, the same permanent identities for 3GPP access are used to identify the subscriber for non-3GPP access authentication, authorization and accounting services.

