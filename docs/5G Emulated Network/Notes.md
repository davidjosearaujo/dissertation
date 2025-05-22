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
# MAC-Based Outbound Traffic for a Specific LAN Device
This guide explains how to configure your Linux router to allow a specific LAN device, identified by its MAC address, to access the internet (WAN). This setup uses a MAC-based firewall rule, so it doesn't require static IP assignment or dynamic IP updates for this specific outbound rule. Port forwarding is not covered in this guide.

```bash
# --- START OF COMMANDS ---

# === 1. Enable IP Forwarding ===
sudo sysctl -w net.ipv4.ip_forward=1

# === 2. Tag Traffic Based on MAC Address ===

# Mark packets from Device 1 (MAC1) to route via WAN1
sudo iptables -t mangle -A PREROUTING -i <LAN_IF> -m mac --mac-source <MAC_DEVICE> -j MARK --set-mark <PDU_SESSION_ID>

# === 3. Create Custom Routing Tables ===
echo "<200+PDU_SESSION_ID> pdu1route" | sudo tee -a /etc/iproute2/rt_tables

# Add routes to custom tables
sudo ip route add default via <PDU_GATEWAY_IP> dev <PDU_IF> table pdu<PDU_SESSION_ID>route

# Create rules to use custom routing tables
sudo ip rule add fwmark <PDU_SESSION_ID> table pdu<PDU_SESSION_ID>route

# === 4. NAT Rules ===
# MASQUERADE for both WANs
sudo iptables -t nat -A POSTROUTING -o <PDU_IF> -j MASQUERADE

# === 5. Firewall Rules (Optional) ===
# Drop all forwarding by default (optional and careful)
sudo iptables -P FORWARD DROP

# Allow related/established traffic
sudo iptables -A FORWARD -m state --state RELATED,ESTABLISHED -j ACCEPT

# Allow forwarding from LAN to each WAN based on MAC
sudo iptables -A FORWARD -i <LAN_IF> -o <PDU_IF> -m mac --mac-source <MAC_DEVICE> -j ACCEPT

# --- END OF COMMANDS ---

# === 6. Persist the Configuration ===
# Debian/Ubuntu systems:
sudo apt install iptables-persistent
sudo netfilter-persistent save
```