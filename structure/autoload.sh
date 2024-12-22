#!/bin/bash
pkill -9 $1
go build -o $1

# linux
servicePort=$2 logFilePath=./logfile.txt ./$1 &>> logfile.txt &

# mac os
# servicePort=$2 logFilePath=./log.txt ./$1 >> log.txt 2>&1 &
