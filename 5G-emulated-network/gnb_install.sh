#!/bin/bash

CORE_IP=$1
GNB_IP_CORE=$2
GNB_IP_UE=$3

sudo apt update
sudo apt install -y \
    lksctp-tools \
    libsctp-dev \
    iproute2

sudo snap install yq

echo -e "\nUpdate gNodeB and Core IPs"
cat /home/vagrant/open5gs-gnb.yaml \
| yq ".linkIp = \"$GNB_IP_UE\"" \
| yq ".ngapIp = \"$GNB_IP_CORE\"" \
| yq ".gtpIp = \"$GNB_IP_CORE\"" \
| yq ".amfConfigs[0].address = \"$CORE_IP\"" \
| sudo tee /home/vagrant/open5gs-gnb.yaml > /dev/null

echo -e "\nMaking gNB start at boot"
echo -e "sudo /home/vagrant/nr-gnb -c open5gs-gnb.yaml &>> /log/gnb.log &" | sudo tee /etc/init.d/gnb > /dev/null
sudo chmod +x /etc/init.d/gnb

echo -e "\nRunning gNB"
/home/vagrant/nr-gnb -c open5gs-gnb.yaml &> /log/gnb.log &