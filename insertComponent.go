package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func insertcomponent(component string) {
	filePath := "internal/pages/main/main.tmpl"
	templateDirectives := []string{`    {{template "` + component + `.tmpl" . }}`}
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
		if strings.Contains(l, "{{") {
			if !strings.Contains(l, ` "head" .`) && !strings.Contains(l, "end") {
				templateDirectives = append(templateDirectives, l)
			}
		}
	}

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}

	final := `<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 3.2 Final//EN">
<html>
{{template "head" . }} 
  <body>`

	for k, d := range templateDirectives {
		final = final + "\n" + d
		if k == len(templateDirectives)-1 {
			final = final + "\n  </body>\n</html>"
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
}
