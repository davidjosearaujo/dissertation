#!/bin/bash

GNB_IP_UE=$1
CLIENT_EAP_IP=$2
AUTH_SERVER_IP=$3
CLIENT_SECRET=$4
UE_LAN_IP=$5

sudo apt update

echo -e "\nInstalling dependencies for 5G Modem"
sudo apt install -y \
  build-essential \
  linux-headers-$(uname -r) \
  lksctp-tools \
  libsctp-dev \
  iproute2 \
  net-tools \
  python3 \
  dnsmasq \
  yq


sudo snap install yq

echo -e "\nUpdate gNodeB IP"
cat /home/vagrant/open5gs-ue.yaml | yq ".gnbSearchList[0] = \"$GNB_IP_UE\"" | sudo tee /home/vagrant/open5gs-ue.yaml

echo -e "\nEnabling 'backhaul' APN"
cat /home/vagrant/open5gs-ue.yaml \
| yq '.sessions[0].apn = "backhaul"' \
| yq '.sessions[0].apn style="single"' \
| sudo tee /home/vagrant/open5gs-ue.yaml

echo -e "\nRunning UE"
/home/vagrant/nr-ue -c /home/vagrant/open5gs-ue.yaml &> /log/ue.log &

sudo sysctl -w net.ipv4.ip_forward=1

echo -e "Writing hostapd configurations"
echo -e "interface=enp0s9
driver=wired
ctrl_interface=/var/run/hostapd
ctrl_interface_group=vagrant

logger_syslog=-1
logger_syslog_level=0

ieee8021x=1 
own_ip_addr=$CLIENT_EAP_IP

auth_server_addr=$AUTH_SERVER_IP
auth_server_port=1812
auth_server_shared_secret=$CLIENT_SECRET" > /home/vagrant/hostapd.conf

cat hostapd.conf > /log/hostapd.log

echo -e "\nRunning hostapd"
sudo /home/vagrant/hostapd -tKdd /home/vagrant/hostapd.conf &>> /log/hostapd.log &

sudo touch /etc/allowed-macs.conf

echo -e "\nConfiguring dnsmasq DHCP server"
echo -e "interface=enp0s9
bind-interfaces

dhcp-range=192.168.59.11,192.168.59.100,255.255.255.0
dhcp-authoritative

dhcp-option=option:router,$UE_LAN_IP
dhcp-option=option:dns-server,8.8.8.8,8.8.4.4

dhcp-ignore=tag:!known
conf-file=/etc/allowed-macs.conf

log-dhcp" | sudo tee /etc/dnsmasq.conf

sudo systemctl restart dnsmasq

echo -e "\nStart Interceptor"
sudo /home/vagrant/interceptor --mode="debug" --interface="/var/run/hostapd/enp0s9" &>> /log/interceptor.log &