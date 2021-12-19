#!/usr/bin/env bash

go mod tidy
go build -o $GOPATH/bin/blqueued . &> logs/build.log
LOG="logs/$(date +%Y.%m.%d_%H:%M).log"
$GOPATH/bin/blqueued &> $LOG
