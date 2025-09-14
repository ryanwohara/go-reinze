#!/bin/bash

set -xe

go fmt
go build

set +x
source .env
set -x

./go-reinze
