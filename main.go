package main

import (
	"GoDiscordAuth/config"
	"GoDiscordAuth/handlers"
	"fmt"
	"log"
	"net/http"
	"os"

)

func main() {
    http.HandleFunc("/", handlers.HandleIndex)
	http.HandleFunc("/api/login", handlers.HandleLogin)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.GetOption("port")), nil); err != nil {
        log.Fatalln(err)
    }
}
