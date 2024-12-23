package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func readConf() *config {
	b, err := os.ReadFile("./bolt.conf.json")
	if err != nil {
		log.Println(err)
	}
	c := config{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		log.Println(err)
	}
	return &c
}

func createPage(params []string) {
	for _, p := range params {
		create(p, "pages")
	}
}

func createComponent(params []string) {
	for _, p := range params {
		create(p, "components")
	}
}

func create(name, structure string) {
	wd := "internal/" + structure + "/"
	ex, err := exists(wd + name)
	if err != nil || ex {
		fmt.Println("Component already exists", err)
		os.Exit(0)
	}
	os.MkdirAll(wd+name, 0755)
	tmpl_, err := os.Create(wd + name + "/" + name + ".html")
	if err != nil {
		log.Println(err)
	}
	tmpl_.WriteString(`<div class="template-wrapper ` + name +
		`-outer" id="` + name + `-outer">` + "\n" +
		`<style>{{ template "` + name + `.css" . }}</style>` + "\n" +
		`<script>{{ template "` + name + `.js" . }}</script>` + "\n" +
		`</div>`)

	_, err = os.Create(wd + name + "/" + name + ".css")
	if err != nil {
		log.Println(err)
	}

	_, err = os.Create(wd + name + "/" + name + ".js")
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Created:", "\n", wd+name+"/"+name+"{.html,.css,.js}")
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

func remoteServiceRestart(args []string) {
	name := args[0]
	log.Println("cd " + name + " && go build -o " + name + " && pkill -f " + name + " && servicePort=$(cat ~/prox.conf | grep $2 | cut -d: -f1) logFilePath=./logfile.txt ./" + name + " &")
	log.Println(cloudCommand([]string{"cd " + name + " && go build -o " + name + " && pkill -f " + name + " && servicePort=$(cat ~/prox.conf | grep $2 | cut -d: -f1) logFilePath=./logfile.txt ./" + name + " &"}))
}

func localCommand(command []string) string {
	var cmd *exec.Cmd = &exec.Cmd{}
	cmd.Env = append(cmd.Env, "GOARCH=amd64")
	cmd.Env = append(cmd.Env, "GOOS=linux")
	cmd = exec.Command(command[0], command[1:]...)
	o, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("local command error: ", err, string(o))
	}
	return string(o)
}

func cloudCommand(command []string) string {
	args := []string{`compute`, `ssh`, `--zone`, `us-central1-a`, `main`, `--project`, `mysterygift`, `--`}

	args = append(args, command...)
	cmd := exec.Command(`gcloud`, args...)
	o, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	return string(o)
}

func insertLineAfter(filepath, opening, insert, closing string) {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(f)

	var lines []string
	var curlyCount int = 0
	var found bool = false
	for fileScanner.Scan() {
		t := fileScanner.Text()
		if strings.Contains(t, opening) {
			found = true
			curlyCount = curlyCount + 1
		}
		if found {
			if strings.Contains(t, closing) {
				curlyCount = curlyCount - 1
			}
			if curlyCount == 0 {
				lines = append(lines, insert)
				found = false
			}
		}
		lines = append(lines, t)
	}
	err = f.Truncate(0)
	if err != nil {
		log.Println(err)
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		log.Println(err)
	}

	_, err = f.WriteString(strings.Join(lines, "\n"))
	if err != nil {
		log.Println(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
