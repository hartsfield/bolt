package main

import (
	"fmt"
	"strings"
)

// name:string,id:string,etc:int
func genStruct(structure string) {
	spl := strings.Split(structure, ",")
	var lines []string
	lines = append(lines, "type boltGenerated struct {")
	for _, property := range spl {
		propAndType := strings.Split(property, ":")
		var ptype string
		if len(propAndType) < 2 {
			ptype = "string"
		} else {
			ptype = propAndType[1]
		}
		line := "    " + (propAndType[0]) + " " + ptype + " `json:\"" + strings.ToLower(propAndType[0]) + "\"" + " redis:\"" + strings.ToLower(propAndType[0]) + "\"`"
		lines = append(lines, line)
	}
	lines = append(lines, "}")
	for _, l := range lines {
		fmt.Println(l)
	}
}

func buildRange(structure string) {
	spl := strings.Split(structure, ":")
	var structName = spl[0]
	var props = strings.Split(spl[1], ",")
	var lines []string
	lines = append(lines, "{{ range ."+structName+" }}")
	for _, property := range props {
		line := "    <div class=\"" + property + "\">{{ ." + property + " }}</div>"
		lines = append(lines, line)
	}
	lines = append(lines, "{{end}}")
	for _, l := range lines {
		fmt.Println(l)
	}
}

func buildForm(structure string) {
	spl := strings.Split(structure, ":")
	var formName = spl[0]
	var props = strings.Split(spl[1], ",")
	var lines []string
	lines = append(lines, "<form class=\""+formName+"\">")
	for _, property := range props {
		line := "    <input class=\"" + property + "\" placeholder=\"" + property + "\"/>"
		lines = append(lines, line)
	}
	lines = append(lines, "</form>")
	for _, l := range lines {
		fmt.Println(l)
	}
}
