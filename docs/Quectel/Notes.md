> Quectel UMTS&LTE&5G modules are USB composite devices with multiple USB interfaces. Each USB interface supports different functionalities, which are implemented by loading different USB interface drivers. After a driver is loaded successfully, the corresponding device node is generated, which can be used by the Linux system to implement the module functionalities, such as AT command, GNSS, DIAG, log and USB network adapter. 
> 
> The following table describes the USB interface information of different modules in the Linux system, including USB driver, interface number, device name and interface function. 
> 
> You can obtain the corresponding VID, PID and interface information of the relevant model, and then port the USB interface driver listed in the following table.

[[Quectel_UMTS_LTE_5G_Linux_USB_Driver_User_Guide_V3.2.pdf#page=12&selection=25,0,113,1|Quectel_UMTS_LTE_5G_Linux_USB_Driver_User_Guide_V3.2, page 12]]

# RG255C series/ RM255C-GL
- VID: 0x2c7c
- PID: 0x0316

| USB Driver | Interface Number | Device Name         | Function                                                                   |
| ---------- | ---------------- | ------------------- | -------------------------------------------------------------------------- |
| USB serial | **2**            | **/dev/ttyUSB2**    | **AT Command**                                                             |
| QMI_WWAN   | 3                | wwan0 /dev/cdc-wdm0 | Configure the type of USBnet interface as RmNet by **AT+QCFG="usbned",0"** |
| MBIM       | 6 and 7          | wwan0 /dev/cdc-wdm0 | Configure the type of USBnet interface as MBIM by **AT+QCFG="usbnet",2**   |
# Dependencies
```bash
sudo apt install -y \ 
  linux-headers-$(uname -r) \
  linux-modules-extra-$(uname -r) \
  build-essential \
  net-tools \
  dwarves
```
# AT Commands
Connect to modem's terminal via serial with `sudo screen /dev/ttyUSB2 9600`(`minicom`acts weird, as in freezes or doesn't show entire response outputs). 
## [[Quectel_RG255C_Series_RM255C-GL_AT_Commands_Manual_V1.0.0_Preliminary_20231218.pdf#page=27&selection=37,0,45,13|Set UE Functionality]]
This command controls the functionality level. It can also be used for resetting the UE.
```
+CFUN: <fun>    // Read command result
AT+CFUN=<fun>[,<rst>]    // Write command
```
- `<fun>` - Functionality level
	- `0` - Minimum functionality
	- `1` - Full functionality
	- `4` - Disable both transmitting and receiving RF signals
- `<rst>` - Whether to reset UE
	- `0` - Do not reset the UE before setting it to `<fun>` power level
	- `1` - Reset UE. Device is fully functional after the reset.
## [[Quectel_RG255C_Series_RM255C-GL_AT_Commands_Manual_V1.0.0_Preliminary_20231218.pdf#page=63&selection=24,0,28,18|Get Operator Selection]]
This command returns the current operators and their status, and allows automatic or manual network selection.
```
+COPS: <mode>[,<format>[,<oper>][,<AcT>]]    // Read command result
AT+COPS: <mode>[,<format>[,<oper>[,<AcT>]]    // Write command
```
- `<mode>`
	- `0` - Automatic. Operator selection
	- `1` - Manual operator selection
	- `2` - Deregister from network
	- `3` - Set only `<format>`, and do not attempt registration/deregistration.
	- `4` - Manual/automatic selection.
- `<oper>` - Operator in format as per `<format>`
- `<format>`
	- `0` - Long format alphanumeric `<oper>` which can be up to 16 characters long
	- `1` - Short format alphanumeric `<oper>`
	- `2` - Numeric `<oper>`. GSM location area identification number
- `<AcT>` - Access technology selected.
	- `7` - E-UTRAN
	- `10`- E-UTRAN conneced to a 5GCN
	- `11`- NR connected to 5GCN
	- `12`- NG-RAN
## [[Quectel_RG255C_Series_RM255C-GL_AT_Commands_Manual_V1.0.0_Preliminary_20231218.pdf#page=66&selection=51,0,64,6|Register to Network]]
This command queries the network registration status
```
+C5GREG: <n>,<stat>[,[<tac>],[<ci>],[<AcT>],[<Allowed _NSSAI_length>],[<Allowed_NSSAI>]]    // Read command result
AT+C5GREG=[<n>]    // Write command
```
- `<n>`
	- `0` - Disable network registration unsolicited result code
	- `1` - Enable network registration unsolicited result code `+C5GREG:<stat>`
	- `2` - Enable network registration and location information unsolicited result code `+C5GREG: <stat>[,[<tac>],[<ci>],[<AcT>],[<Allowe d_NSSAI_length>],[<Allowed_NSSAI>]]`
- `<stat>`
	- `0` - Not registered, MT is not currently searching an operator to register to
	- `1` - Registered, home network
	- `2` - Not registered, but MT is currently trying to attach or searching an operator to register to
	- `3` - Registration denied
	- `4` - Unknown
	- `5` - Registered, roaming
	- `8` - Registered for emergency services only
## [[Quectel_RG255C_Series_RM255C-GL_AT_Commands_Manual_V1.0.0_Preliminary_20231218.pdf#page=73&selection=20,0,20,3|Define a Packet Data Protocol (PDP) Context]]
This command specifies PDP context parameters for a specific context `<cid>`. A special form of the Write Command (`AT+CGDCONT=<cid>`) causes the values for context `<cid>` to become undefined. It is not allowed to change the definition of an already activated context.
```
CGDCONT: <cid>,<PDP_type>,<APN>,<PDP_ad dr>,<d_comp>,<h_comp>[,<IPv4AddrAlloc>[,<req uest_type>,,,,,,,,[,<SSC_mode>[,<S-NSSAI>[,<Pref _access_type>,,[,<Always-on_req>]]]]]] […]    // Read command result

AT+CGDCONT=[<cid>[,<PDP_type>[,<APN>[,<PDP_addr>[,<d_comp>[,<h_comp>[,<IPv4AddrAlloc>[,<request_type>,,,,,,,,[,<SSC_mode>[,<S-NSSAI>[,<Pref_access_type>,,[,<Always-on_req>]]]]]]]]]]]]    // Write command
```
- `<cid>` - PDP context identifier. A numeric parameter which specifies a particular PDP context definition. The parameter is local to the TE-MT interface and is used in other PDP context-related commands. The range of supported values (minimum value =1) is returned by the test form of the command. Range: 1-16
- `<PDP_type>` - Packet data protocol type, a string parameter which specifies the type of packet data protocol
	- "IP" - IPv4
	- "PPP" - Point to Point Protocol
	- "IPV6" - IPv6
	- "IPV4V6" - Virtual introduced to handle dual IP stack UE capability
- `<APN>` - Access point name, which is a logical name used to select the GGSN or the external packet data network. If the value is null or omitted, then the subscription value will be requested
- `<PDP_addr>` - Identify the MT in the address space applicable to the PDP. If the value is null or omitted, then a value may be provided by the TE during the PDP startup procedure or, failing that, a dynamic address will be requested. The allocated address may be read using the `AT+CGPADDR`
- `<d_comp>` - Controls PDP data compression (applicable for SNDCP only) 
	- `0` - Off
	- `2`- V.42bis
- `<h_comp>` -  Controls PDP header compression 
	- `0` - Off
	- `4` - RFC3095
- `<IPv4AddrAlloc>` - Control how the MT/TA requests to get the IPv4 address information
	- `0` - IPv4 address allocation through NAS signalling
	- `1` - IPv4 address allocated through DHCP
- `<request_type>` - Indicate the type of PDP context activation request for the PDP context.
	- `0` - PDP context is for new PDP context establishment or for handover from a non-3GPP access network (how the MT decides whether the PDP context is for new PDP context establishment or for handover is implementation specific)
	- `1` - PDP context is for emergency bearer services
- `<SSC_mode>` - Indicate the session and service continuity (SSC) mode for the PDU session in 5GS
	- `0` - SSC mode 1
	- `1` - SSC mode 2
	- `2` - SSC mode 3
- `<S-NSSAI>` - Dependent of the form, the string can be separated by dot(s) and semicolon(s). This parameter is associated with the PDU session for identifying a network slice in 5GS. The parameter has one of the forms:
	- **sst** - only slice/service type (SST) is present
	- **sst;mapped_sst** - SST and mapped configured SST are present
	- **sst.sd** - SST and slice differentiator (SD) are present
	- **sst.sd;mapped_sst** - SST, SD and mapped configured SST are present
	- **sst.sd;mapped_sst.mapped_sd** - SST, SD, mapped configured SST and mapped configured SD are present
- `<Pref_access_type>` - Indicate the preferred access type for the PDU session in 5GS
	- `0` - 3GPP
	- `1` - non-3GPP
- `<Always-on_req>` - Indicate whether the UE requests to establish the PDU session as an always-on PDU session.
	- `0` - always-on not requested
	- `1` - always-on requested
## [[Quectel_RG255C_Series_RM255C-GL_AT_Commands_Manual_V1.0.0_Preliminary_20231218.pdf#page=150&selection=48,0,51,0|(De)Activate PDP Contexts]]
This command activates or deactivates the specified PDP context(s). If any PDP context is already in the requested state, the state for that context remains unchanged.
```
+CGACT: <cid>,<state>    // Read command result
AT+CGACT=<state>[,<cid>]    // Write command
```
- `<state>` - Indicate the state of PDP context activation
	- `0` - Deactivated
	- `1` - Activated
## [[Quectel_RG255C_Series_RM255C-GL_AT_Commands_Manual_V1.0.0_Preliminary_20231218.pdf#page=151&selection=72,0,76,18|Get PDP Address]]
This command returns a list of PDP addresses for the specified context identifiers. If no `<cid>` is specified, the addresses for all defined contexts are returned.
```
+CGPADDR: list of defined <cid>    // Read command result
AT+CGPADDR=[<cid>[,<cid>[,…]]]    // Write command
```
# Steps to create interfaces and PDP contexts
1. To generated new interfaces, first unload the kernel module `qmi_wwan_q` and load it with the desired number of new interfaces.
```
sudo rmod qmi_wwan_q
sudo modprobe qmi_wwan_q qmap_mode=4 # or other number
```

Now in the modem module:
2. Define functionalities - `AT+CFUN=1`
3. Get operator selection mode - `AT+COPS?`
	1. Connect and register to operator - `AT+COPS=<cid>`
4. Get network registration status - `AT+C5GREG?`
5. Get existing PDP Contexts - `AT+CGDCONT?`
	1. Define new PDP context:
		1. `AT+CGDCONT=1,"IPV4V6","backhaul"`
		2. `AT+CGDCONT=4,"IPV4V6","client"`, one for each client
 
**NOTE:** Contexts with `cid` **2** and **3** are reserved for `ims` and `sos`respectively. We are able to redefine the first context from `internet` to `backhaul` but for new contexts in the `clients` DNN we will have to define them it a `cid` starting at 4 or greater. 

7. Back in the machine, start the `quectel_qmi_proxy` and the `quectel-CM` tools while also defining the name of the DNNs to use.
```
sudo ./quectel-qmi-proxy
sudo ./quectel-CM -n 1 -s <dnn1>
sudo ./quectel-CM -n 4 -s <dnn2>
# The -n flag corresponds to the PDP cid configured in the previous step
```