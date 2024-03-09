package cookie

import (
	"net/http"
	"time"
)

func IsUnexpired(cookie *http.Cookie) bool {
	if cookie.Expires.After(time.Now()) {
		return false
	}

	return true
}

func SetWithExpiration(w http.ResponseWriter, cookie http.Cookie, duration time.Duration) {
	cookie.Expires = time.Now().Add(duration)
	http.SetCookie(w, &cookie)
}

func Clear(w http.ResponseWriter, cookie *http.Cookie) {
	// Invalidate the cookie in standard HTTP way
	cookie.MaxAge = -1 
	http.SetCookie(w, cookie)
}
