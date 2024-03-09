package handlers

import (
	"GoDiscordAuth/cookie"
	"GoDiscordAuth/discordApi"
	"GoDiscordAuth/templates"
	"fmt"
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
        http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleLoginFail(w http.ResponseWriter, r *http.Request) {
	if err := templates.RenderByFilename(w, "login_fail.html", nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
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

	cookie.SetWithExpiration(w, accessTokenCookie, time.Minute * 30)
	cookie.SetWithExpiration(w, refreshTokenCookie, time.Minute * 30)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func HandleAPIRefreshTokenRefresh(w http.ResponseWriter, r *http.Request) {
	fmt.Println("refresh handler")
	refreshTokenCookie, err := r.Cookie("RefreshToken")
	if err != nil {
		fmt.Fprintln(w, "TODO: Handle no refresh token")
	}

	payload, _ := discordapi.RefreshAccessToken(refreshTokenCookie.Value)

	accessTokenCookie := http.Cookie{Name: "AccessToken", Value: payload.AccessToken, Path: "/", HttpOnly: true}
	cookie.SetWithExpiration(w, accessTokenCookie, time.Second * time.Duration(payload.ExpiresIn))
	// Refresh refresh token expiration
	cookie.SetWithExpiration(w, *refreshTokenCookie, time.Hour * 24 * 30) // 1 month

	whereCameFromURL := r.URL.Path
	http.Redirect(w, r, whereCameFromURL, http.StatusSeeOther)
}

func HandleAPILogout(w http.ResponseWriter, r *http.Request) {
	accessTokenCookie, err := r.Cookie("AccessToken")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return;
	}

	err = discordapi.RevokeTokens(accessTokenCookie.Value)
	if err != nil {
		panic(err)
	}

	// Make both cookies invalid so new ones will be retrieved when needed
	refreshTokenCookie, _ := r.Cookie("RefreshToken")
	cookie.Clear(w, accessTokenCookie)
	cookie.Clear(w, refreshTokenCookie)
	
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

