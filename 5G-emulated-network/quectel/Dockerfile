FROM debian:bookworm-slim AS build
COPY quectel/kernel_drivers/v6.1.1 /kernel_drivers
COPY quectel/qconnect /qconnect
COPY quectel/qmi_wwan /qmi_wwan

RUN apt-get update && sudo apt install -y \
    build-essential \
    net-tools 
    
RUN cd /kernel_drivers && \
    make clean && \
    make install

RUN cd qmi_wwan && \
    make clean && \
    make install

RUN cd qconnect && \
    make clean && \
    make

FROM debian:bookworm-slim
COPY --from=build /kernel_drivers /kernel_drivers
COPY --from=build /qconnect/quectel-qmi-proxy /qconnect/quectel-qmi-proxy
COPY --from=build /qmi_wwan /qmi_wwan