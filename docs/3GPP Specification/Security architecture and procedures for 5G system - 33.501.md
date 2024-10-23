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