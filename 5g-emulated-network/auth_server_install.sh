#!/bin/bash

CLIENT_EAP_IP=$1
CERT_CA_PASSWD=$2
CERT_SERVER_PASSWD=$3
CERT_CLIENT_PASSWD=$4
CLIENT_SECRET=$5

sudo apt update

sudo apt install -y freeradius

sudo -s -u freerad

echo -e "\nStopping FreeRADIUS service"
sudo systemctl stop freeradius.service

echo -e "\nAdding UE as client"
echo -e "client UE {\n\tipaddr = $CLIENT_EAP_IP\n\tsecret = $CLIENT_SECRET\n}" >> /etc/freeradius/3.0/clients.conf

echo -e "\nRemoving old certificates"
cd /etc/freeradius/3.0/certs
make destroycerts
make clean

echo -e "\nGenerating new certificates"
sed -i "s/input_password.*/input_password = $CERT_CA_PASSWD/" /etc/freeradius/3.0/certs/ca.cnf
sed -i "s/output_password.*/output_password = $CERT_CA_PASSWD/" /etc/freeradius/3.0/certs/ca.cnf
sed -i "s/input_password.*/input_password = $CERT_SERVER_PASSWD/" /etc/freeradius/3.0/certs/server.cnf
sed -i "s/output_password.*/output_password = $CERT_SERVER_PASSWD/" /etc/freeradius/3.0/certs/server.cnf
sed -i "s/input_password.*/input_password = $CERT_CLIENT_PASSWD/" /etc/freeradius/3.0/certs/client.cnf
sed -i "s/output_password.*/output_password = $CERT_CLIENT_PASSWD/" /etc/freeradius/3.0/certs/client.cnf
make

echo -e "\nVerifying certs"
make verify

echo -e "\nEnabling EAP-TLS"
rm /etc/freeradius/3.0/mods-enabled/eap
sed -i "s/default_eap_type = .*/default_eap_type = tls/" /etc/freeradius/3.0/mods-available/eap
sed -i "s/private_key_password = .*/private_key_password = $CERT_SERVER_PASSWD/" /etc/freeradius/3.0/mods-available/eap
sed -i "s#private_key_file = .*#private_key_file = /etc/freeradius/3.0/certs/server.key#" /etc/freeradius/3.0/mods-available/eap
sed -i "s#certificate_file = .*#certificate_file = /etc/freeradius/3.0/certs/server.pem#" /etc/freeradius/3.0/mods-available/eap
sed -i "s#ca_file = .*#ca_file = /etc/freeradius/3.0/certs/ca.pem#" /etc/freeradius/3.0/mods-available/eap
ln -s /etc/freeradius/3.0/mods-available/eap /etc/freeradius/3.0/mods-enabled/eap

chmod 640 /etc/freeradius/3.0/certs/*
chown freerad:freerad /etc/freeradius/3.0/certs/*

echo -e "\nCopying certificates to shared folder"
cp /etc/freeradius/3.0/certs/* /certs

echo -e "\nStarting FreeRADIUS in full debugging"
sudo freeradius -X &>> /log/freeradius.log &

echo -e "\nMaking freeradius start at boot"
echo -e "
sudo freeradius -X &>> /log/freeradius.log &
" | sudo tee /etc/init.d/freeradius > /dev/null
sudo chmod +x /etc/init.d/freeradius