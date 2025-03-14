
import os
import sys
import time
import ctrl

wpas_ctrl_dir = '/var/run/hostapd'

def connect(host, port=9877):
    if host != None:
        try:
            htpd = ctrl.Ctrl(host, port)
            return htpd
        except:
            print("Could not connect to host: ", host)
            return None

def main(host=None, port=9877):
    # Connecting to socket 
    mon = connect(host, port)
    if mon is None:
        print("Could not open event monitor connection")
        return

    # Attaching to socket file
    mon.attach()

    # Receiving events
    while True:
        while mon.pending():
            ev = mon.recv()
            print(ev)


if __name__ == "__main__":
    main(host=sys.argv[1], port=int(sys.argv[2]))