FROM debian:bookworm-slim AS build
COPY quectel/kernel_drivers/v6.1.1 /kernel_drivers
COPY quectel/qconnect /qconnect
COPY quectel/qmi_wwan /qmi_wwan

RUN apt-get update && sudo apt install -y \
    build-essential \
    net-tools \
    minicom

RUN cd /kernel_drivers && \
    make clean && \
    make install

RUN cd /qmi_wwan && \
    make clean && \
    make install

RUN cd /qconnect && \
    make clean && \
    make

RUN usermod -aG $USER dialout