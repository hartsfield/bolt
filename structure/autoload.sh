#!/bin/bash
pkill -9 $1
go build -o $1
servicePort=$2 logFilePath=./logfile.txt ./$1 &>> logfile.txt &
