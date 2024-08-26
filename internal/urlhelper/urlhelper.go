package urlhelper

import (
	"fmt"
	"net/http"
)

func GetFullURLOverridePath(r *http.Request, overridePath ...string) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if len(overridePath) != 0 {
		return fmt.Sprintf("%s://%s%s", scheme, r.Host, overridePath[0])
	}
	return fmt.Sprintf("%s://%s%s?%s#%s", scheme, r.Host, r.URL.Path, r.URL.RawQuery, r.URL.Fragment)
}
