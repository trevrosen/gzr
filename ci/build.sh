#!/usr/bin/env bash

set -eu

source ~/.bash_profile

pushd $GOPATH/src/github.com/bypasslane/gzr
    make build
popd

mkdir -p release
cp $GOPATH/src/github.com/bypasslane/gzr/gzr release
