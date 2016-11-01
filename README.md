# gzr

*"Gozer"*

## What
A CLI tool for managing containerized services and contexts at Bypass, assuming [Kubernetes](http://kubernetes.io) (K8S) as the orchestration tool and [Docker](http://kubernetes.io) as the container runtime.

## Why?
We need an easy way for developers to manage Kubernetes-based resources, whether in a remote environment or locally. Think of `gzr` as sort of like the [Heroku toolbelt](https://blog.heroku.com/the_heroku_toolbelt) for Bypass stuff.

This will replace the other Gozer repo.

## How
In Go, written with the [Cobra CLI framework](https://github.com/spf13/cobra) that powers CLIs for Kubernetes, etcd, Docker etc.