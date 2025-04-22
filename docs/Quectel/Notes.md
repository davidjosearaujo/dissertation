> Quectel UMTS&LTE&5G modules are USB composite devices with multiple USB interfaces. Each USB interface supports different functionalities, which are implemented by loading different USB interface drivers. After a driver is loaded successfully, the corresponding device node is generated, which can be used by the Linux system to implement the module functionalities, such as AT command, GNSS, DIAG, log and USB network adapter. 
> 
> The following table describes the USB interface information of different modules in the Linux system, including USB driver, interface number, device name and interface function. 
> 
> You can obtain the corresponding VID, PID and interface information of the relevant model, and then port the USB interface driver listed in the following table.

[[Quectel_UMTS_LTE_5G_Linux_USB_Driver_User_Guide_V3.2.pdf#page=12&selection=25,0,113,1|Quectel_UMTS_LTE_5G_Linux_USB_Driver_User_Guide_V3.2, page 12]]

# RG255C series/ RM255C-GL
- VID: 0x2c7c
- PID: 0x316

| USB Driver | Interface Number | Device Name         | Function                                                                   |
| ---------- | ---------------- | ------------------- | -------------------------------------------------------------------------- |
| USB serial | **2**            | **/dev/ttyUSB2**    | **AT Command**                                                             |
| QMI_WWAN   | 3                | wwan0 /dev/cdc-wdm0 | Configure the type of USBnet interface as RmNet by **AT+QCFG="usbned",0"** |
| MBIM       | 6 and 7          | wwan0 /dev/cdc-wdm0 | Configure the type of USBnet interface as MBIM by **AT+QCFG="usbnet",2**   |
