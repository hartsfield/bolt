package main // viewData represents the root model used to dynamically update the app

import (
	"fmt"
	"log"
	"math/rand"
	"time"
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
