#!/usr/bin/env bash

set -eu

source ~/.bash_profile

pushd $GOPATH/src/github.com/bypasslane/gzr
    make install_build_deps
    make install
    make build_web
popd

mkdir -p release
cp $GOPATH/src/github.com/bypasslane/gzr/gzr release
