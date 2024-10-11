# 3.6 Support for Wi-Fi Only Devices
Current 3GPP standard does not define architecture to support Wi-Fi only devices without USIM connecting to 5G Core.

The UE is shown as a dual radio device supporting both 3GPP and non- 3GPP radio per reference architectures in TS 23.501 clause 4.2.8.2, even though connectivity over the 3GPP access is shown as optional

**A UE is required to support SIM based identity (SUPIbased or IMSI) and SIM credentials to connect to 5G Core over WLAN access.**

However, **most Wi-Fi only devices**, e.g. devices in enterprise deployments, **would not have a USIM included** even in cases where these devices can be upgraded to support 5G control plane (NAS) and user plane functionality. Hence, **it is important to enable Wi-Fi only devices without 3GPP IMSI based identity and SIM credentials to connect to 5G Core**, to provide 5G services and experiences across different enterprise and verticals deployments.

Also, for N5CW UE connecting to 5G Core via TWIF, 3GPP spec requires such UEs to support SIM based credentials since such devices are authenticated using EAP-AKA’ as described in clause 7A.2.4 in TS 33.501 and clause 7.3A.4.2 in TS 24.502

The Wi-Fi only UEs may or may not support 5G NAS functionality and typically do not include SIM based identity and credentials. The **Wi-Fi only UEs with 5G NAS** will need to **support EAP-5G, IKEv2 and 5G NAS** signalling. To support Wi-Fi only devices without USIM, the 5G Core network needs to support non-IMSI based identity in the form of NAI (username@realm) and non-AKA based authentication methods such as EAP-TLS or EAP-TTLS. ^478de6

# 4 Summary and Recommendations
## Support for Wi-Fi only devices
To facilitate adoption of Wi-Fi only UEs without USIM for private networks in the stand-alone non-public network (SNPN) deployments, it is highlighted that 3GPP needs to define architecture and procedures for supporting Wi-Fi only UEs with non-IMSI based identity and EAP-TLS/EAP-TTLS based authentication to access 5G SNPN core via N3IWF, TNGF or TWIF functions.