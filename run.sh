#!/bin/bash

set -xe

go fmt
go build
source .env
./reinze
