#!/bin/sh

set -e

# TODO(trevors): Auto-detect CrOS or other Linux
URL_LAUNCHER=/usr/bin/garcon-url-handler

go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
$URL_LAUNCHER coverage.html
