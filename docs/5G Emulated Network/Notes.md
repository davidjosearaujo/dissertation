- Different subnets for each PDU Session can be defined in the SMF and UPF.
- In the UPF the interface attributed to each must be defined.
- Apparently, the UPF configuration does not allow for Open5GS to automatically create the tun interfaces.
	- [x] Tun interfaces must be created manually: `ip tuntap add dev intun0 mode tun`
	- [x] Tun interfaces must be brought up manually: `ip link set intun0 up`
	- [x] Open5GS automatically attributes the IP based on the config

# Design a topology diagram for my POC
## PDU Session Topology Perspective
![[topology.png]]
# Implement EAP-TLS
Try to follow instructions from here:
- https://www.tp-link.com/us/support/faq/3456/
- https://wiki.alpinelinux.org/wiki/FreeRadius_EAP-TLS_configuration
- https://simplificandoredes.com/en/freeradius-installation-and-configuration
## Testing configurations
We need to stop the running service
``` bash
sudo systemctl stop freeradius
```
Then run in debug mode
```bash
sudo freeradius -X
```
## Define UE as a authenticator (client)
- [ ] Add new entry to `clients` file
```bash
$ sudo nano /etc/freeradius/clients.conf
client AP1 {                  #’AP1’ is the alias of your access point
	ipaddr = 192.168.0.100/24 #The IP address of UE
	secret = testing123     # The ’secret’ will be the ‘Authentication Password’
}```
- [ ] It shouldn't be necessary, but if requests are not being received, we may need to configure FreeRADIUS to listen on the specific interface:
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
## Enable EAP-TLS as a supported authentication method
- [x] Edit `etc/freeradius/mods-enabled/eap`
```bash
default _eap_type = tls
```
## Make the certificates
```bash
$ sudo -s
$ cd /etc/freeradius/certs
```
Note that you need to **clean up all the CAs each time before you recreate them**, or `openssl` Swill output ‘Nothing to be done’ and it won’t regenerate new CAs. Delete the existing files by the following command:
```bash
$ rm -f *csr *key *p12 *pem *crl *crt *der *mk *txt *attr *old serial dh
```
You can edit those \*.cnf files to meet your requirements. Here we just leave them all to default for testing purpose. After cleaning up the CAs, run make command to generate new CAs.
```bash
$ make
```
## Restart the FreeRADIUS server
```bash
$ sudo systemctl restart freeradius
```
## Install the certificates on Users
Copy the generated ca.der and client.p12 file
Install ca.der. Then, install client.p12. Note that the password of the private key is ‘_whatever_’ by default (if you haven’t changed the configurations by editing /etc/freeradius/certs/\*.cnf).
## Define NAUN3 as a suplicant (user)
- [ ] Add new entry to `users` file, with support form EAP-TLS authentication method