package discordapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func SendAuthenticatedRequest[T any](method, url, accessToken string, data any) (T, error) {
    req, _ := http.NewRequest(http.MethodGet, url, nil)
    req.Header.Add("Authorization", "Bearer " + accessToken)

    var payload = new(T)

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return *payload, err
    }

    if !isOk(res.StatusCode) {
        return *payload, &ApiError{code: res.StatusCode, msg: res.Status}
    }

    defer res.Body.Close()
    bodyBytes, _ := io.ReadAll(res.Body)

    if err = json.Unmarshal(bodyBytes, payload); err != nil {
        return *payload, err
    }

    return *payload, nil 
}

func isOk(statusCode int) bool {
    return statusCode >= 200 && statusCode <= 299
}
