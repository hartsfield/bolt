package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func readConf() *config {
	b, err := os.ReadFile("./bolt.conf.json")
	if err != nil {
		log.Println(err)
	}
	c := config{}
	json.Unmarshal(b, &c)
	fmt.Println(c)
	return &c
}
