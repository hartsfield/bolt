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
	streamjs string              = `async function submitPost() {
        const form = document.getElementById("uploadForm");
        const data = new FormData(form);
        let response = await fetch("/uploadItem", {
            method: "POST",
            body: data,
        });

        let res = await response.json();
        handleResponse(res);
    }

    function handleResponse(res) {
        if (res.success == "true") {
            window.location = window.location.origin;
        } else {
            document.getElementById("errorField").innerHTML = res.error;
        }
    }`
	streamcss string = `body, html {
    margin: 3em 0.5em;
}
input, textarea {
    border: none;
}
::placeholder {
    color: var(--html-bg);
    opacity: 1; /* Firefox */
}
.stream {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    justify-content: flex-start;
    margin: 0.5em;
}
.uploadForm {
    display: flex;
    flex-direction: column;
    margin-bottom: 1em;
    margin: 1em;
}
.uploadForm > * {
    padding: 0.3em;
    width: 100%;
    border-radius: 0.3em;
    margin-top: 0.5em;
    background: #FFFFFF;
    border: 1px solid #e8e8e8;
}
.form-submit {
    text-align: center;
    padding-left: 0;
    padding-right: 0;
    background: #f38d1c;
    border: 1px solid orange;
    color: white;
}
.stream > div {
    margin: 0.5em;
    padding: 0.5em;
    border-radius: 0.4em;
    width: 25%;
    display: flex;
    flex-direction: column;
    flex-wrap: nowrap;
    flex-grow: 1;
    justify-content: space-between;
    align-items: stretch;
    background: #f1f1f1;
    cursor: pointer;
}
.media-item > img {
    width: 100%;
    border-radius: 0.4em;
}
.next-lines {
    margin-top: 0.8rem;
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    justify-content: space-around;
    background: #e06767;
    color: #e8e8e8;
    padding: 0.6em;
    border-radius: 0.3em;
    font-size: 0.9em;
    align-items: center;
}
.About {
    width: 100%;
}
@media screen and (orientation:landscape) {
    body, html {
        max-width: 80ch;
    }
    .stream {
        max-width: 80ch;
    }
}

`
)

type data struct {
	Inputs          map[string][]string
	Lowered         map[string][]string
	Items           map[string]string
	StreamDirective string
}

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
	b_go := makeCode(data{Inputs: inputs, Lowered: inputs}, globals_streamable_go)
	writeGo(b_go, componentName+"Handler.go")
	b_html := makeCode(data{Items: elements, StreamDirective: buildDataStream(inputs)}, globals_streamable_html)
	writeTmpl(b_html, componentName, ".html")
	b_css := makeCode(data{}, streamcss)
	writeTmpl(b_css, componentName, ".css")
	b_js := makeCode(data{}, streamjs)
	writeTmpl(b_js, componentName, ".js")
	insertcomponent([]string{componentName, "main"})

	insertViewDirective([]string{"Stream", "[]*item"})
	newRoute([]string{"/uploadItem", "uploadHandler"})
	// newHandler("uploadHandler", nil, []string{"\"net/http\"", "\"strings\""})
	newRoute([]string{"/view", "viewItem"})
	newHandler("viewItem",
		[]byte("readDB()\n        id := strings.Split(r.RequestURI, \""+
			"/\")[2]\n\tfmt.Println(r.RequestURI, id, itemsMap[id], itemsMap)"+
			"\n\texeTmpl(w, r, &viewData{AppName:"+
			" appConf.App.Name, "+
			"Stream: []*item{itemsMap[id]}}, \"main.html\")"),
		[]string{"\"net/http\"", "\t\"strings\"", "\t\"fmt\""},
	)

}

func insertViewDirective(vd []string) {
	open := "type viewData struct"
	closer := "}"
	insert := "\t" + vd[0] + " " + vd[1]
	insertLineAfter("viewdata.go", open, insert, closer)
	insert = "\tItem *item"
	insertLineAfter("viewdata.go", open, insert, closer)

	open = "view = &viewData"
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
		"\t\t<div class='item-outer' onclick=\"window.location = '/view/{{$v.ID}}'\">",
	}
	var mediaElm string
	var next_lines []string = []string{"<div class='next-lines'>"}
	for _, in := range ins {
		for _, name := range in {
			if name == "Media" {
				mediaElm = "\t\t\t<div class='" +
					"item-part media-item " + name +
					"'><img src='/{{$v.TempFileName}}'/>{{ $v." +
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

func writeTmpl(btxt []byte, filename, ext string) {
	err := os.MkdirAll("internal/components/"+filename, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	mkMultiFile(btxt, "internal/components/"+filename+"/"+filename, ext)
}

func mkMultiFile(btxt []byte, dir, ext string) {
	err := os.WriteFile(dir+ext, btxt, 0644)
	if err != nil {
		log.Println(err)
	}
}

func writeGo(btxt []byte, filename string) {
	err := os.WriteFile(filename, btxt, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
