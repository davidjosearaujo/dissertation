# 7 Security for non-3GPP access to the 5G core network

Security for untrusted non-3GPP access to the 5G Core network is achieved by a procedure using IKEv2 as defined in RFC 7296 to set up one or more IPsec ESP security associations. The role of IKE initiator (or client) is taken by the UE, and the role of IKE responder (or server) is taken by the N3IWF.

During this procedure, the AMF delivers a key K_N3IWF to the N3IWF. The AMF derives the key K_N3IWF from the key K_AMF. The key K_N3IWF is then used by UE and N3IWF to complete authentication within IKEv2.

###  7.2.1 Security procedures

This clause specifies how a UE is authenticated to 5G network via an untrusted non-3GPP access network. It uses a vendor-specific EAP method called "EAP-5G", utilizing the "Expanded" EAP type and the existing 3GPP Vendor-Id, registered with IANA under the SMI Private Enterprise Code registry. The "EAP-5G" method is used between the UE and the N3IWF and is utilized for encapsulating NAS messages. If the UE needs to be authenticated by the 3GPP home network, any of the authentication methods as described in clause 6.1.3 can be used. The method is executed between the UE and AUSF as shown below. 

When possible, the UE shall be authenticated by reusing the existing UE NAS security context in AMF.

![[Pasted image 20240930145000.png]]
1. The UE connects to an untrusted non-3GPP access network with procedures outside the scope of 3GPP. When the UE decides to attach to 5GC network, the UE selects an N3IWF in a 5G PLMN, as described in TS 23.501 clause 6.3.6.
2. The UE proceeds with the establishment of an IPsec Security Association (SA) with the selected N3IWF by initiating an IKE initial exchange according to RFC 7296. After step 2 all subsequent IKE messages are encrypted and integrity protected by using the IKE SA established in this step.
3. The UE shall initiate an IKE_AUTH exchange by sending an IKE_AUTH request message. The AUTH payload is not included in the IKE_AUTH request message, which indicates that the IKE_AUTH exchange shall use EAP signalling (in this case EAP-5G signalling). As per the RFC 7296 [25], in the IDi the UE shall set the ID type as ID_KEY-ID in this message and set its value equal to any random number. The UE shall not use its GUTI/SUCI/SUPI as the Id in this step. If the UE is provisioned with the N3IWF root certificate, it shall include the CERTREQ payload within the IKE_AUTH request message to request N3IWF’s certificate.