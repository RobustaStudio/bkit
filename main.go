package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	file, err := os.Open("bot.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	bot, err := NewBotFromReader(file)
	if err != nil {
		log.Fatal(err)
	}

	m := NewMessenger(bot)

	http.ListenAndServe(":8000", m)
}
