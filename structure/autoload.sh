#!/bin/bash
servicePort=5654
logFilePath=./logfile.txt
pkill -9 $1
go build -o $1
servicePort=5652 logFilePath=./logfile.txt ./$1 & 
