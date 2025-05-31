#!/bin/bash

CERT_CLIENT_PASSWD=$1

sudo apt-get update
sudo apt-get install -y wpasupplicant iperf3

echo -e "\nRemove IP address"
sudo ip addr flush enp0s8

echo -e "\nWriting wpa_supplicant configurations"
echo -e 'ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=vagrant
eapol_version=2
ap_scan=0

network={
    eapol_flags=0

    key_mgmt=IEEE8021X
    eap=TLS

    identity="user@example.org"
    ca_cert="/certs/ca.pem"
    private_key="/certs/client.p12"
    private_key_passwd="'$CERT_CLIENT_PASSWD'"
}' > wpa_supplicant.conf

LOG_FILE_PATH="/log/$(cat /etc/hostname).log"
cat wpa_supplicant.conf > ${LOG_FILE_PATH}

echo -e "\nMaking wpa_supplicant start at boot"
echo -e "
sudo wpa_supplicant -tKdd -ienp0s8 -Dwired -c./wpa_supplicant.conf &>> ${LOG_FILE_PATH} &
" | sudo tee /etc/init.d/wpa_supplicant > /dev/null
sudo chmod +x /etc/init.d/wpa_supplicant

for i in $(seq 1 10); do
    startTime=$(date +%s)

    sudo wpa_supplicant -tKdd -ienp0s8 -Dwired -c./wpa_supplicant.conf &>> ${LOG_FILE_PATH} &

    sudo dhclient enp0s8

    elapsedTime=$(($(date +%s) - startTime))
    echo -e "
    Attempt #$i elapsed time: $elapsedTime seconds
    " | sudo tee -a /log/$(cat /etc/hostname)_connection_delay.log

    sudo killall wpa_supplicant
    sudo ip addr flush enp0s8

    sleep 20
done
