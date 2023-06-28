#!/bin/bash
git pull
go build -o $1
pkill -f $1
./$1
# nohup ./$1 > /dev/null & disown
# sleep 2
