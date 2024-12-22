#!/bin/bash
pkill $1
go build -o $1

# linux
servicePort=$2 logFilePath=./log.txt ./$1 &>> log.txt &

# mac os
# ./$1 >> log.txt 2>&1 &
