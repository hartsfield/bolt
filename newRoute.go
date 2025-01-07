package main

import (
	"os"
	"strings"
)

func newRoute(params []string) {
	route := params[0]
	if !strings.Contains(route, "/") {
		route = "/" + route
	}
	handler := params[1]
	routeLine := "\t" + `mux.HandleFunc("` + route + `", ` + handler + `)`
	insertLineAfter("router.go", "func registerRoutes(mux *http.ServeMux)", routeLine, "}")
	// if len(params) > 2 {
	// 	newHandler(handler, []byte(params[2]), nil)
	// 	return
	// }
	// newHandler(handler, nil, []string{"\"net/http\"", "\"strings\""})
}

func newHandler(handlerName string, fileBytes []byte, imports []string) {
	f, err := os.OpenFile(handlerName+".go", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var inner string
	if fileBytes != nil {
		inner = string(fileBytes)
	}
	var imstring string = strings.Join(imports, "\n")

	txt := "package main\n\nimport (\n\t" + imstring + "\n)\n\nfunc " + handlerName +
		"(w http.ResponseWriter, r *http.Request) {\n\t" + inner + "\n}"

	if _, err = f.WriteString(txt); err != nil {
		panic(err)
	}
}
