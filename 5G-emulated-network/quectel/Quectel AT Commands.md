##### Sometimes the minicom/serial output is not being sent/received. To fix it type:  `ate`

### Miscellaneous commands
```
AT+CEREG=2;+C5GREG=2  # Something regarding registration
AT+CGDCONT?  
AT+CFUN=1  
AT+COPS=?  
AT+COPS=1,2,99908,12  
AT+QPING=1,"8.8.8.8"
```

### To select a `dnn` through AT commands:
```
AT+CGDCONT=1,"IPV4v6","internet"
```

### List available networks
```
AT+COPS=?
```

### Register in one of the previously listed networks
```
AT+COPS=1,2,99908,12 
```