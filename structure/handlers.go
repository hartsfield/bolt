package main

import (
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	exeTmpl(w, r, nil, "main.tmpl")
}
