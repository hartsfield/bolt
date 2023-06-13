package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type stringFlag struct {
	set   bool
	value string
	do    func(string)
}

func (sf *stringFlag) Set(x string) error {
	sf.value = x
	sf.set = true
	return nil
}
func (sf *stringFlag) String() string {
	return sf.value
}

// embed the file templates
var (
	//go:embed structure/internal/shared/css/main.css
	css_shared string
	//go:embed structure/internal/shared/js/main.js
	js_shared string
	//go:embed structure/internal/pages/main/main.tmpl
	page_tmpl string
	//go:embed structure/internal/components/head/head.tmpl
	head_tmpl string
	//go:embed structure/internal/components/footer/footer.tmpl
	foot_tmpl string
	//go:embed structure/internal/components/footer/footer.css
	foot_css string
	//go:embed structure/internal/components/footer/footer.js
	foot_js string
	//go:embed structure/Dockerfile
	docker_tmpl string
	//go:embed structure/handlers.go
	handlers_go_tmpl string
	//go:embed structure/helpers.go
	helpers_go_tmpl string
	//go:embed structure/logging.go
	logging_go_tmpl string
	//go:embed structure/main.go
	main_go_tmpl string
	//go:embed structure/restart-service.sh
	rssh_tmpl string
	//go:embed structure/server.go
	server_go_tmpl string
	//go:embed structure/bolt_conf.json
	bolt_conf_tmpl string
)

// other globals
var (
	// flagMap maps flag strings to a *stringFlag{}
	fMap map[string]*stringFlag = make(map[string]*stringFlag)
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	rand.Seed(time.Now().UTC().UnixNano())

	fMap["init"] = &stringFlag{do: boltInit}
	fMap["add-component"] = &stringFlag{}
	fMap["new-component"] = &stringFlag{do: createComponent}
	fMap["add-page"] = &stringFlag{}
	fMap["new-page"] = &stringFlag{do: createPage}

	flag.Var(fMap["init"], "init", "Initializes a new bolt project")
	flag.Var(fMap["add-component"], "add-component", "Installs a component from a remote git repository")
	flag.Var(fMap["new-component"], "new-component", "Initializes a new component with the given name")
	flag.Var(fMap["add-page"], "add-page", "Installs a page template from a remote git repository")
	flag.Var(fMap["new-page"], "new-page", "Initializes a new page with the given name")
}

func main() {
	readFlags()
}

func readFlags() {
	flag.Parse()
	noFlagsSet := true
	for _, clf := range fMap {
		if clf.set {
			noFlagsSet = false
			clf.do(clf.value)
		}
	}
	if noFlagsSet {
		boltInit("")
	}
}
func createPage(name string) {
	create(name, "pages")
}
func createComponent(name string) {
	create(name, "components")
}
func create(name, structure string) {
	wd := "internal/" + structure + "/"
	os.MkdirAll(wd+name, 0755)
	tmpl_, err := os.Create(wd + name + "/" + name + ".tmpl")
	if err != nil {
		log.Println(err)
	}
	tmpl_.WriteString(`{{ define "` + name + `.tmpl" . }}` + "\n" + `{{end}}`)

	css_, err := os.Create(wd + name + "/" + name + ".css")
	if err != nil {
		log.Println(err)
	}
	css_.WriteString(`{{ define "` + name + `.css" . }}` + "\n" + `{{end}}`)

	js_, err := os.Create(wd + name + "/" + name + ".js")
	if err != nil {
		log.Println(err)
	}
	js_.WriteString(`{{ define "` + name + `.js" . }}` + "\n" + `{{end}}`)
	fmt.Println("Created:", "\n", wd+name+"/"+name+"{.tmpl,.css,.js}")
}

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

func boltInit(appName string) {
	empty, err := isEmpty(".")
	if err != nil {
		log.Fatalln(err)
	}
	if !empty && len(appName) < 1 {
		_, err = os.ReadFile("bolt_conf.json")
		if err == nil {
			fmt.Println("\n   > bolt_conf.json detected, this directory already appears to contain a bolt project, exiting.")
			fmt.Println()
			os.Exit(0)
		}
		log.Fatalln("Directory not empty, exiting...")
	}

	var appdir string = "./"
	if len(appName) > 0 {
		if _, err := os.Stat(appdir + appName + "/"); os.IsNotExist(err) {
			err := os.MkdirAll(appName, 0755)
			if err != nil {
				log.Println(err)
			}

			appdir = appdir + appName + "/"
			fmt.Println("Created app directory:", appName+"/")
		} else {
			log.Fatalln("\n    >  Directory", appName, "already exists, exiting.")
		}
	} else {
		fmt.Println("\n    > Initializing app in current directory...")
	}

	os.MkdirAll(appdir+"internal/components/footer", 0755)
	os.MkdirAll(appdir+"internal/components/head", 0755)
	os.MkdirAll(appdir+"internal/pages/main", 0755)
	os.MkdirAll(appdir+"internal/shared/css", 0755)
	os.MkdirAll(appdir+"internal/shared/js", 0755)
	os.MkdirAll(appdir+"public/media", 0755)

	main_go, err := os.Create(appdir + "main.go")
	if err != nil {
		log.Println(err)
	}
	main_go.WriteString(main_go_tmpl)

	bctmpl, err := os.Create(appdir + "bolt_conf.json")
	if err != nil {
		log.Println(err)
	}
	bctmpl.WriteString(bolt_conf_tmpl)

	ftmpl, err := os.Create(appdir + "internal/components/footer/footer.tmpl")
	if err != nil {
		log.Println(err)
	}
	ftmpl.WriteString(foot_tmpl)

	fcss, err := os.Create(appdir + "internal/components/footer/footer.css")
	if err != nil {
		log.Println(err)
	}
	fcss.WriteString(foot_css)

	fjs, err := os.Create(appdir + "internal/components/footer/footer.js")
	if err != nil {
		log.Println(err)
	}
	fjs.WriteString(foot_js)

	htmpl, err := os.Create(appdir + "internal/components/head/head.tmpl")
	if err != nil {
		log.Println(err)
	}
	htmpl.WriteString(head_tmpl)

	pmain, err := os.Create(appdir + "internal/pages/main/main.tmpl")
	if err != nil {
		log.Println(err)
	}
	pmain.WriteString(page_tmpl)

	shcss, err := os.Create(appdir + "internal/shared/css/main.css")
	if err != nil {
		log.Println(err)
	}
	shcss.WriteString(css_shared)

	shjs, err := os.Create(appdir + "internal/shared/js/main.js")
	if err != nil {
		log.Println(err)
	}
	shjs.WriteString(js_shared)

	logging_go, err := os.Create(appdir + "logging.go")
	if err != nil {
		log.Println(err)
	}
	logging_go.WriteString(logging_go_tmpl)

	server_go, err := os.Create(appdir + "server.go")
	if err != nil {
		log.Println(err)
	}
	server_go.WriteString(server_go_tmpl)

	handlers_go, err := os.Create(appdir + "handlers.go")
	if err != nil {
		log.Println(err)
	}
	handlers_go.WriteString(handlers_go_tmpl)

	helpers_go, err := os.Create(appdir + "helpers.go")
	if err != nil {
		log.Println(err)
	}
	helpers_go.WriteString(helpers_go_tmpl)

	dockerfile, err := os.Create(appdir + "Dockerfile")
	if err != nil {
		log.Println(err)
	}
	dockerfile.WriteString(docker_tmpl)

	rs_sh, err := os.Create(appdir + "restart-service.sh")
	if err != nil {
		log.Println(err)
	}
	rs_sh.WriteString(rssh_tmpl)

	cmd := exec.Command("tree", "--dirsfirst", ".")
	b, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}
	fmt.Print("\n    > ", strings.ReplaceAll(string(b), "\n", "\n         "))

	fmt.Println("\n    ##############################\n    # > Initialization complete. #\n    ##############################")
	fmt.Println()
}
