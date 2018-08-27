FROM debian:jessie-slim

LABEL Description="goproxy docker image" Vendor="Tyk" Version=$TYKVERSION

RUN apt-get update \
 && apt-get upgrade -y \
 && apt-get install -y --no-install-recommends \
            curl ca-certificates apt-transport-https \
 && curl -L https://packagecloud.io/tyk/tyk-pump/gpgkey | apt-key add - \
 && apt-get autoremove -y \
 && rm -rf /root/.cache


RUN mkdir -p /opt/goproxy
COPY goproxy /opt/goproxy/goproxy
COPY ingress.sh /opt/goproxy/ingress.sh

RUN chmod +x /opt/goproxy/goproxy
RUN chmod +x /opt/goproxy/ingress.sh

VOLUME ["/opt/goproxy"]

WORKDIR /opt/goproxy

CMD ["/opt/goproxy/ingress.sh"]
