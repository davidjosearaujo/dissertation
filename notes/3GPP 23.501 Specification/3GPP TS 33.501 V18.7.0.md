# 7 Security for non-3GPP access to the 5G core network

Security for untrusted non-3GPP access to the 5G Core network is achieved by a procedure using IKEv2 as defined in RFC 7296 to set up one or more IPsec ESP security associations. The role of IKE initiator (or client) is taken by the UE, and the role of IKE responder (or server) is taken by the N3IWF.

During this procedure, the AMF delivers a key N3IWF to the N3IWF. The AMF derives the key KN3IWF from the key KAMF. The key KN3IWF is then used by UE and N3IWF to complete authentication within IKEv2.