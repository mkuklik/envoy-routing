# FROM lyft/envoy:latest
# FROM envoyproxy/envoy:latest
FROM envoyproxy/envoy-debug:v1.20-latest
# FROM dwflynn/envoy-debug:20170501

RUN apt-get update && apt-get -q install -y \
    curl \
    dnsutils
COPY envoy.yaml /etc/envoy.yaml
CMD /usr/local/bin/envoy -c /etc/envoy.yaml -l debug 
