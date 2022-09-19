package v1

import "net/http"

func isContentTypeMatch(r *http.Request, contentType string) bool {
	if r.Header.Get("Content-Type") != contentType {
		return false
	}
	return true
}
