#!/bin/bash
pkill $1
go build -o $1
case "$OSTYPE" in
  linux*) ./$1 &>> log.txt & ;;
  darwin*) ./$1 >> log.txt 2>&1 & ;;
esac
