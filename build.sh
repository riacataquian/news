#!/usr/bin/env bash

set -e

vgo test ./...
vgo test -race ./...
vgo vet ./...
vgo fmt ./...
golint ./...
