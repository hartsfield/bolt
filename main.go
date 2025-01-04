// NOTE: Need:
// bolt bootstrap
// bolt config
// bolt interactive
package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// NOTE: see init.go

func main() {
	readFlags()
}

func readFlags() {
	flag.Parse()
	noFlagsSet := true
	for _, com := range flagCommands {
		if com.set {
			noFlagsSet = false
			com.do(strings.Split(com.value, ","))
		}
	}
	if noFlagsSet {
		boltInit([]string{""})
	}
}

func boltInit(params []string) {
	appName := params[0]

	empty, err := isEmpty(".")
	fmt.Println(empty, err)
	if empty && err != nil && os.MkdirAll(appName, 0755) == nil {
		fmt.Println(empty, err)
		log.Fatalln("Directory not empty, exiting...", err)
	}

	os.Chdir(appName)
	copyFiles(appName)
	b := localCommand("tree -C --dirsfirst .")
	fmt.Println(localCommand("go mod init example.com/m/v2"))
	fmt.Print("\n    > ", strings.ReplaceAll(string(b), "\n", "\n         "))
	fmt.Println("\n    ##############################\n    # > Initialization complete. #\n    ##############################")
	fmt.Println()
	if len(params) >= 2 {
		mkGitHubRepo := params[1]
		if mkGitHubRepo == "git" {
			localCommand("git init")
			localCommand("git add .")
			localCommand("git commit -m init_commit")
		}
		if len(params) >= 3 {
			privateOrPublic := params[2]
			localCommand("gh repo create --source=. --" + privateOrPublic)
			localCommand("git push origin master")
		}
	}
	fmt.Println()
}

func copyFiles(appdir string) {
	os.MkdirAll("internal/components", 0755)
	os.MkdirAll("internal/pages/main", 0755)
	os.MkdirAll("internal/shared/head", 0755)
	os.MkdirAll("public/media", 0755)

	for fileString, path := range files {
		f, err := os.Create(path)
		if err != nil {
			log.Println(err)
		}
		f.WriteString(fileString)
	}

	writeConf(defaultConf([]string{
		appdir,
		"9343",
		"us-central1-a",
		"mysterygift",
		"main",
		"john",
	}))

	f, err := os.Create("autoload.sh")
	if err != nil {
		log.Println(err)
	}
	_, err = f.Write([]byte("# use in n/vim to restart on save:\n" +
		"# :autocmd BufWritePost * silent! !./autoload.sh\n" +
		"#!/bin/bash\n" +
		"pkill " + appdir + " || true\n" +
		"go build -o " + appdir + "\n" +
		"./" + appdir + " >> log.txt 2>&1 &"))
	if err != nil {
		log.Println(err)
	}

	err = os.Chmod("autoload.sh", 0755)
	if err != nil {
		log.Println(err)
	}
}
