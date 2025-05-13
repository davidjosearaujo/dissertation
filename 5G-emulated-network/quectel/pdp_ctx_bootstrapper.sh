# Copyright 2025 David Araújo
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     https://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Testing creating multiple PDP context in the same DNN
sudo qmicli -d /dev/cdc-wdm0 --devcie-open-qmi --wds-create-profile="3gpp,name=naun3_1,apn=client,pdp-type=IPV4V6,auth=NONE"
sudo qmicli -d /dev/cdc-wdm0 --device-opne-qmi --wds-create-profile="3gpp,name=naun3_2,apn=client,pdp-type=IPV4V6,auth=NONE"

echo Y | sudo tee /sys/class/net/wwan0/qmi/raw_ip

echo 1 | sudo tee /sys/class/net/wwan0/qmi/add_mux	# backhaul
echo 2 | sudo tee /sys/class/net/wwan0/qmi/add_mux	# NAUN3 1
echo 3 | sudo tee /sys/class/net/wwan0/qmi/add_mux	# NAUN3 2

sudo ip link set wwan0 up
sudo ip link set qmimux0 up

sudo qmicli -p -d /dev/cdc-wdm0 --device-open-qmi --client-no-release-cid --wds-noop
sudo qmicli -p -d /dev/cdc-wdm0 --device-open-qmi --client-no-release-cid --wda-noop

sudo qmicli -p -d /dev/cdc-wdm0 --device-open-qmi --wda-set-data-format="link-layer-protocol=raw-ip,ul-protocol=qmap,dl-protocol=qmap,dl-max-datagrams=32,dl-datagram-max-size=32768,ep-type=hsusb,ep-iface-number=4" --client-cid=$wda_cid --client-no-release-cid

# mux-id shall by the ID returned by the previous command
sudo qmicli -p -d /dev/cdc-wdm0 --device-open-qmi --wds-bind-mux-data-port="mux-id=$id,ep-iface-number=4" --client-cid=$wds_cid --client-no-release-cid

sudo qmicli -p -d /dev/cdc-wdm0 --device-open-qmi --wds-start-network="3gpp-profile=5,apn=client,ip-type=4" --wds-follow-network --client-cid=$wds_cid --client-no-release-cid

sudo qmicli -p -d /dev/cdc-wdm0 --device-open-qmi --wds-get-current-settings