package main

func insertcomponent(params []string) {
	component := params[0]
	var page string
	if len(params) > 1 {
		page = params[1]
	} else {
		page = "main"
	}
	filePath := "internal/pages/" + page + "/" + page + ".html"
	opener := "<body"
	insert := "\t{{template \"" + component + ".html\" . }}"
	closer := "</body"
	insertLineAfter(filePath, opener, insert, closer)
}
