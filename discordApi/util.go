package discordapi

func isOk(statusCode int) bool {
    return statusCode >= 200 && statusCode <= 299
}
