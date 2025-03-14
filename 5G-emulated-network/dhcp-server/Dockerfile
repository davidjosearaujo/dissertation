FROM debian:12-slim

COPY dhcp-server /dhcp-server

# Install dependencies
RUN apt-get update && apt-get install python3 python3-dev -y \
    git build-essential libnl-3-dev libnl-genl-3-dev \
    && rm -rf /var/lib/apt/lists/*

# Clone and build hostapd
WORKDIR /root
RUN git clone git://w1.fi/hostap.git && \
    cd hostap/hostapd && \
    cp defconfig .config && \
    sed -i "s/#CONFIG_DRIVER_WIRED=y/CONFIG_DRIVER_WIRED=y/" .config && \
    make

WORKDIR /root/hostap/hostapd
ENTRYPOINT ./hostapd -tKdd