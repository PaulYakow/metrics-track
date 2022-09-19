package v1

import "net/http"

func isContentTypeMatch(r *http.Request, contentType string) bool {
	return r.Header.Get("Content-Type") == contentType
}
