#!/bin/bash

go build -o bin/trace-server github.com/koykov/demo/traceID/cmd/server
go build -o bin/trace-client github.com/koykov/demo/traceID/cmd/client
