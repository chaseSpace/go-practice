#!/bin/sh

if ! command -v golangci-lint >/dev/null; then
    echo "golangci-lint not installed or available in the PATH" >&2
    echo "install golangci-lint ..." >&2
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.58.1
fi

#goland 可直接定位文件
golangci-lint run ./... |sed 's/\\/\//g'