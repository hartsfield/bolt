package main

import (
	"net/http"
	"strings"
	"fmt"
)

func viewItem(w http.ResponseWriter, r *http.Request) {
	readDB()
id := strings.Split(r.RequestURI, "/")[2]
	fmt.Println(r.RequestURI, id, itemsMap[id], itemsMap)
	exeTmpl(w, r, &viewData{AppName: itemsMap[id].Message, Stream: []*item{itemsMap[id]}}, "main.html")
}