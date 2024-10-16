# 4.12.2 Registration via Untrusted non-3GPP Access

The EAP-5G packets utilize the "Expanded" EAP type and the existing 3GPP Vendor-Id registered with IANA under the SMI Private Enterprise Code registry. The "EAP-5G" method is used between the UE and the N3IWF and is utilized only for encapsulating NAS messages (not for authentication). If the UE needs to be authenticated, mutual authentication is executed between the UE and AUSF. The details of the authentication procedure are specified in TS 33.501

## 4.12.2.2 Registration procedure for untrusted non-3GPP access

![[Pasted image 20241001102256.png]]