FROM ubuntu

RUN DEBIAN_NONINTERACTIVE=1 apt update && apt upgrade -y
RUN apt install -y ca-certificates

COPY root/ /

EXPOSE 3000
ENTRYPOINT ["/usr/bin/grocery-price-fetcher", "/etc/grocery-price-fetcher/config.yaml"]