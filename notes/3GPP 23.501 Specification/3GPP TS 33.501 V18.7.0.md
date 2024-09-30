# 7 Security for non-3GPP access to the 5G core network

Security for untrusted non-3GPP access to the 5G Core network is achieved by a procedure using IKEv2 as defined in RFC 7296 to set up one or more IPsec ESP security associations. The role of IKE initiator (or client) is taken by the UE, and the role of IKE responder (or server) is taken by the N3IWF.

During this procedure, the AMF delivers a key K_N3IWF to the N3IWF. The AMF derives the key K_N3IWF from the key K_AMF. The key K_N3IWF is then used by UE and N3IWF to complete authentication within IKEv2.

###  7.2.1 Security procedures

This clause specifies how a UE is authenticated to 5G network via an untrusted non-3GPP access network. It uses a vendor-specific EAP method called "EAP-5G", utilizing the "Expanded" EAP type and the existing 3GPP Vendor-Id, registered with IANA under the SMI Private Enterprise Code registry. The "EAP-5G" method is used between the UE and the N3IWF and is utilized for encapsulating NAS messages. If the UE needs to be authenticated by the 3GPP home network, any of the authentication methods as described in clause 6.1.3 can be used. The method is executed between the UE and AUSF as shown below. 

When possible, the UE shall be authenticated by reusing the existing UE NAS security context in AMF.

![[Pasted image 20240930145000.png]]
# 7A Security for trusted non-3GPP access to the 5G core network

## 7A.2.1	Authentication for trusted non-3GPP access

This clause specifies how a UE is authenticated to 5G network via a trusted non-3GPP access network. 

This is based on the specified procedure in [[3GPP TS 23.502 V18.7.0]] clause 4.12a.2.2 "Registration procedure for trusted non-3GPP access". The authentication procedure is similar to the authentication procedure for Untrusted non-3GPP access defined in clause 7.2.1 with few differences, which are mentioned below: