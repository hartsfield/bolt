package main

import (
	"log"
	"os"
	"strings"
)

func insertcomponent(component string) {
	f, err := os.OpenFile("internal/pages/main/main.tmpl", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	var b []byte
	_, err = f.Read(b)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(b))

	lines := strings.Split(string(b), "\n")
	for k, l := range lines {
		log.Println(l)
		if strings.Contains(l, "footer") {
			var temp []string = append(lines[:k], `{{template "`+component+`"}}`)
			lines = append(temp, lines[k:]...)
		}
	}

	for _, l := range lines {
		log.Println(l)
		if _, err = f.WriteString(l); err != nil {
			panic(err)
		}
	}
}
