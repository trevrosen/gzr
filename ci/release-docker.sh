#!/usr/bin/env bash

set -eu

pushd gozer-web
    # very inception
    gzr build .
popd
