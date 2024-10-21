For an untrusted non-3GPP access network, the communication between the UE and the 5GCN is not trusted to be secure.

For an untrusted non-3GPP access network, to secure communication between the UE and the 5GCN, **a UE establishes secure connection to the 5G core network over untrusted non-3GPP access via the N3IWF**. The UE performs registration to the 5G core network during the IKEv2 SA establishment procedure as specified in TS 24.501 and RFC 7296. After the registration, the UE supports NAS signalling with 5GCN using the N1 reference point as specified in TS 24.501. The N3IWF interfaces the 5GCN CP function via the N2 interface to the AMF and the 5GCN UP functions via N3 interface to the UPF as described in TS 23.501.
# 4.2 Untrusted access
For an untrusted non-3GPP access network, the communication between the UE and the 5GCN is not trusted to be secure.

For an untrusted non-3GPP access network, to secure communication between the UE and the 5GCN, a UE establishes secure connection to the 5G core network over untrusted non-3GPP access via the N3IWF. The UE performs registration to the 5G core network during the IKEv2 SA establishment procedure as specified in TS 24.501 and RFC 7296. After the registration, the UE supports NAS signalling with 5GCN using the N1 reference point as specified in TS 24.501. The N3IWF interfaces the 5GCN CP function via the N2 interface to the AMF and the 5GCN UP functions via N3 interface to the UPF as described in TS 23.501.
# 4.3 Identities
## 4.3.1 User identities
When the UE accesses the 5GCN over non-3GPP access networks, the same permanent identities for 3GPP access are used to identify the subscriber for non-3GPP access authentication, authorization and accounting services.

The Subscription Permanent Identifier (SUPI) is defined in TS 33.501. The SUPI can contain an IMSI, a network specific identifier, a GCI or a GLI as specified in TS 23.501. A SUPI containing an IMSI is defined in TS 23.003. A SUPI containing a network specific identifier, a GCI or a GLI always takes the form of a NAI as defined in TS 23.003.

The Subscription Concealed Identifier (SUCI) is a privacy preserving identifier containing the concealed SUPI as specified in TS 33.501. SUCI is calculated from SUPI. When the SUPI contains an IMSI, the corresponding SUCI is derived as specified in TS 23.003. When the SUPI contains a network specific identifier, a GCI or a GLI, the corresponding SUCI in NAI format is derived as specified in TS 23.003.
# 4.5 Trusted Access
For a trusted non-3GPP access network, the communication between the UE and the 5GCN is secure. A trusted non-3GPP access network is connected to the 5GCN via a trusted non-3GPP gateway function (TNGF) as specified in TS 23.501. The TNGF interfaces the 5GCN CP function via the N2 interface to the AMF and the 5GCN UP functions via N3 interface to the UPF as described in TS 23.501.

or a trusted non-3GPP access network, the UE establishes secure connection to the 5GCN over trusted non-3GPP access to the TNGF. The UE uses 3GPP-based authentication for connecting to a non-3GPP access and establishes an IPsec Security Association (SA) with the TNGF in order to register to the 5GCN by using the registration procedure as specified in TS 24.501. After the registration, the UE supports NAS signalling with the 5GCN using the N1 reference point as specified in TS 24.501.
# 6 UE - 5GC network protocols
