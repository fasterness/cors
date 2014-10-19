package cors

import (
	"fmt"
	"net/http"
	"strings"
)

type CorsHandler struct {
	ALLOWED_METHODS   []string
	ALLOWED_ORIGINS   []string
	ALLOWED_HEADERS   []string
	EXPOSED_HEADERS   []string
	ALLOW_CREDENTIALS string
	MAX_AGE           float64
	handler           http.Handler
}

func New(handler http.Handler) *CorsHandler {
	return &CorsHandler{
		ALLOWED_METHODS:   []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"},
		ALLOWED_ORIGINS:   []string{"*"},
		ALLOWED_HEADERS:   []string{"Content-Type"},
		EXPOSED_HEADERS:   []string{"Content-Type"},
		ALLOW_CREDENTIALS: "true",
		MAX_AGE:           0,
		handler:           handler,
	}
}
func (cors *CorsHandler) AllowOrigin(origin string) {
	if origin == "*" {
		cors.ALLOWED_ORIGINS = []string{"*"}
	}
	for i := 0; i < len(cors.ALLOWED_ORIGINS); i++ {
		if origin == cors.ALLOWED_ORIGINS[i] {
			return
		}
	}
	cors.ALLOWED_ORIGINS = append(cors.ALLOWED_ORIGINS, origin)
}
func (cors *CorsHandler) AllowMethod(method string) {
	method = strings.ToUpper(method)
	for i := 0; i < len(cors.ALLOWED_METHODS); i++ {
		if method == cors.ALLOWED_METHODS[i] {
			return
		}
	}
	cors.ALLOWED_METHODS = append(cors.ALLOWED_METHODS, method)
}
func (cors *CorsHandler) AllowHeader(header string) {
	for i := 0; i < len(cors.ALLOWED_HEADERS); i++ {
		if header == cors.ALLOWED_HEADERS[i] {
			return
		}
	}
	cors.ALLOWED_HEADERS = append(cors.ALLOWED_HEADERS, header)

}

func (cors *CorsHandler) ExposeHeader(header string) {
	for _, exposedHeader := range cors.EXPOSED_HEADERS {
		if header == exposedHeader {
			return
		}
	}
	cors.EXPOSED_HEADERS = append(cors.EXPOSED_HEADERS, header)
}

func (cors *CorsHandler) AllowCredentials(creds bool) {
	if creds {
		cors.ALLOW_CREDENTIALS = "true"
	} else {
		cors.ALLOW_CREDENTIALS = "false"
	}
}
func (cors *CorsHandler) RemoveOrigin(origin string) {
	for i := 0; i < len(cors.ALLOWED_ORIGINS); i++ {
		if origin == cors.ALLOWED_ORIGINS[i] {
			cors.ALLOWED_ORIGINS = cors.ALLOWED_ORIGINS[:i+copy(cors.ALLOWED_ORIGINS[i:], cors.ALLOWED_ORIGINS[i+1:])]
		}
	}

}
func (cors *CorsHandler) RemoveMethod(method string) {
	method = strings.ToUpper(method)
	for i := 0; i < len(cors.ALLOWED_METHODS); i++ {
		if method == cors.ALLOWED_METHODS[i] {
			cors.ALLOWED_METHODS = cors.ALLOWED_METHODS[:i+copy(cors.ALLOWED_METHODS[i:], cors.ALLOWED_METHODS[i+1:])]
		}
	}
}
func (cors *CorsHandler) RemoveHeader(header string) {
	for i := 0; i < len(cors.ALLOWED_HEADERS); i++ {
		if header == cors.ALLOWED_HEADERS[i] {
			cors.ALLOWED_HEADERS = cors.ALLOWED_HEADERS[:i+copy(cors.ALLOWED_HEADERS[i:], cors.ALLOWED_HEADERS[i+1:])]
		}
	}
}

func (cors *CorsHandler) RemoveExposedHeader(header string) {
	for i, elem := range cors.EXPOSED_HEADERS {
		if header == elem {
			cors.EXPOSED_HEADERS = cors.EXPOSED_HEADERS[:i+copy(cors.EXPOSED_HEADERS[i:], cors.EXPOSED_HEADERS[i+1:])]
		}
	}
}

func (cors *CorsHandler) IsOriginAllowed(origin string) bool {
	for i := 0; i < len(cors.ALLOWED_ORIGINS); i++ {
		if "*" == cors.ALLOWED_ORIGINS[i] {
			return true
		} else if origin == cors.ALLOWED_ORIGINS[i] {
			return true
		}
	}
	return false
}
func (cors *CorsHandler) IsMethodAllowed(method string) bool {
	method = strings.ToUpper(method)
	for i := 0; i < len(cors.ALLOWED_METHODS); i++ {
		if method == cors.ALLOWED_METHODS[i] {
			return true
		}
	}
	return false
}
func (cors *CorsHandler) IsHeaderAllowed(header string) bool {
	for i := 0; i < len(cors.ALLOWED_HEADERS); i++ {
		if header == cors.ALLOWED_HEADERS[i] {
			return true
		}
	}
	return false
}
func (cors *CorsHandler) AllowedMethods() string {
	methods := ""
	for i := 0; i < len(cors.ALLOWED_METHODS); i++ {
		if methods == "" {
			methods = cors.ALLOWED_METHODS[i]
		} else {
			methods = fmt.Sprintf("%s, %s", methods, cors.ALLOWED_METHODS[i])
		}
	}
	return methods
}
func (cors *CorsHandler) AllowedHeaders() string {
	headers := ""
	for i := 0; i < len(cors.ALLOWED_HEADERS); i++ {
		if headers == "" {
			headers = cors.ALLOWED_HEADERS[i]
		} else {
			headers = fmt.Sprintf("%s, %s", headers, cors.ALLOWED_HEADERS[i])
		}
	}
	return headers
}

func (cors *CorsHandler) ExposedHeaders() string {
	headers := ""
	for _, header := range cors.EXPOSED_HEADERS {
		if headers == "" {
			headers = header
		} else {
			headers = fmt.Sprintf("%s, %s", headers, header)
		}
	}
	return headers
}

func (cors *CorsHandler) SetMaxAge(age float64) {
	cors.MAX_AGE = age
}
func (cors *CorsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	origin := req.Header.Get("Origin")
	if origin != "" && cors.IsOriginAllowed(origin) {
		w.Header().Add("Access-Control-Allow-Origin", origin)
		w.Header().Add("Access-Control-Allow-Methods", cors.AllowedMethods())
		w.Header().Add("Access-Control-Allow-Headers", cors.AllowedHeaders())
		w.Header().Add("Access-Control-Expose-Headers", cors.ExposedHeaders())
		w.Header().Add("Access-Control-Allow-Credentials", cors.ALLOW_CREDENTIALS)
		if cors.MAX_AGE > 0 {
			w.Header().Add("Access-Control-Max-Age", fmt.Sprintf("%9.f", cors.MAX_AGE))
		}
	}
	cors.handler.ServeHTTP(w, req)
}
