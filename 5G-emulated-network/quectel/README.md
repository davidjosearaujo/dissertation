# Setup and installation

After uploading the `quectel` folder to any Ubuntu system, you should have 3 directories.
	- kernel_drivers
	- qcnonnect
	- qmi_wwan

### Steps:
Dependencies: make (`sudo apt install build-essential net-tools`)

1) Install `usb-serial` option drivers.
	-`cd kernel_drivers` and `cd` again to the closest version to your current kernel version (`uname -r`)
	- do `sudo make install`

2) Install the `wwan` drivers
	- `cd qmi_wwan`
	- do `sudo make install`

3) Install QConnect Manager
	- `cd qconnect`
	- do `sudo make`

4) Reboot the system (`sudo reboot now`)

5) Connect the board to the APU, turn it ON and press the PWRKEY
	- You can check its connection status with `sudo dmesg -w` to see when it connects.

6) After it successfully connects:
	- `cd qconnect`
	- `sudo ./quectel-CM`

7) Check it is working by pinging 1.1.1.1

#### Note:
On some APUs, for some reason, ttyUSB* serial devices won't show up. If you need them:
- `sudo modprobe usb_wwan`
- Plug in the 5G module and wait until it connects. *AFTER* perfrom the below command
- `sudo insmod /lib/modules/6.8.0-52-generic/kernel/drivers/usb/serial/option.ko`


##### DEBUG
I don't know why, but `modprobe` loads the compressed (.ko.zst) `option` module, while the `insmod` loads the uncompressed file (.ko).
I've narrowed down the problem to the compression of that module, but I don't know how to fix it, nor if there is even a point.
