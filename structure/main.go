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
		setupLogging()
	}

	ctx, srv := bolt()

	fmt.Println("Waiting for connections @ http://localhost" + srv.Addr)
	log.Println("Waiting for connections @ http://localhost" + srv.Addr)

	<-ctx.Done()
}
