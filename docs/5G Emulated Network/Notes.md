- Different subnets for each PDU Session can be defined in the SMF and UPF.
- In the UPF the interface attributed to each must be defined.
- Apparently, the UPF configuration does not allow for Open5GS to automatically create the tun interfaces.
	- [x] Tun interfaces must be created manually: `ip tuntap add dev intun0 mode tun`
	- [x] Tun interfaces must be brought up manually: `ip link set intun0 up`
	- [x] Open5GS automatically attributes the IP based on the config

# Design a topology diagram for my POC
## PDU Session Topology Perspective
![[topology.png]]
- [ ] Use only two DNNs
	- [ ] One is the default one, for control
	- [ ] Other is for NAUN3s
		- [ ] In this one, the multiple PDU Sessions will be established
# Implement EAP-TLS
## EAP Framework Topology Perspective
![[topology-EAP.png]]
Try to follow instructions from here:
- https://www.tp-link.com/us/support/faq/3456/
- https://wiki.alpinelinux.org/wiki/FreeRadius_EAP-TLS_configuration
- https://simplificandoredes.com/en/freeradius-installation-and-configuration
- https://www.systutorials.com/docs/linux/man/5-wpa_supplicant/ - Usefull for wpa_supplicant in wired mode
- https://w1.fi/wpa_supplicant/devel/
## Testing configurations
We need to stop the running service
``` bash
sudo systemctl stop freeradius.service
```
Then run in debug mode
```bash
sudo freeradius -X
```
## Define UE as a authenticator (client)
- [x] Add new entry to `/etc/freeradius/3.0/clients.conf` file
```bash
$ sudo nano /etc/freeradius/3.0/clients.conf
...
client UE {                  #’UE’ is the alias of your access point
	ipaddr = 192.168.58.100 #The IP address of UE
	secret = testing123     # The ’secret’ will be the ‘Authentication Password’
}
...
```
- [x] It shouldn't be necessary, but if requests are not being received, we may need to configure FreeRADIUS to listen (on `radiusd.conf` file) on the specific interface:
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
- [x] Note that you need to **clean up all the CAs each time before you recreate them**, or `openssl` Swill output ‘Nothing to be done’ and it won’t regenerate new CAs. Delete the existing files by the following command:
```bash
$ rm -f *csr *key *p12 *pem *crl *crt *der *mk *txt *attr *old serial dh
```
- [x] You can edit those \*.cnf files to meet your requirements. 
	- If you wish to change the certificate password, do it in `ca.cnf` in field `output_password`. **ATTENTION**, if you do it, you must also change the password in `mods-available/eap` in the field o `tls > private_key_password`. By default the password should be `whatever`.
- [x] After cleaning up the CAs, run make command to generate new CAs.
```bash
$ make
```
## Enable EAP-TLS as a supported authentication method
- [x] Edit `/etc/freeradius/3.0/mods-available/eap`
```bash
default_eap_type = tls
```
- [x] Delete old and create new symlink (do it as freerad user)
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
- [x] Configuring [`hostapd`](https://wireless.docs.kernel.org/en/latest/en/users/documentation/hostapd.html)
The `hostapd.conf` file will have the following configurations.
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
- [x] We will need to use the `secret` we've defined on `/etc/freeradius/3.0/clients.conf`. In a Vagrantfile can be done by just using a variable.
## Install the certificates on Users
- [x] Copy the generated `ca.pem` and `client.pem` file
> *Note that the password of the private key is ‘whatever’ by default (if you haven’t changed the configurations by editing /etc/freeradius/3.0/certs/\*.cnf).*
## Copy certificate to NAUN3
- [ ] Disable IP on the NAUN3
- [ ] Install `ca.crt`, `client.crt` and `client.key`