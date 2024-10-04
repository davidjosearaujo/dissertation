# 3.6 Support for Wi-Fi Only Devices
Current 3GPP standard does not define architecture to support Wi-Fi only devices without USIM connecting to 5G Core.

**A UE is required to support SIM based identity (SUPIbased on IMSI) and SIM credentials to connect to 5G Core over WLAN access.**

However, **most Wi-Fi only devices**, e.g. devices in enterprise deployments, **would not have a USIM included** even in cases where these devices can be upgraded to support 5G control plane (NAS) and user plane functionality. Hence, **it is important to enable Wi-Fi only devices without 3GPP IMSI based identity and SIM credentials to connect to 5G Core**, to provide 5G services and experiences across different enterprise and verticals deployments.

The Wi-Fi only UEs may or may not support 5G NAS functionality and typically do not include SIM based identity and credentials. The **Wi-Fi only UEs with 5G NAS** will need to **support EAP-5G, IKEv2 and 5G NAS** signalling. To support Wi-Fi only devices without USIM, the 5G Core network needs to support non-IMSI based identity in the form of NAI (username@realm) and non-AKA based authentication methods such as EAP-TLS or EAP-TTLS. ^478de6