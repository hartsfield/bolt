package main

import "strings"

// name:string,id:string,etc:int
func genStruct(structure string) {
	spl := strings.Split(structure, ",")
	var lines []string
	for _, property := range spl {
		propAndType := strings.Split(property, ":")
		var ptype string
		if len(propAndType) < 2 {
			ptype = "string"
		} else {
			ptype = propAndType[1]
		}
		line := strings.Title(propAndType[0]) + " " + ptype + "`json\"" + strings.ToLower(propAndType[0]) + "\"`"
		lines = append(lines, line)
	}
}
