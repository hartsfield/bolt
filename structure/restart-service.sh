#!/bin/bash
git pull
go build -o $1
pkill -f $1
./$1
