#!/bin/bash
pkill -f $1
go build -o $1
servicePort=$2 logFilePath=./logfile.txt ./$1 & 
