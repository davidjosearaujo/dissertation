docker run -it --rm -v "$(pwd)":/home --entrypoint=/bin/bash debian -c "
    apt update && \
    apt install -y make lsb-release && \
    cd /home && \
    make setup && \
    make && \ 
    make print"