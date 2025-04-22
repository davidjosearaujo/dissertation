FROM debian:bookworm-slim AS build
COPY quectel/kernel_drivers/v6.1.1 /kernel_drivers
COPY quectel/qconnect /qconnect
COPY quectel/qmi_wwan /qmi_wwan

RUN sudo apt install -y \
    build-essential \
    net-tools \