#!/bin/bash

CORE_IP=$1
GNB_IP_CORE=$2
GNB_IP_UE=$3

cd ~/UERANSIM

echo -e "\nUpdate gNodeB and Core IPs"
cat config/open5gs-gnb.yaml \
| yq ".linkIp = \"$GNB_IP_UE\"" \
| yq ".ngapIp = \"$GNB_IP_CORE\"" \
| yq ".gtpIp = \"$GNB_IP_CORE\"" \
| yq ".amfConfigs[0].address = \"$CORE_IP\"" \
| sudo tee config/open5gs-gnb.yaml

echo -e "\nRunning gNB"
build/nr-gnb -c config/open5gs-gnb.yaml &> /log/$(cat /etc/hostname)_gnb_$(date +%s).log &