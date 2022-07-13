package main

import (
	"log"
	"os"
	"url_shortener/pkg/url_shortener"
)

var mappings = make(map[string]string)
var hostName string

func main() {
	if val := os.Getenv("hostname"); val == "" {
		panic("Need hostname env value")
	}
	hostName = os.Getenv("hostname")
	a := url_shortener.NewApp(hostName)
	err := a.Start()
	if err != nil {
		log.Fatal(err)
	}
}
