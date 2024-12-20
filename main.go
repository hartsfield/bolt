package main

import (
	_ "embed"
)

func main() {
	rc = readConf()
	// see init.go
	readFlags()
}
