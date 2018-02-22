package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

var (
	HTML_FILE   = flag.String("html", "", "the html bot file")
	HTTP_SERVER = flag.String("http", ":8976", "the http listen address")
)

func main() {
	flag.Parse()
	file, err := os.Open(*HTML_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	bot, err := NewBotFromReader(file)
	if err != nil {
		log.Fatal(err)
	}
	m := NewMessenger(bot)
	http.ListenAndServe(*HTTP_SERVER, m)
}
