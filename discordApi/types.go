package discordapi

type DiscordIdentity struct {
    UserId string `json:"id"`
    Username string `json:"username"`
    Discriminator string `json:"discriminator"`
    Avatar string `json:"avatar"`
}

type SuccesfulAuthenticationPayload struct {
    AccessToken string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn int `json:"expires_in"`
}

type ApiError struct {
    code int
    msg string
}

func (e *ApiError) Error() string {
    return e.msg
}
