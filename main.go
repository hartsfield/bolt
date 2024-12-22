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
	"os/exec"
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
	if !empty || err != nil || os.MkdirAll(appName, 0755) == nil {
		log.Fatalln("Directory not empty, exiting...", err)
	}

	// if _, err := os.Stat("./" + appName + "/"); os.IsNotExist(err) {
	// 	log.Fatalln("\n    > Directory", appName, "already exists, exiting.")
	// }

	// if os.MkdirAll(appName, 0755) == nil {
	// 	log.Fatalln("Directory Exists. Exiting.")
	// }

	copyFiles("./" + appName + "/")

	cmd := exec.Command("tree", "-C", "--dirsfirst", ".")
	cmd.Dir = "./" + appName
	b, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(localCommand([]string{"go", "mod", "init", "example.com/m/v2"}))
	fmt.Print("\n    > ", strings.ReplaceAll(string(b), "\n", "\n         "))
	fmt.Println("\n    ##############################\n    # > Initialization complete. #\n    ##############################")
	fmt.Println()
}

func copyFiles(appdir string) {
	os.MkdirAll(appdir+"internal/components", 0755)
	os.MkdirAll(appdir+"internal/pages/main", 0755)
	os.MkdirAll(appdir+"internal/shared/head", 0755)
	os.MkdirAll(appdir+"public/media", 0755)

	for fileString, path := range files {
		f, err := os.Create(appdir + path)
		if err != nil {
			log.Println(err)
		}
		f.WriteString(fileString)
	}

	err := os.Chmod(appdir+"autoload.sh", 0755)
	if err != nil {
		log.Println(err)
	}
}
