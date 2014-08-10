package httputil

import (
	"net/http"
	"strings"

	"code.google.com/p/xsrftoken"
)

// IsHttps is a helper function that evaluates the http.Request
// and returns True if the Request uses HTTPS. It is able to detect,
// using the X-Forwarded-Proto, if the original request was HTTPS and
// routed through a reverse proxy with SSL termination.
func IsHttps(r *http.Request) bool {
	switch {
	case r.URL.Scheme == "https":
		return true
	case strings.HasPrefix(r.Proto, "HTTPS"):
		return true
	case r.Header.Get("X-Forwarded-Proto") == "https":
		return true
	default:
		return false
	}
}

// GetScheme is a helper function that evaluates the http.Request
// and returns the scheme, HTTP or HTTPS. It is able to detect,
// using the X-Forwarded-Proto, if the original request was HTTPS
// and routed through a reverse proxy with SSL termination.
func GetScheme(r *http.Request) string {
	switch {
	case r.URL.Scheme == "https":
		return "https"
	case strings.HasPrefix(r.Proto, "HTTPS"):
		return "https"
	case r.Header.Get("X-Forwarded-Proto") == "https":
		return "https"
	default:
		return "http"
	}
}

// GetHost is a helper function that evaluates the http.Request
// and returns the hostname. It is able to detect, using the
// X-Forarded-For header, the original hostname when routed
// through a reverse proxy.
func GetHost(r *http.Request) string {
	switch {
	case len(r.Host) != 0:
		return r.Host
	case len(r.URL.Host) != 0:
		return r.URL.Host
	case len(r.Header.Get("X-Forwarded-For")) != 0:
		return r.Header.Get("X-Forwarded-For")
	case len(r.Header.Get("X-Host")) != 0:
		return r.Header.Get("X-Host")
	case len(r.Header.Get("XFF")) != 0:
		return r.Header.Get("XFF")
	default:
		return "localhost:8080"
	}
}

// GetURL is a helper function that evaluates the http.Request
// and returns the URL as a string. Only the scheme + hostname
// are included; the path is excluded.
func GetURL(r *http.Request) string {
	return GetScheme(r) + "://" + GetHost(r)
}

// GetCookie retrieves and verifies the cookie value.
func GetCookie(r *http.Request, name string) (value string) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return
	}
	value = cookie.Value
	return
}

// SetCookie writes the cookie value.
func SetCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   r.URL.Host,
		HttpOnly: true,
		Secure:   IsHttps(r),
	}

	http.SetCookie(w, &cookie)
}

// DelCookie deletes a cookie.
func DelCookie(w http.ResponseWriter, r *http.Request, name string) {
	cookie := http.Cookie{
		Name:   name,
		Value:  "deleted",
		Path:   "/",
		Domain: r.URL.Host,
		MaxAge: -1,
	}

	http.SetCookie(w, &cookie)
}

// SetXsrf writes the cookie value.
func SetXsrf(w http.ResponseWriter, r *http.Request, token, login string) {
	cookie := http.Cookie{
		Name:     "XSRF-TOKEN",
		Value:    xsrftoken.Generate(token, login, "/"),
		Path:     "/",
		Domain:   r.URL.Host,
		HttpOnly: false,
		Secure:   IsHttps(r),
	}

	http.SetCookie(w, &cookie)
}

// CheckXsrf verifies the xsrf value.
func CheckXsrf(r *http.Request, token, login string) bool {
	if r.Method == "GET" {
		return true
	}
	return xsrftoken.Valid(
		r.Header.Get("X-XSRF-TOKEN"), token, login, "/")
}
