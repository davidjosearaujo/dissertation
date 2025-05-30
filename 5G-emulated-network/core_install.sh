#!/bin/bash

CORE_IP=$1
UE_IMSI=$2
UE_KEY=$3
UE_OPC=$4

sudo add-apt-repository ppa:open5gs/latest
sudo apt-get update && apt-get install -y \
    linux-headers-$(uname -r) \
    build-essential \
    libqmi-glib-dev \
    net-tools \
    minicom \
    iperf3 \
    gnupg \
    curl

sudo snap install yq

echo -e "\nAdding MongoDB repos"
curl -fsSL https://www.mongodb.org/static/pgp/server-8.0.asc | sudo gpg -o /usr/share/keyrings/mongodb-server-8.0.gpg --dearmor

echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-8.0.gpg ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/8.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-8.0.list

echo -e "\nAdding Node repos"
curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg

NODE_MAJOR=20
echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | sudo tee /etc/apt/sources.list.d/nodesource.list

sudo apt-get update

echo -e "\nInstalling MongoDB, Nginx and Open5GS"
sudo apt-get install -y mongodb-org open5gs nodejs nginx

echo -e "\nReload MongoDB"
sudo systemctl start mongod
sudo systemctl enable mongod

echo -e "\nUpdate AMF NGAP and UPF GTPU IPs"
cat /etc/open5gs/amf.yaml | yq ".amf.ngap.server[0].address = \"$CORE_IP\"" | sudo tee /etc/open5gs/amf.yaml
cat /etc/open5gs/upf.yaml | yq ".upf.gtpu.server[0].address = \"$CORE_IP\"" | sudo tee /etc/open5gs/upf.yaml

echo -e "\nAttributing DNN name to default APN"
cat /etc/open5gs/smf.yaml \
| yq '.smf.session[0].dnn = "backhaul"' \
| sudo tee /etc/open5gs/smf.yaml
cat /etc/open5gs/upf.yaml \
| yq '.upf.session[0].dnn = "backhaul"' \
| yq '.upf.session[0].dev = "ogstun"' \
| sudo tee /etc/open5gs/upf.yaml

echo -e "\nAdding DNN for NAUN3 PDU Sessions attribution"
sudo ip tuntap add dev clientun0 mode tun
sudo ip link set clientun0 up
sudo ip addr add 10.46.0.1/24 brd 10.46.0.255 dev clientun0
cat /etc/open5gs/smf.yaml \
| yq '.smf.session[1].subnet = "10.46.0.0/24"' \
| yq '.smf.session[1].gateway = "10.46.0.1"' \
| yq '.smf.session[1].dnn = "clients"' \
| sudo tee /etc/open5gs/smf.yaml
cat /etc/open5gs/upf.yaml \
| yq '.upf.session[1].subnet = "10.46.0.0/24"' \
| yq '.upf.session[1].gateway = "10.46.0.1"' \
| yq '.upf.session[1].dnn = "clients"' \
| yq '.upf.session[1].dev = "clientun0"' \
| sudo tee /etc/open5gs/upf.yaml

sudo systemctl restart open5gs-upfd
sudo systemctl restart open5gs-smfd 
sudo systemctl restart open5gs-amfd

echo -e "\nInstalling WebUI"
curl -fsSL https://open5gs.org/open5gs/assets/webui/install | sudo -E bash -

echo -e "\nExposing WebUI through a reverse proxy"
echo "server {
    listen 8080;

    location / {
        proxy_pass http://127.0.0.1:9999;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
    }
}" | sudo tee /etc/nginx/sites-available/webui

sudo ln -s /etc/nginx/sites-available/webui /etc/nginx/sites-enabled/
sudo systemctl restart nginx

echo -e "\nPre-loading UE in the AMF"
sudo chmod +x open5gs-dbctl
./open5gs-dbctl add $UE_IMSI $UE_KEY $UE_OPC

echo -e "\nPre-loading UE APN sessions"
./open5gs-dbctl update_apn $UE_IMSI backhaul 0
./open5gs-dbctl update_apn $UE_IMSI clients 0