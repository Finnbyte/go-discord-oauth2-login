package handlers

import (
	"GoDiscordAuth/config"
	"GoDiscordAuth/cookie"
	"GoDiscordAuth/discordApi"
	"GoDiscordAuth/templates"
	"fmt"
	"log"
	"net/http"
	"time"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	var userData discordapi.DiscordIdentity

	if accessToken, err := r.Cookie("AccessToken"); err == nil {
		// TODO Handle outdated access token
		newUserData, err := discordapi.GetOwnDiscordIdentity(accessToken.Value)
		if err == nil {
			userData = newUserData
		}
	}

	if err := templates.RenderByFilename(w, "index.html", userData); err != nil {
		log.Fatalln(err)
	}
}

func HandleLoginFail(w http.ResponseWriter, r *http.Request) {
	if err := templates.RenderByFilename(w, "login_fail.html", nil); err != nil {
		log.Fatalln(err)
	}
}

	queryParams := r.URL.Query()
	if oauthCode := queryParams.Get("code"); oauthCode != "" {
		accessToken, refreshToken, err := discordapi.RequestToken(oauthCode)
		if err != nil {
			// TODO Handle Discord refusing authorization
			panic(err)
		}

		refreshTokenCookie := http.Cookie{Name: "RefreshToken", Value: refreshToken, HttpOnly: true, Path: "/"}
		http.SetCookie(w, &refreshTokenCookie)

		accessTokenCookie := http.Cookie{Name: "AccessToken", Value: accessToken, HttpOnly: true, Path: "/"}
		http.SetCookie(w, &accessTokenCookie)
	}

	http.Redirect(w, r, "http://127.0.0.1"+config.GetOption("port"), http.StatusFound)
}
