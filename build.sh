#!/bin/sh

GOEXPERIMENT=jsonv2 go build -ldflags="-s -w" -o goaurrpc cmd/aur_rpc_service/main.go