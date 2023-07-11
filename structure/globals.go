package main

import (
	"html/template"
	"os"
)

// ckey/ctxkey is used as the key for the HTML context and is how we retrieve
// token information and pass it around to handlers
type ckey int

const ctxkey ckey = iota

var (
	servicePort                    = ":" + os.Getenv("servicePort")
	logFilePath                    = os.Getenv("logFilePath")
	templates   *template.Template = template.New("main")
	companyName string             = "BoltApp"
)
