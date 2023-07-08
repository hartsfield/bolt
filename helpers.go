package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func createPage(name string) {
	create(name, "pages")
}
func createComponent(name string) {
	create(name, "components")
}
func create(name, structure string) {
	wd := "internal/" + structure + "/"
	ex, err := exists(wd + name)
	if err != nil || ex {
		fmt.Println("Component already exists", err)
		os.Exit(0)
	}
	os.MkdirAll(wd+name, 0755)
	tmpl_, err := os.Create(wd + name + "/" + name + ".tmpl")
	if err != nil {
		log.Println(err)
	}
	tmpl_.WriteString(`{{ define "` + name + `.tmpl" }}` + "\n" +
		`<style>{{ template "` + name + `.css" }}</style>` + "\n" +
		`<script>{{ template "` + name + `.js"}}</script>` + "\n" +
		`{{end}}`)

	css_, err := os.Create(wd + name + "/" + name + ".css")
	if err != nil {
		log.Println(err)
	}
	css_.WriteString(`{{ define "` + name + `.css" }}` + "\n" + `{{end}}`)

	js_, err := os.Create(wd + name + "/" + name + ".js")
	if err != nil {
		log.Println(err)
	}
	js_.WriteString(`{{ define "` + name + `.js" }}` + "\n" + `{{end}}`)
	fmt.Println("Created:", "\n", wd+name+"/"+name+"{.tmpl,.css,.js}")
}

func autoNav(name string, sections []string) {
	wd := "internal/components/"
	ex, err := exists(wd + name)
	if err != nil || ex {
		fmt.Println("Component already exists", err)
		os.Exit(0)
	}
	os.MkdirAll(wd+name, 0755)
	tmpl_, err := os.Create(wd + name + "/" + name + ".tmpl")
	if err != nil {
		log.Println(err)
	}
	tmpl_.WriteString(`
{{ define "` + name + `.tmpl" }}
<div class="navbar-outer">
  <div class="logo-nav" onclick="window.location = window.location.origin">{{ .CompanyName }} - Security Enhancements&nbsp;🛡️</div>

  <div class="nav-landscape">
    <ul>
       <li onclick="jumpTo('section-services')">Services & Products</li>
       <li onclick="jumpTo('section-software')">Business Software</li>
       <li onclick="jumpTo('section-location')">Locations</li>
       <li onclick="jumpTo('section-contact')">Contact</li>
    </ul>
  </div>

  <div class="nav-portrait" id="nav-portrait">
    <div class="nav-portrait-logo">{{ .CompanyName }}</div>
    <ul>
       <li onclick="jumpTo('section-services')">Services & Products</li>
       <li onclick="jumpTo('section-software')">Business Software</li>
       <li onclick="jumpTo('section-location')">Locations</li>
       <li onclick="jumpTo('section-contact')">Contact</li>
    </ul>
  </div>

  <div class="ham-outer" onclick="showNavPortrait()">
    <div class="hamburger"></div>
    <div class="hamburger"></div>
    <div class="hamburger"></div>
  </div>

</div>

<style>{{ template "` + name + `.css" }}</style>
<script>{{ template "` + name + `.js"}}</script>
{{end}}
                `)

	// tmpl_.WriteString(`{{ define "` + name + `.tmpl" }}` + "\n" +
	// 	`<style>{{ template "` + name + `.css" }}</style>` + "\n" +
	// 	`<script>{{ template "` + name + `.js"}}</script>` + "\n" +
	// 	`{{end}}`)

}
func autoList(name string, listItems []string) {}
func autoFlex(name string)                     {}

func isEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func addStyle(name string) {
	s := strings.SplitN(name, ":", 3)
	empty, err := isEmpty(components_dir + s[0])
	if err != nil {
		log.Fatalln(err)
	}

	if empty {
		log.Fatalln("Component not found:", s[0])
	}

	cssfile, err := os.OpenFile(components_dir+s[0]+"/"+s[0]+".css", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	defer cssfile.Close()
	_, err = cssfile.WriteString(s[1] + " {}")
	if err != nil {
		log.Println(err)
	}
	fmt.Println("\n    > css rule added to component:")
	fmt.Println("\n            rule:", s[1], "{}")
	fmt.Println("       component:", s[0])
	fmt.Println("        modified:", components_dir+s[0]+"/"+s[0]+".css")
	fmt.Println()
	fmt.Println()
}

// bolt --add-style component:element:position,display,z-index

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
