#!/bin/bash
export servicePort=9394
export logFilePath=./log_file.txt
trap -- '' SIGTERM
git pull
go build -o boiler_plate
pkill -f boiler_plate
nohup ./boiler_plate > /dev/null & disown
sleep 2
