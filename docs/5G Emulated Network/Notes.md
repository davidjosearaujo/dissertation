- Different subnets for each PDU Session can be defined in the SMF and UPF.
- In the UPF the interface attributed to each must be defined.
- Apparently, the UPF configuration does not allow for Open5GS to automatically create the tun interfaces.
	- [x] Tun interfaces must be created manually: `ip tuntap add dev intun0 mode tun`
	- [x] Tun interfaces must be brought up manually: `ip link set intun0 up`
	- [x] Open5GS automatically attributes the IP based on the config

# Design a topology diagram for my POC
## PDU Session Topology Perspective
![[topology.png]]
- [x] Use only two DNNs
	- [x] One is the default one, for control
	- [x] Other is for NAUN3s
		- [x] In this one, the multiple PDU Sessions will be established
I've created an [issue](https://github.com/aligungr/UERANSIM/issues/756) in UERANSIM GitHub repo regarding the faillure in creating the multiple PDU Sessions.
# Implement EAP-TLS with FreeRADIUS and `hostapd` 
## EAP Framework Topology Perspective
![[topology_EAP.png]]
Try to follow instructions from here:
- [Configuration Guide on EAP-TLS authentication for WPA-Enterprise (with FreeRADIUS)](https://www.tp-link.com/us/support/faq/3456/)
- [Alpine FreeRadius EAP-TLS configuration](https://wiki.alpinelinux.org/wiki/FreeRadius_EAP-TLS_configuration)
- [FreeRadius: Installation and Configuration](https://simplificandoredes.com/en/freeradius-installation-and-configuration)
- [wpa_supplicant (5) - Linux Manuals](https://www.systutorials.com/docs/linux/man/5-wpa_supplicant/) - Useful for `wpa_supplicant` in wired mode
- [`wpa_supplicant` Documentation](https://w1.fi/wpa_supplicant/devel/)
## Testing FreeRADIUS configurations
We need to stop the running service
``` bash
sudo systemctl stop freeradius.service
```
Then run in debug mode
```bash
sudo freeradius -X
```
## Define UE as a authenticator (client)
- Add new entry to `/etc/freeradius/3.0/clients.conf` file
```bash
$ sudo nano /etc/freeradius/3.0/clients.conf
...
client UE {                  #’UE’ is the alias of your access point
	ipaddr = 192.168.58.100 #The IP address of UE
	secret = testing123     # The ’secret’ will be the ‘Authentication Password’
}
...
```
- It shouldn't be necessary, but if requests are not being received, we may need to configure FreeRADIUS to listen (on `radiusd.conf` file) on the specific interface:
```bash
listen {
	type = auth
	ipaddr = *
	port = 1812
	interface = eth0
}
```
or we can just allow it to bind to all interfaces:
```bash
listen {
  type = auth
  ipaddr = *
  port = 1812
}
```
## Make the certificates
```bash
$ sudo -s freerad
$ cd /etc/freeradius/3.0/certs
```
- Note that you need to **clean up all the CAs each time before you recreate them**, or `openssl` Swill output ‘Nothing to be done’ and it won’t regenerate new CAs. Delete the existing files by the following command:
```bash
$ rm -f *csr *key *p12 *pem *crl *crt *der *mk *txt *attr *old serial dh
```
- You can edit those `\*.cnf` files to meet your requirements. 
	- If you wish to change the certificate password, do it in `ca.cnf` in field `output_password`. **ATTENTION**, if you do it, you must also change the password in `mods-available/eap` in the field o `tls > private_key_password`. By default the password should be `whatever`.
- After cleaning up the CAs, run make command to generate new CAs.
```bash
$ make
```
## Enable EAP-TLS as a supported authentication method
- Edit `/etc/freeradius/3.0/mods-available/eap`
```bash
default_eap_type = tls
```
- Delete old and create new symlink (do it as freerad user)
```bash
$ sudo -s -u freerad
$ rm /etc/freeradius/3.0/mods-enabled/eap
$ ln -s /etc/freeradius/3.0/mods-available/eap /etc/freeradius/3.0/mods-enabled/eap
```
## Restart the FreeRADIUS server
```bash
$ sudo systemctl restart freeradius
```
## Config the UE as AP with client secret
### Configuring [`hostapd`](https://wireless.docs.kernel.org/en/latest/en/users/documentation/hostapd.html)
- The `hostapd.conf` file will have the following configurations.
```
interface=enp0s10
logger_syslog=-1
driver=wired  # THIS NEEDS TO BE CHANGED TO USE 80211nl DRIVER
#IEEE 802.11 Configs
#ssid=ue_ap_1
#channel=1
own_ip_addr="UE_EAP_IP"
auth_server_addr="AUTH_SERVER_IP"
auth_server_port=1812
auth_server_shared_secret="CLIENT_SECRET"
```
- We will need to use the `secret` we've defined on `/etc/freeradius/3.0/clients.conf`. In a Vagrantfile can be done by just using a variable.
## Install the certificates on Users
-  Copy the generated `ca.pem` and `client.pem` file
> *Note that the password of the private key is ‘whatever’ by default (if you haven’t changed the configurations by editing /etc/freeradius/3.0/certs/\*.cnf).*
## Copy certificate to NAUN3
- Disable IP on the NAUN3
- Install `ca.crt`, `client.crt` and `client.key`
# DHCP for authenticated hosts
Create my on DHCP server that can check if devices have been successfully authenticated before leasing an IP.
- Hostapd flags [docs](https://w1.fi/wpa_supplicant/devel/wpa__ctrl_8h.html)
- I need to access `hostapd`to check authenticated devices. Check [this](https://w1.fi/wpa_supplicant/devel/hostapd_ctrl_iface_page.html) documentation.
- Needs to deauth users and request and close PDU Sessions
# Establishing Multiple PDU Sessions/PDP Contexts/PDNs with Quectel Cellular Modem
## How to use QMAP
### With bridge and via Quectel-CM
1. Enable `QUECTEL_BRIDGE_MODE` in `qmi_wwan_q.c`, line 134
``` C
#define QUECTEL_BRIDGE_MODE
```

2. Set `qmap_mode`to 4
```C
#define QUECTEL_WWAN_QMAP 4
```

3. Add direct interface to bridge mapping.
``` C
#ifdef QUECTEL_BRIDGE_MODE
static uint __read_mostly bridge_mode = BIT(1)|BIT(2);
module_param( bridge_mode, uint, S_IRUGO );
#endif
```

4. Compile with `make install`

5. Load module to kernel with `sudo modprobe qmi_wwan_q qmap_mode=4 bridge_mode=6`

6. Activate bridge
```bash
sudo ip link add name lan_bridge type bridge
sudo ip link set dev <wwan0_idx> master br maslan_bridge
```

7. Check `bridge`interfaces
```bash
ip link show type bridge
```

8. Check interfaces connected to the bridge
```bash
sudo bridge link show br lan_brige
```

9. Create new PDP Context if needed
```bash
sudo qmicli -d /dev/cdc-wdm0 --device-open-qmi --wds-create-profile="3gpp,name=naun3_1,apn=clients,pdp-type=IPV4V6,auth=NONE"
```

10. Active QMI proxy
```bash
./quectel-qmi-proxy -d /dev/cdc-wdm0
```

11. Use `quectel-CM` to setup data call with proper PDN and interface binding
```bash
./quectel-CM -n 1 -m 2 -s internet # backhaul DNN
./quectel-CM -n 4 -m 3 -s client
```

Flags:
- `-b` enables network interface bridge function
- `-n` specifies which PDN to setup data call;
- `-m` binds a QMI data call to `wwan0_<iface_idx>` when QMAP is used. E.g  `-n 1 -m 1`, it binds the PDN 1 to `wwan0_1` .
- `-s` flag allows us to specify which APN to connect to.

12. Get IP address **via** DHCP
```bash
udhcpc -i br2
udhcpc -i br3
...
``` 

13. If the QMI data call is left running in the background, you can later kill the connection, by **specifying the PDN ID number**
```bash
./quectel-CM -k 1
```