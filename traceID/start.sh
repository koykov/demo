#!/bin/bash

./build.sh

source .env

rm -f logs/*.log
bin/trace-server -port $SERVER_PORT -tport $TRACED_PORT -cbport $CB_PORT -pbport $PB_PORT -cport $CLIENT_PORTS &> logs/server.log & SERVER_PID=$!

IFS=',' read -ra PORTS <<< "$CLIENT_PORTS"
for PORT in "${PORTS[@]}"; do
  bin/trace-client -port $PORT &> "logs/client-$PORT.log" & CLIENT_PID=$!
done
