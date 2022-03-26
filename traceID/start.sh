#!/bin/bash

source .env

rm -f log/*.log

bin/trace-server -port $SERVER_PORT -tport $TRACED_PORT -cbport $CB_PORT -pbport $PB_PORT -cport $CLIENT_PORTS &> log/server.log & SERVER_PID=$!

IFS=',' read -ra PORTS <<< "$CLIENT_PORTS"
for PORT in "${PORTS[@]}"; do
  bin/trace-client -port $PORT &> "log/client-$PORT.log" & CLIENT_PID=$!
done
