FROM golang:1.8-alpine

ADD . /go/src/github.com/bypasslane/gzr
WORKDIR /go/src/github.com/bypasslane/gzr

RUN apk update && \
    apk add git nodejs make ca-certificates && \
    update-ca-certificates && \
    git config --global url."https://".insteadOf git:// && \
    git config --global url."https://".insteadOf ssh:// && \
    git config --global url."https://github.com/".insteadOf git@github.com: && \
    go get github.com/Masterminds/glide && \
    make install_deps && \
    make install_build_deps && \
    make build_web
