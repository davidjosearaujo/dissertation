First step is to add your user to the dialout group so that you don't need sudo.

```
sudo usermod -a $USER -G dialout
```

Then, some tools (like pyserial) can output the error of `device busy` or something else. In this specific module, disable ModemManager
```
sudo systemctl stop ModemManager
sudo systemctl disable ModemManager
```

If you just intend to use `minicom` or `screen`, the baudrate is 9600. Just edit the minicom specification accordingly.

For the commands themselves go to [[Quectel AT Commands]]