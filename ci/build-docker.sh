#!/usr/bin/env bash

set -eu

docker build -f Dockerfile.builder -t tmp/tmp:tmp .
docker run tmp/tmp:tmp /bin/echo .
docker cp $(docker ps -lq):/go/src/github.com/bypasslane/gzr/gzr .
gzr build .
