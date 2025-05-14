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

# Get initial client IDs for configuration



wds_cid=$(sudo qmicli -d /dev/cdc-wdm0 --client-no-release-cid --wds-noop | awk -F 'CID:' '{print $2}' | cut -c 3-4)
wda_cid=$(sudo qmicli -d /dev/cdc-wdm0 --client-no-release-cid --wda-noop | awk -F 'CID:' '{print $2}' | cut -c 3-3)

# Testing creating multiple PDP context in the same DNN
sudo qmicli -d /dev/cdc-wdm0 --wds-create-profile="3gpp,name=naun3_1,apn=client,pdp-type=IPV4V6,auth=NONE" --client-cid=$wds_cid --client-no-release-cid
sudo qmicli -d /dev/cdc-wdm0 --wds-create-profile="3gpp,name=naun3_2,apn=client,pdp-type=IPV4V6,auth=NONE" --client-cid=$wds_cid --client-no-release-cid
sudo qmicli -d /dev/cdc-wdm0 --wds-get-profile-list=3gpp --client-cid=$wds_cid --client-no-release-cid

# Enable raw ip encapsulation
sudo qmicli -d /dev/cdc-wdm0 --set-expected-data-format=raw-ip

# Using fixed number of interfaces
sudo qmicli -d /dev/cdc-wdm0 --link-add="iface=wwan0,prefix=backhaul,mux-id=1"
sudo qmicli -d /dev/cdc-wdm0 --link-add="iface=wwan0,prefix=naun3,mux-id=2"
sudo qmicli -d /dev/cdc-wdm0 --link-add="iface=wwan0,prefix=naun3,mux-id=3"

# Check the links
sudo qmicli -d /dev/cdc-wdm0 --link-list=wwan0

# Bring up the link layer
sudo ip link set wwan0 up
sudo ip link set qmimux0 up
sudo ip link set qmimux1 up
sudo ip link set qmimux2 up

# Repeat the following commands for each mux-id

sudo qmicli -p -d /dev/cdc-wdm0 --client-no-release-cid --wds-noop
sudo qmicli -p -d /dev/cdc-wdm0 --client-no-release-cid --wda-noop

sudo qmicli -p -d /dev/cdc-wdm0 --wda-set-data-format="link-layer-protocol=raw-ip,ul-protocol=qmap,dl-protocol=qmap,dl-max-datagrams=32,dl-datagram-max-size=32768,ep-type=hsusb,ep-iface-number=4" --client-cid=$wda_cid --client-no-release-cid

## mux-id shall by the ID returned by the previous command
sudo qmicli -p -d /dev/cdc-wdm0 --wds-bind-mux-data-port="mux-id=$id,ep-iface-number=4" --client-cid=$wds_cid --client-no-release-cid

sudo qmicli -p -d /dev/cdc-wdm0 --wds-start-network="3gpp-profile=4,apn=client,ip-type=4" --wds-follow-network --client-cid=$wds_cid --client-no-release-cid

sudo qmicli -p -d /dev/cdc-wdm0  --wds-get-current-settings