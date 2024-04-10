package handlers

import (
	"GoDiscordAuth/cookies"
	"GoDiscordAuth/discordApi"
	"GoDiscordAuth/templates"
	"fmt"
	"net/http"
	"time"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	var userData discordapi.DiscordIdentity

	maybeAccessTokenCookie := cookies.TryGetValidCookie(r, "AccessToken")
	_, refreshTokenNotFoundErr := r.Cookie("RefreshToken")
	if maybeAccessTokenCookie == nil {
		if refreshTokenNotFoundErr == nil {
			http.Redirect(w, r, "/api/login/refresh", http.StatusSeeOther)
			return
		}
	} else {
		userData, _ = discordapi.GetOwnDiscordIdentity(maybeAccessTokenCookie.Value)
	}

	if err := templates.Get().ExecuteTemplate(w, "index.html", userData); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleLoginFail(w http.ResponseWriter, r *http.Request) {
	if err := templates.Get().ExecuteTemplate(w, "login_fail.html", nil); err != nil {
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

	payload, err := discordapi.RequestToken(code)
	if err != nil {
		fmt.Fprintf(w, "Token were not given. Reason: %s", err.Error())
		return
	}

	accessTokenCookie := http.Cookie{Name: "AccessToken", Value: payload.AccessToken, HttpOnly: true, Path: "/"}
	refreshTokenCookie := http.Cookie{Name: "RefreshToken", Value: payload.RefreshToken, HttpOnly: true, Path: "/"}

	cookies.SetWithExpiration(w, accessTokenCookie, time.Second * time.Duration(payload.ExpiresIn))
	cookies.SetWithExpiration(w, refreshTokenCookie, time.Hour * 24 * 30) // 1 month

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func HandleAPIAccessTokenRefresh(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie("RefreshToken")
	if err != nil {
		fmt.Fprintln(w, "TODO: Handle no refresh token")
		return;
	}

	payload, _ := discordapi.RefreshAccessToken(refreshTokenCookie.Value)

	accessTokenCookie := http.Cookie{Name: "AccessToken", Value: payload.AccessToken, Path: "/", HttpOnly: true}
	cookies.SetWithExpiration(w, accessTokenCookie, time.Second * time.Duration(payload.ExpiresIn))
	// Refresh refresh token expiration
	cookies.SetWithExpiration(w, *refreshTokenCookie, time.Hour * 24 * 30) // 1 month

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
	cookies.Clear(w, accessTokenCookie)
	cookies.Clear(w, refreshTokenCookie)

	// Fixes this handler getting cached so logout doesn't work after first success
	w.Header().Add("Expires", "0")

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

