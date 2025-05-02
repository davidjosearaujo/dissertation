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
# PDP Context
## [[Quectel_RG255C_Series_RM255C-GL_AT_Commands_Manual_V1.0.0_Preliminary_20231218.pdf#page=63&selection=24,0,28,18|5.1 AT+COPS Operator Selection]]
This command returns the current operators and their status, and allows automatic or manual network selection.
```
AT+COPS=?
+COPS: (3,"268 02","268 02","26802",7),(3,"NOS","NOS","26803",7),(3,"vodafone P","voda P","26801",7),(3,"MEO","MEO","26806",7)
```
`+COPS: <mode>[,<format>[,<oper>][,<AcT>]]`
- `<mode>`
	- **0** - Automatic. Operator selection
	- **1** - Manual operator selection
	- **2** - Deregister from network
	- **3** - Set only `<format>`, and do not attempt registration/deregistration.
	- **4** - Manual/automatic selection.
## [[Quectel_RG255C_Series_RM255C-GL_AT_Commands_Manual_V1.0.0_Preliminary_20231218.pdf#page=65&selection=94,0,98,31|5.2 AT+CEREG EPS Network Registration Status]]
This command queries the network registration status
```
AT+CEREG?
+CEREG: <n>,<stat>[,<tac>,<ci>[,<AcT]]
```
- `<n>`
	- **0** - Disable network registration unsolicited result code
	- **1** - Enable network registration unsolicited result code `+CEREG:<stat>`
	- **2** - Enable network registration and location information unsolicited result code `+CEREG: <stat>[,[<lac>],[<ci>],[<AcT>]]`
- `<stat>`
	- **0** - Not registered, MT is not currently searching an operator to register to
	- **1** - Registered, home network
	- **2** - Not registered, but MT is currently trying to attach or searching an operator to register to
	- **3** - Registration denied
	- **4** - Unknown
	- **5** - Registered, roaming

## [[Quectel_RG255C_Series_RM255C-GL_AT_Commands_Manual_V1.0.0_Preliminary_20231218.pdf#page=73&selection=20,0,20,3|5.6 AT+CGDCONT Define PDP Contexts]]
This command specifies PDP context parameters for a specific context
## [[Quectel_RG255C_Series_RM255C-GL_AT_Commands_Manual_V1.0.0_Preliminary_20231218.pdf#page=151&selection=72,0,76,18|9.3 AT+CGPADDR Show PDP Addresses]]
- Defining a PDP context: `AT+CGDCONT=1,"IP","UNINET"`
- Activating the PDP: `AT+CGACT=1,1`
- Showing the PDP address: `AT+CGPADDR=1`
