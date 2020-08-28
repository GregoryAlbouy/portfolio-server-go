package utl

import (
	"net/http"
)

// RequestIP returns the IP of input request. It first checks for proxies
// through X-FORWARDED-FOR header and returns its value if found, if not
// it returns the remote address.
func RequestIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
