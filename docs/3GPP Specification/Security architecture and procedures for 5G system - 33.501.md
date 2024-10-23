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
