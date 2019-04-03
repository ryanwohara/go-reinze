#!/bin/bash

set -xe

go fmt *.go
go build *.go
./main
