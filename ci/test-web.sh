#!/usr/bin/env bash

set -eu

source ~/.bash_profile

pushd gozer-web
    npm install
    # tests don't exist yet, uncomment when they do @JonathanTech
    #npm test
popd
