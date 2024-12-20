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

type env map[string]string

type config struct {
	App    app    `json:"app"`
	GCloud gcloud `json:"gcloud"`
}

type app struct {
	Name       string `json:"name"`
	DomainName string `json:"domain_name"`
	Version    string `json:"version"`
	Env        env    `json:"env"`
	Port       string `json:"port"`
	AlertsOn   string `json:"alertsOn"`
	TLSEnabled bool   `json:"tls_enabled"`
	Repo       string `json:"repo"`
}

type gcloud struct {
	Command   string `json:"command"`
	Zone      string `json:"zone"`
	Project   string `json:"project"`
	User      string `json:"user"`
	LiveDir   string `json:"livedir"`
	ProxyConf string `json:"proxyConf"`
}

type stringFlag struct {
	set   bool
	value string
	do    func([]string)
}

// type funcany[T any] func(doAny)

// type doAny struct {
// 	fileName       string
// 	FileData       string
// 	ComponentName  string
// 	PageName       string
// 	Sections       []string
// 	AutoSplashName string
// 	ProxConfig     string
// 	RouteName      string
// 	HandlerName    string
// 	AppName        string
// }

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
	//go:embed structure/internal/shared/head/head.css
	css_shared string
	//go:embed structure/internal/shared/head/head.js
	js_shared string
	//go:embed structure/internal/pages/main/main.html
	page_tmpl string
	//go:embed structure/internal/shared/head/head.html
	head_tmpl string
	////go:embed structure/Dockerfile
	//docker_tmpl string
	//go:embed structure/handlers.go
	handlers_go_tmpl string
	//go:embed structure/helpers.go
	helpers_go_tmpl string
	//go:embed structure/logging.go
	logging_go_tmpl string
	//go:embed structure/main.go
	main_go_tmpl string
	//go:embed structure/router.go
	router_tmpl string
	//go:embed structure/server.go
	server_go_tmpl string
	//go:embed structure/bolt.conf.json
	bolt_conf_tmpl string
	//go:embed structure/viewdata.go
	viewdata_go_tmpl string
	//go:embed structure/globals.go
	globals_go_tmpl string
	//go:embed structure/autoload.sh
	globals_autoload_sh string
	//go:embed streamable/tmpls/submitform.html
	globals_streamable_html string
	//go:embed streamable/tmpls/submithandler.go
	globals_streamable_go string
)

// other globals
var (
	// f(lag)Map maps flag strings to a *stringFlag{}
	fMap map[string]*stringFlag = make(map[string]*stringFlag)

	components_dir string = "internal/components/"
	rc             *config
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fMap["init"] = &stringFlag{do: boltInit}
	flag.Var(fMap["init"], "init", "\nEx.: --init appName\n\nInitializes a new bolt project in the directory 'appName'\n")

	fMap["new-page"] = &stringFlag{do: createPage}
	flag.Var(fMap["new-page"], "new-page", "Initializes a new page with the given name")

	fMap["new-component"] = &stringFlag{do: createComponent}
	flag.Var(fMap["new-component"], "new-component", "Initializes a new component with the given name")

	fMap["new-route"] = &stringFlag{do: newRoute}
	flag.Var(fMap["new-route"], "new-route", "Initializes a new route")

	fMap["insert-component"] = &stringFlag{do: insertcomponent}
	flag.Var(fMap["insert-component"], "insert-component", "Inserts a new component into a page")

	fMap["streamable"] = &stringFlag{do: genStream}
	flag.Var(fMap["streamable"], "streamable", "unfinished")

	fMap["deploy"] = &stringFlag{do: deploy}
	flag.Var(fMap["deploy"], "deploy", "Deploys project to server")

	fMap["autonav"] = &stringFlag{do: autonav}
	flag.Var(fMap["autonav"], "autonav", "Initializes a new navbar component")

	fMap["autosplash"] = &stringFlag{do: autosplash}
	flag.Var(fMap["autosplash"], "autosplash", "Initializes a splash screen component")

	fMap["remote-service-restart"] = &stringFlag{do: remoteServiceRestart}
	flag.Var(fMap["remote-service-restart"], "remote-service-restart", "Restarts a remote service")

	// flag.Var(fMap["add-style"], "add-style", "Adds a style to the stylesheet of the given component, usage: bolt --add-style=component:rulename")
	// flag.Var(fMap["build-form"], "build-form", "Genrates an HTML form based on input")
	// flag.Var(fMap["install-component"], "install-component", "Installs a component from a git hub repo")
}

func readFlags() {
	flag.Parse()
	noFlagsSet := true
	for _, clf := range fMap {
		if clf.set {
			noFlagsSet = false
			clf.do(strings.Split(clf.value, ","))
		}
	}
	if noFlagsSet {
		boltInit([]string{""})
	}
}

func boltInit(params []string) {
	appName := params[0]
	empty, err := isEmpty(".")
	if err != nil {
		log.Fatalln(err)
	}
	if !empty && len(appName) < 1 {
		_, err = os.ReadFile("bolt.conf.json")
		if err == nil {
			fmt.Println("\n   > bolt.conf.json detected, this directory already appears to contain a bolt project, exiting.")
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

	os.MkdirAll(appdir+"internal/components", 0755)
	os.MkdirAll(appdir+"internal/pages/main", 0755)
	os.MkdirAll(appdir+"internal/shared/head", 0755)
	os.MkdirAll(appdir+"public/media", 0755)

	main_go, err := os.Create(appdir + "main.go")
	if err != nil {
		log.Println(err)
	}
	main_go.WriteString(main_go_tmpl)

	bctmpl, err := os.Create(appdir + "bolt.conf.json")
	if err != nil {
		log.Println(err)
	}
	bctmpl.WriteString(bolt_conf_tmpl)

	htmpl, err := os.Create(appdir + "internal/shared/head/head.html")
	if err != nil {
		log.Println(err)
	}
	htmpl.WriteString(head_tmpl)

	pmain, err := os.Create(appdir + "internal/pages/main/main.html")
	if err != nil {
		log.Println(err)
	}
	pmain.WriteString(page_tmpl)

	shcss, err := os.Create(appdir + "internal/shared/head/head.css")
	if err != nil {
		log.Println(err)
	}
	shcss.WriteString(css_shared)

	shjs, err := os.Create(appdir + "internal/shared/head/head.js")
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

	globals_go, err := os.Create(appdir + "globals.go")
	if err != nil {
		log.Println(err)
	}
	globals_go.WriteString(globals_go_tmpl)

	viewdata_go, err := os.Create(appdir + "viewdata.go")
	if err != nil {
		log.Println(err)
	}
	viewdata_go.WriteString(viewdata_go_tmpl)

	autoloadSh, err := os.Create(appdir + "autoload.sh")
	if err != nil {
		log.Println(err)
	}
	autoloadSh.WriteString(globals_autoload_sh)
	err = os.Chmod(autoloadSh.Name(), 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	// dockerfile, err := os.Create(appdir + "Dockerfile")
	// if err != nil {
	// 	log.Println(err)
	// }
	// dockerfile.WriteString(docker_tmpl)

	router, err := os.Create(appdir + "router.go")
	if err != nil {
		log.Println(err)
	}
	router.WriteString(router_tmpl)

	cmd := exec.Command("tree", "-C", "--dirsfirst", ".")
	b, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}
	log.Println("\n", localCommand([]string{"go", "mod", "init", "example.com/m/v2"}))
	// log.Println(localCommand([]string{"go", "mod", "tidy"}))
	fmt.Print("\n    > ", strings.ReplaceAll(string(b), "\n", "\n         "))

	fmt.Println("\n    ##############################\n    # > Initialization complete. #\n    ##############################")
	fmt.Println()
}
