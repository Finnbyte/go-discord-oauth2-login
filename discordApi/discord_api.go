package discordapi

import (
	"GoDiscordAuth/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	URL "net/url"
)

const discordApiUrl = "https://discord.com/api"

func GetOwnDiscordIdentity(accessToken string) (DiscordIdentity, error) {
    const url = discordApiUrl + "/users/@me"

	payload, err := SendAuthenticatedRequest[DiscordIdentity](http.MethodGet, url, accessToken, nil)
	if err != nil {
		return DiscordIdentity{}, err
	}

    return payload, nil
}

func RequestToken(oauthCode string) (SuccesfulAuthenticationPayload, error) {
    const url = discordApiUrl + "/oauth2/token"
    var redirectUri = fmt.Sprintf("http://127.0.0.1:%s/api/login/callback", config.GetOption("port"))

    payload := URL.Values{}
    payload.Set("client_id", config.GetOption("clientId"))
    payload.Set("client_secret", config.GetOption("clientSecret"))
    payload.Set("grant_type", "authorization_code")
    payload.Set("code", oauthCode)
    payload.Set("redirect_uri", redirectUri)

    res, err := http.PostForm(url, payload)
    if err != nil {
       return SuccesfulAuthenticationPayload{}, err 
    }

    if !isOk(res.StatusCode) {
        return SuccesfulAuthenticationPayload{}, &ApiError{code: res.StatusCode, msg: res.Status}
    }

	var authPayload SuccesfulAuthenticationPayload

    defer res.Body.Close()

    bodyBytes, _ := io.ReadAll(res.Body)
    json.Unmarshal(bodyBytes, &authPayload)

    return authPayload, nil
}

func RevokeTokens(accessToken string) error {
	const url = discordApiUrl + "/oauth2/token/revoke"

	payload := URL.Values{}
    payload.Set("client_id", config.GetOption("clientId"))
    payload.Set("client_secret", config.GetOption("clientSecret"))
	payload.Set("token", accessToken)
	payload.Set("token_type_hint", "access_token")

	res, err := http.PostForm(url, payload)
	if err != nil {
		return err
	}

	if !isOk(res.StatusCode) {
		bodyBytes, _ := io.ReadAll(res.Body)
		fmt.Println(string(bodyBytes))
		return &ApiError{msg: res.Status, code: res.StatusCode}
	}

	return nil
}

func RefreshAccessToken(refreshToken string) (SuccesfulAuthenticationPayload, error) {
    const url = discordApiUrl + "/oauth2/token"

    payload := URL.Values{}
    payload.Set("client_id", config.GetOption("clientId"))
    payload.Set("client_secret", config.GetOption("clientSecret"))
    payload.Set("grant_type", "refresh_token")
    payload.Set("refresh_token", refreshToken)
    
    res, err := http.PostForm(url, payload)
    if err != nil {
       return SuccesfulAuthenticationPayload{}, err 
    }

	if !isOk(res.StatusCode) {
		return SuccesfulAuthenticationPayload{}, &ApiError{code: res.StatusCode, msg: res.Status}
	}

	var authPayload SuccesfulAuthenticationPayload

    defer res.Body.Close()
    bodyBytes, _ := io.ReadAll(res.Body)
    json.Unmarshal(bodyBytes, &authPayload)

    return authPayload, nil
}

