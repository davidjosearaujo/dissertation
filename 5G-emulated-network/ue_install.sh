#!/bin/bash

GNB_IP_UE=$1
CLIENT_EAP_IP=$2
AUTH_SERVER_IP=$3
CLIENT_SECRET=$4
UE_LAN_IP=$5

sudo apt update

echo -e "\nInstalling dependencies for 5G Modem"
sudo apt install -y \
  linux-headers-$(uname -r) \
  build-essential \
  lksctp-tools \
  libsctp-dev \
  net-tools \
  iproute2 \
  python3 \
  dnsmasq

sudo snap install yq

echo -e "\nUpdate gNodeB IP"
cat /home/vagrant/open5gs-ue.yaml | yq ".gnbSearchList[0] = \"$GNB_IP_UE\"" | sudo tee /home/vagrant/open5gs-ue.yaml > /dev/null

echo -e "\nEnabling 'backhaul' DNN"
cat /home/vagrant/open5gs-ue.yaml \
| yq '.sessions[0].apn = "backhaul"' \
| yq '.sessions[0].apn style="single"' \
| sudo tee /home/vagrant/open5gs-ue.yaml > /dev/null

echo -e "\nRunning UE"
/home/vagrant/nr-ue -c /home/vagrant/open5gs-ue.yaml &>> /log/ue.log &

echo -e "\nEnabling IP forwarding"
sudo sysctl -w net.ipv4.ip_forward=1

echo -e "\nWriting hostapd configurations"
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

echo -e "\nSetting 'backhaul' DNN channel as default route"
sudo ip route delete default
sudo ip route add default via 10.45.0.1 dev uesimtun0

echo -e "\nConfiguring dnsmasq DHCP server"

sudo touch /etc/allowed-macs.conf

echo -e "interface=enp0s9
bind-interfaces

dhcp-range=192.168.59.11,192.168.59.100,255.255.255.0
dhcp-authoritative

dhcp-option=option:router,$UE_LAN_IP
dhcp-option=option:dns-server,8.8.8.8,8.8.4.4

dhcp-ignore=tag:!known
conf-file=/etc/allowed-macs.conf

log-dhcp" | sudo tee /etc/dnsmasq.conf > /dev/null

sudo systemctl restart dnsmasq

echo -e "\nStart Interceptor"
sudo /home/vagrant/interceptor --mode="debug" --interface="/var/run/hostapd/enp0s9" &>> /log/interceptor.log &

echo -e "\nMaking UE, hostapd and inteceptor start at boot"
echo -e "
sudo /home/vagrant/nr-ue -c /home/vagrant/open5gs-ue.yaml &>> /log/ue.log &
sleep 5
sudo ip route delete default
sudo ip route add default via 10.45.0.1 dev uesimtun0
sudo /home/vagrant/hostapd -tKdd /home/vagrant/hostapd.conf &>> /log/hostapd.log &
sudo /home/vagrant/interceptor --mode="debug" --interface="/var/run/hostapd/enp0s9" &>> /log/interceptor.log &
" | sudo tee /etc/init.d/ue > /dev/null
sudo chmod +x /etc/init.d/ue