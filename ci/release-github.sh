#!/usr/bin/env bash

set -eu

source ~/.bash_profile
go get github.com/aktau/github-release

TAG=`git describe --tags $(git rev-list --tags --max-count=1)`

pushd release
    tar -czvf gzr-linux-amd64.tar.gz gzr

    github-release release \
        --user bypasslane \
        --repo gzr \
        --tag ${TAG} \

    github-release upload \
        --user bypasslane \
        --repo gzr \
        --tag ${TAG} \
        --name "gzr-linux-amd64" \
        --file gzr-linux-amd64.tar.gz
popd
