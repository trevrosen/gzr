#!/usr/bin/env bash

set -eu

# Start kube proxy for accessing local kube, ignore if fail.
kubectl proxy --port=8080 & 2>/dev/null

gzr web --namespace=integration
