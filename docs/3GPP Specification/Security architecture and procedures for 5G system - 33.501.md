# 4.1 Security domains 
![[2024-10-23_10-39.png]]
# 4.3 Security entities in the 5G Core network
The 5G System architecture introduces the following security entities in the 5G Core network:
- **AUSF**: 	AUthentication Server Function;
- **ARPF**: 	Authentication credential Repository and Processing Function;
- **SIDF**:	Subscription Identifier De-concealing Function;
- **SEAF**:    SEcurity Anchor Function.
# 5.2 Requirements on the UE
## 5.2.5 Subscriber privacy
- The UE shall support 5G-GUTI
- The Home Network Public Key shall be stored in the USIM.
- The protection scheme identifier shall be stored in the USIM.
- The Home Network Public Key Identifier shall be stored in the USIM.
- The SUCI calculation indication, either USIM or ME calculating the SUCI, shall be stored in USIM.
# 5.5 Requirements on the AMF
## 5.5.3 Subscriber privacy
- The AMF shall support to trigger primary authentication using the SUCI.
- The AMF shall support assigning 5G-GUTI to the UE.
- The AMF shall support reallocating 5G-GUTI to UE.
# 5.6	Requirements on the SEAF
The security anchor function (SEAF) provides the authentication functionality via the AMF in the serving network. The SEAF shall fulfil the following requirements:
- The SEAF shall support primary authentication using SUCI.
# 5.8	Requirements on the UDM 
## 5.8.1	Generic requirements
The long-term key(s) used for authentication and security association setup purposes shall be protected from physical attacks and shall never leave the secure environment of the UDM/ARPF unprotected.
## 5.8.2	Subscriber privacy related requirements to UDM and SIDF
The SIDF is responsible for de-concealment of the SUCI and shall fulfil the following requirements:
- The SIDF shall be a service offered by UDM.
- The SIDF shall resolve the SUPI from the SUCI based on the protection scheme used to generate the SUCI.

The Home Network Private Key used for subscriber privacy shall be protected from physical attacks in the UDM. 

The UDM shall hold the Home Network Public Key Identifier(s) for the private/public key pair(s) used for subscriber privacy.

The algorithm used for subscriber privacy shall be executed in the secure environment of the UDM.
# 5.8a Requirements on AUSF
The Authentication server function (AUSF) shall handle authentication requests for both, 3GPP access and non-3GPP access.
# 6 Security procedures between UE and 5G network functions
## 6.1 Primary authentication and key agreement
### 6.1.1 Authentication framework
#### 6.1.1.1 General
UE and serving network shall support EAP-AKA' and 5G AKA authentication methods.

UE and serving network shall support EAP-AKA' and 5G AKA authentication methods.

> Note 4: **EAP-AKA' and 5G AKA are the only authentication methods that are supported in UE and serving network**, hence only they are described in sub-clause 6.1.3 of the present document. For a private network using the 5G system as specified in an example of how additional authentication methods can be used with the EAP framework is given in the informative Annex B.
#### 6.1.1.2 EAP framework
The EAP framework is specified in RFC 3748. It defines the following roles: **peer**, **pass-through authenticator** and **back-end authentication server**. The back-end authentication server acts as the EAP server, which terminates the EAP authentication method with the peer. In the 5G system,  the EAP framework is supported in the following way:
- The UE takes the role of the peer.
- The SEAF takes the role of pass-through authenticator.
- The AUSF takes the role of the backend authentication server.
## 6.1.2  Initiation of authentication and selection of authentication method
![[Pasted image 20241023152054.png]]
## 6.1.3 Authentication procedures
### 6.1.3.1 Authentication procedure for EAP-AKA'
The 3GPP 5G profile for EAP-AKA' is specified in the normative Annex F. The selection of using EAP-AKA' is described in sub-clause 6.1.2 of the present document.
![[Pasted image 20241023153444.png]]
### 6.1.3.2 Authentication procedure for 5G AKA
5G AKA enhances EPS AKA by providing the home network with proof of successful authentication of the UE from the visited network. The proof is sent by the visited network in an Authentication Confirmation message.
![[Pasted image 20241023153659.png]]
## 6.2 Key hierarchy, key derivation, and distribution scheme
### 6.2.1 Key hierarchy
![[Pasted image 20241023161036.png]]
### 6.2.2	Key derivation and distribution scheme
#### 6.2.2.2 Keys in the UE
For every key in a network entity, there is a corresponding key in the UE.
![[Pasted image 20241023172927.png]]
# 7 Security for non-3GPP access to the 5G core network
## 7.2 Security procedures
### 7.2.1 Authentication for Untrusted non-3GPP Access
This clause specifies how a UE is authenticated to 5G network via an untrusted non-3GPP access network. It uses a vendor-specific EAP method called "EAP-5G", utilizing the "Expanded" EAP type and the existing 3GPP Vendor-Id, registered with IANA under the SMI Private Enterprise Code registry.

 TheEAP-5G" method is used between the UE and the N3IWF and is utilized for encapsulating NAS messages. If the UE needs to be authenticated by the 3GPP home network, any of the authentication methods as described in clause 6.1.3 can be used. The method is executed between the UE and AUSF as shown below.
 ![[Pasted image 20241024105055.png]]
### 7A.2.1	Authentication for trusted non-3GPP access
This clause specifies how a UE is authenticated to 5G network via a trusted non-3GPP access network.

This is based on the specified procedure in TS 23.502 clause 4.12a.2.2 "Registration procedure for trusted non-3GPP access". The authentication procedure is similar to the authentication procedure for Untrusted non-3GPP access defined in clause 7.2.1 with few differences, which are mentioned below:
![[Pasted image 20241024112734.png]]
### 7A.2.4 Authentication for devices that do not support 5GC NAS over WLAN access
A N5CW device is capable to register to 5GC **with 3GPP credentials** and to establish 5GC connectivity via a trusted WLAN access network. The reference architecture is captured in clause 4.2.8.5.2 of TS 23.501. The 3GPP credentials are stored as defined in clause 6.1.1.1. The Trusted WLAN Interworking Function (TWIF) provides interworking functionality that enables connectivity with 5GC and implements the NAS protocol stack and exchanges NAS messages with the AMF on behalf of the N5CW device. A single EAP-AKA’ authentication procedure is executed for connecting the N5CW device both to the trusted WLAN access network and to the 5G core network.
![[Pasted image 20241024113650.png]]
# 7B Security for wireline access to the 5G core network
To support Wireless and Wireline Convergence for the 5G system, two new network entities, 5G-RG and FN-RG, are introduced in the architecture specificaction TS 23.501.

The 5G-RG acts as a 5G UE and can connect to 5GC via wireline access network (W-5GAN) or via Fixed Wireless Access (FWA). Existing security procedures defined in this document are reused. The 5G-RG also acts as end point of N1 and provides the NAS signaling connection to the 5GC on behalf of the AUN3 devices behind the 5G-RG.

The FN-RG can connect to 5GC via wireline access network (W-5GAN). The W-AGF performs the registration procedure on behalf of the FN-RG. It acts as end point of N1 and provides the NAS signalling connection to the 5GC on behalf of the FN-RG.

A 5G-capable UE can connect to 5GC through an RG that’s connected to the 5GC via wireline access network (W-5GAN) or NG-RAN. The UE supports untrusted non-3GPP access and/or trusted non-3GPP access.
## 7B.2	Authentication for 5G-RG
The 5G-RG can be connected to 5GC via W-5GAN, NG RAN or via both accesses. The registration procedure for the 5G-RG connecting to 5GC via NG-RAN is specified in TS 23.316 clause 4.11. The registration procedure for the 5G-RG connecting to 5GC via W-5GAN is specified in TS 23.316 clause 7.2.1.

The Untrusted non-3GPP access procedure defined in clause 7.2.1 is used as the basis for registration of the 5G-RG. The 5G-RG shall support both 5G-AKA and EAP-AKA’ and it shall be authenticated by the 3GPP home network. The 5G-RG is equivalent to a normal UE.

As 5G-RG is a UE from 5GC point of view, the authentication framework defined in clause 6.1.3 shall be used to authenticate the 5G-RG.

In case of 5G-RG connects to 5GC via 5G-RAN, comparing to clause 6.1, the difference is:

- UE is replaced by 5G-RG.

In case of 5G-RG connects to 5GC via W-5GAN, a W-CP protocol stack message shall be used between the 5G-RG and the W-5GAN for encapsulating NAS message. The authentication method is executed between the 5G-RG and AUSF as shown below.
![[Pasted image 20241024114720.png]]
## 7B.3 Authentication for FN-RG
The FN-RG connects to 5GC via W-5GAN, which has the W-AGF function that provides connectivity to the 5GC via N2 and N3 reference points. Since the FN-RG is a non-wireless entity defined by BBF or CableLabs, it doesn’t support N1. The W-AGF provides N1 connectivity on behalf of the FN-RG. The authentication method is executed between the FN-RG and AUSF as shown in Figure 7B.c.

The W-AGF may authenticate the FN-RG; this is controlled by local policies.

It is assumed that there is a trust relationship between the wireline operator that manages the W-5GAN and the PLMN operator managing the 5GC. The AMF trusts the W-5GAN based on mutual authentication executed when security is established on the interface between the two using NDS/IP or DTLS.
![[Pasted image 20241024114851.png]]
## 7B.4 Authentication for UE behind 5G-RG and FN-RG
A UE that is connected to a 5G-RG or FN-RG, can access the 5GC via the N3IWF or via the TNGF.

A **UE behind a FN-RG can use the untrusted non-3GPP access procedure** as defined in TS 23.502 clause 4.12.2.2 to access the 5GC via the N3IWF.

A **UE behind a 5G-RG can use either the untrusted non-3GPP access** as defined in TS 23.502 clause 4.12.2.2, or trusted N3GPP-access as defined in TS 23.502 clause 4.12a.2.2.

A **UE connecting to the 5G-RG or FN-RG via WLAN supporting IEEE 802.1X can use the NSWO authentication procedure** as specified in Annex S of the present document.
## 7B.7 Authentication for AUN3 devices behind 5G-RG
An AUN3 device behind 5G-RG, as defined in TS 23.316, shall be registered to the 5GC by the 5G-RG and shall be authenticated by 5GC using EAP-AKA’, as specified in RFC 5448
### 7B.7.2 Authentication for AUN3 devices not supporting 5G key hierarchy
![[Pasted image 20241024121648.png]]
### 7B.7.3 Authentication for AUN3 devices supporting 5G key hierarchy
![[Pasted image 20241024121712.png]]
# Annex O (Informative): Authentication for non-5G capable devices behind residential gateways
This annex describes the authentication procedure, using EAP-TLS as an example, for Non-5G Capable (N5GC) devices behind Residential Gateways (RGs) in private networks or in isolated deployment scenarios (i.e., roaming is not considered) with wireline access. The registration procedure of N5GC devices behind Cable RGs is described in clause 4.10a of TS 23.316

N5GC devices lack some key 5G capabilities, including NAS and the derivation of 5G key hierarchy. Those devices exist in wireline networks and need to be able to access the converged 5G core. The authentication of N5GC devices can be based on additional EAP methods other than EAP-AKA’.

The procedure in O.3 uses EAP-TLS as in Annex B as an example, but it differs from the Annex B in the following:
1. the W-AGF creates the registration request on behalf of the N5GC device,
2. 5G related parameters (including ngKSI and ABBA) are not sent to the N5GC device. When received from the AMF, these parameters are ignored by the W-AGF, and
3. Neither the N5GC device nor the AUSF derives any 5G related keys after EAP authentication.
## O.3 Authentication procedure
It uses EAP-TLS as an example, but other EAP methods can also be supported. The W-AGF acts on behalf of the N5GC device during the registration process. The link between the N5GC device and the RG, and between the RG and the W-AGF can be any data link (L2) that supports EAP encapsulation.
![[Pasted image 20241028143737.png]]