package discordapi

type DiscordIdentity struct {
    UserId string `json:"id"`
    Username string `json:"username"`
    Discriminator string `json:"discriminator"`
    Avatar string `json:"avatar"`
}

type SuccesfulAuthentication struct {
    AccessToken string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

type ApiError struct {
    code int
    msg string
}

func (e *ApiError) Error() string {
    return e.msg
}
