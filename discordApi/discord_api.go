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

    req, _ := http.NewRequest(http.MethodGet, url, nil)
    req.Header.Add("Authorization", "Bearer " + accessToken)

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return DiscordIdentity{}, err
    }

    if !isOk(res.StatusCode) {
        return DiscordIdentity{}, &ApiError{code: res.StatusCode, msg: res.Status}
    }

    var discordIdentity DiscordIdentity

    bodyBytes, err := io.ReadAll(res.Body)
    if err != nil {
        return DiscordIdentity{}, err
    }

    fmt.Println(string(bodyBytes))

    if err = json.Unmarshal(bodyBytes, &discordIdentity); err != nil {
        return DiscordIdentity{}, err
    }

    return discordIdentity, nil
}

func RequestToken(oauthCode string) (accessToken, refreshToken string, err error) {
    const url = discordApiUrl + "/oauth2/token"
    var redirectUri = fmt.Sprintf("http://127.0.0.1:%s/api/login", config.GetOption("port"))

    payload := URL.Values{}
    payload.Set("client_id", config.GetOption("clientId"))
    payload.Set("client_secret", config.GetOption("clientSecret"))
    payload.Set("grant_type", "authorization_code")
    payload.Set("code", oauthCode)
    payload.Set("redirect_uri", redirectUri)

    res, err := http.PostForm(url, payload)
    if err != nil {
       return "", "", err 
    }

    if !isOk(res.StatusCode) {
        return "", "", &ApiError{code: res.StatusCode, msg: res.Status}
    }

    // Body always has to be closed
    defer res.Body.Close()

    var succesfulAuth SuccesfulAuthentication

    bodyBytes, _ := io.ReadAll(res.Body)
    json.Unmarshal(bodyBytes, &succesfulAuth)

    accessToken = succesfulAuth.AccessToken
    refreshToken = succesfulAuth.RefreshToken
    return
}

func RefreshAccessToken(refreshToken string) (accessToken string, err error) {
    const url = discordApiUrl + "/oauth2/token"

    payload := URL.Values{}
    payload.Set("client_id", config.GetOption("clientId"))
    payload.Set("client_secret", config.GetOption("clientSecret"))
    payload.Set("grant_type", "refresh_token")
    payload.Set("refresh_token", refreshToken)
    
    res, err := http.PostForm(url, payload)
    if err != nil {
       return "", err 
    }

    if !isOk(res.StatusCode) {
        return "", &ApiError{code: res.StatusCode, msg: res.Status}
    }

    // Body always has to be closed
    defer res.Body.Close()

    var succesfulAuth SuccesfulAuthentication

    bodyBytes, _ := io.ReadAll(res.Body)
    json.Unmarshal(bodyBytes, &succesfulAuth)

    accessToken = succesfulAuth.AccessToken
    return
}
