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

func HandleAPILogin(w http.ResponseWriter, r *http.Request) {
	const discordLoginUrl = "https://discord.com/oauth2/authorize?client_id=1215381408772788234&response_type=code&redirect_uri=http%3A%2F%2F127.0.0.1%3A8080%2Fapi%2Flogin%2Fcallback&scope=identify"
	http.Redirect(w, r, discordLoginUrl, http.StatusSeeOther)
}

func HandleAPILoginCallback(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	code := queryParams.Get("code")

	if code == "" {
		fmt.Fprint(w, "Code was empty")
		return
	}

	accessToken, refreshToken, err := discordapi.RequestToken(code)
	if err != nil {
		fmt.Fprint(w, "Token were not given. Reason: %s", err.Error())
		return
	}

	accessTokenCookie := http.Cookie{Name: "AccessToken", Value: accessToken, HttpOnly: true, Path: "/"}
	refreshTokenCookie := http.Cookie{Name: "RefreshToken", Value: refreshToken, HttpOnly: true, Path: "/"}

	cookie.SetWithExpiration(w, accessTokenCookie, time.Second * 5)
	cookie.SetWithExpiration(w, refreshTokenCookie, time.Second * 5)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

