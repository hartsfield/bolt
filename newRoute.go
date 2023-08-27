package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func newRoute(routeAndHandler string) {
	filePath := "router.go"
	var route, handler = strings.Split(routeAndHandler, ",")[0], strings.Split(routeAndHandler, ",")[1]
	routes := []string{`        mux.HandleFunc("` + route + `", ` + handler + `)`}
	readFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var lines []string

	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}

	readFile.Close()

	for _, l := range lines {
		fmt.Println(l)
		if strings.Contains(l, "mux.HandleFunc(") {
			routes = append(routes, l)
		}
	}

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}

	final := `package main

import "net/http"

// registerRoutes registers the routes with the provided *http.ServeMux
func registerRoutes(mux *http.ServeMux) {`

	for k, d := range routes {
		final = final + "\n" + d
		if k == len(routes)-1 {
			final = final + "\n}"
		}
	}
	log.Println(final)

	_, err = f.WriteString(final)
	if err != nil {
		log.Println(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	newHandler(handler)
}

func newHandler(handlerName string) {
	f, err := os.OpenFile(handlerName+".go", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	text := `package main

import "net/http"

func ` + handlerName + `(w http.ResponseWriter, r *http.Request) {

}`

	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}
