package main // viewData represents the root model used to dynamically update the app

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type viewData struct {
	CSS_Shared template.URL
	JS_Shared  template.URL
	PageTitle  string
}

//go:embed internal/shared/css/main.css
var shared_css string

//go:embed internal/shared/js/main.js
var shared_js string

// ckey/ctxkey is used as the key for the HTML context and is how we retrieve
// token information and pass it around to handlers
type ckey int

const ctxkey ckey = iota

var (
	servicePort = ":" + os.Getenv("servicePort")
	logFilePath = os.Getenv("logFilePath")

	redisIP = os.Getenv("redisIP")
	rdb     = redis.NewClient(&redis.Options{
		Addr:     redisIP + ":6379",
		Password: "",
		DB:       0,
	})
	rdx = context.Background()

	templates *template.Template = template.New("main")
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	if len(logFilePath) > 1 {
		logFile := setupLogging()
		defer logFile.Close()
	}

	ctx, srv := bolt()

	fmt.Println("Server started @ http://localhost" + srv.Addr)
	log.Println("Server started @ " + srv.Addr)

	<-ctx.Done()
}
