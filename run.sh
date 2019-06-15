#!/bin/bash

set -xe

go fmt
go build
set +x
source .env
set -x
./reinze
