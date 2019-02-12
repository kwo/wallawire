#!/bin/sh

BUILD_VERSION=$(git describe --always --tags --dirty='*')
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

CGO_ENABLED=0 go build -o ./wallawire -tags netgo -ldflags "-w -s -X main.Version=${BUILD_VERSION} -X main.BuildTime=${BUILD_TIME}" .
