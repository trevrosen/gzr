FROM alpine:latest

ADD gzr /usr/local/bin/gzr
ADD docker-config/kube.conf /root/.kube/config
ADD docker-config/gzr.json /root/.gzr.json
ADD start.sh /usr/local/bin/start.sh
ADD https://storage.googleapis.com/kubernetes-release/release/v1.5.3/bin/linux/amd64/kubectl /usr/local/bin/kubectl

RUN apk update && \
    apk add bash && \
    chmod +x /usr/local/bin/kubectl

ENTRYPOINT ["start.sh"]
