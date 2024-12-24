#!/bin/bash
pkill $1 || true
go build -o $1
./$1 >> log.txt 2>&1 &
