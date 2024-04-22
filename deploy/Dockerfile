FROM ubuntu

RUN DEBIAN_NONINTERACTIVE=1 apt update && apt upgrade -y
RUN apt install -y ca-certificates

RUN mkdir -p /etc/groceries

COPY etc /etc/groceries
COPY bin /usr/bin

EXPOSE 8080
ENTRYPOINT ["/usr/bin/grocery-server", "/etc/groceries/manifest.yaml"]