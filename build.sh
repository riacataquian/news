#!/usr/bin/env bash

set -e

env GO111MODULE=on vgo test -race ./...
env GO111MODULE=on golint ./...
