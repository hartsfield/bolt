package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"strings"
	"text/template"
)

var (
	inputs   map[string][]string = make(map[string][]string)
	elements map[string]string   = make(map[string]string)
)

type data struct {
	Inputs          map[string][]string
	Lowered         map[string][]string
	Items           map[string]string
	StreamDirective string
}

// var jsn string = `{ "file": ["FileElement"], "text": ["Title","Year","Price"], "textarea": ["About"] }`

// The model is just a json file
func genStream(model_ []string) {
	model, err := os.ReadFile(model_[0])
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal([]byte(model), &inputs)
	if err != nil {
		log.Println(err)
	}
	var componentName string = "upload"
	for elm, names := range inputs {
		for _, name := range names {
			switch elm {
			case "text":
				elements[name] = "input"
			case "file":
				elements[name] = "file"
			case "textarea":
				elements[name] = "textarea"
			}
		}
	}
	b_html := makeCode(data{Items: elements, StreamDirective: buildDataStream(inputs)}, globals_streamable_html)
	b_go := makeCode(data{Inputs: inputs, Lowered: lowerMap(inputs)}, globals_streamable_go)
	writeTmpl(b_html, componentName)
	writeGo(b_go, componentName+"Handler.go")
	insertViewDirective([]string{"Stream", "[]*item"})
}

func insertViewDirective(vd []string) {
	open := "type viewData struct"
	closer := "}"
	insert := "\t" + vd[0] + " " + vd[1]
	insertLineAfter("viewdata.go", open, insert, closer)

	route := "/uploadItem"
	handler := "uploadHandler"
	routeLine := "\t" + `mux.HandleFunc("` + route + `", ` + handler + `)`
	insertLineAfter("router.go", "func registerRoutes(mux *http.ServeMux)", routeLine, closer)

	open = "view = &viewData"
	closer = "}"
	insert = "\t\t\tStream: stream,"
	insertLineAfter("helpers.go", open, insert, closer)

}

func lowerMap(ins map[string][]string) (lowered map[string][]string) {
	lowered = make(map[string][]string)
	for typ, names := range ins {
		for _, name := range names {
			lowered[typ] = append(lowered[typ], strings.ToLower(name))
		}
	}
	return
}

func makeCode(data data, glob string) []byte {
	b := bytes.NewBuffer([]byte{})
	t := template.Must(template.New("").Parse(glob))
	err := t.ExecuteTemplate(b, "", &data)
	if err != nil {
		log.Println(err)
	}
	return b.Bytes()
}

// func makeHTML() []byte {
// 	data := data{Items: elements, StreamDirective: buildDataStream(inputs)}
// 	b := bytes.NewBuffer([]byte{})
// 	t := template.Must(template.New("").Parse(globals_streamable_html))
// 	err := t.ExecuteTemplate(b, "", &data)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return b.Bytes()
// }

func buildDataStream(ins map[string][]string) string {
	build := []string{
		"<div class='stream'>",
		"\t{{ range $k, $v :=  .Stream }}",
		"\t\t<div class='item-outer'>",
	}
	var mediaElm string
	var next_lines []string = []string{"<div class='next-lines'>"}
	for _, in := range ins {
		for _, name := range in {
			if name == "FileElement" {
				mediaElm = "\t\t\t<div class='" +
					"item-part media-item " + name +
					"'><img src='{{$v.TempFileName}}'/>{{ $v." +
					name + " }}</div>"
			} else {
				next_lines = append(next_lines, "\t\t\t<div class='item-part "+name+"'>{{ $v."+name+" }}</div>")
			}
		}
	}
	next_lines = append(next_lines, "</div>")
	end := []string{
		"\t\t</div>",
		"\t{{ end }}",
		"</div>",
	}
	build = append(build, mediaElm)
	build = append(build, next_lines...)
	build = append(build, end...)
	return strings.Join(build, "\n")
}

func writeTmpl(btxt []byte, filename string) {
	err := os.MkdirAll("internal/components/"+filename+"/", os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	err = os.WriteFile("internal/components/"+filename+"/"+filename+".html", btxt, 0644)
	if err != nil {
		log.Fatal(err)
	}
	insertcomponent([]string{filename, "main"})
}

func writeGo(btxt []byte, filename string) {
	err := os.WriteFile(filename, btxt, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
