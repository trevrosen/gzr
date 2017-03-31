#!/usr/bin/env bash

set -eu

source ~/.bash_profile

git config --global url."git@github.com:".insteadOf "https://github.com/"
ssh-keyscan -H github.com >> ~/.ssh/known_hosts
go get github.com/Masterminds/glide

rm -rf $GOPATH/src/github.com/bypasslane/gzr
mkdir -p $GOPATH/src/github.com/bypasslane/gzr
cp -R . $GOPATH/src/github.com/bypasslane/gzr

pushd $GOPATH/src/github.com/bypasslane/gzr
    make install
    make test
popd
