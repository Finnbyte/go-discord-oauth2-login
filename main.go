package main

import (
	"GoDiscordAuth/config"
	"GoDiscordAuth/handlers"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Could not load .env file")
	}

	config.SetOption("clientId", os.Getenv("CLIENT_ID"))
	config.SetOption("clientSecret", os.Getenv("CLIENT_SECRET"))
	config.SetOption("port", os.Getenv("PORT"))

    http.HandleFunc("/", handlers.HandleIndex)
	http.HandleFunc("/api/login", handlers.HandleAPILogin)
	http.HandleFunc("/api/login/callback", handlers.HandleAPILoginCallback)
	http.HandleFunc("/api/logout", handlers.HandleAPILogout)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.GetOption("port")), nil); err != nil {
        log.Fatalln(err)
    }
}
