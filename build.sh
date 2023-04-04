#!/usr/bin/env bash

export GOPROXY=https://goproxy.cn,direct
export GO111MODULE=on

go mod tidy
go mod vendor #  将下载的包保存到  vendor/ 才能正常引用