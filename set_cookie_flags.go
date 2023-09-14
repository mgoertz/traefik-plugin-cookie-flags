// Package traefik_plugin_cookie_flags a traefik plugin adding flags to cookies in the response.
package traefik_plugin_cookie_flags //nolint

import (
	"context"
	"net/http"
)

const setCookieHeader string = "Set-Cookie"

// Config the plugin configuration.
type Config struct {
	SameSite string `json:"sameSite,omitempty" toml:"sameSite,omitempty" yaml:"sameSite,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// CookieFlagger an plugin with a possible configuration.
type CookieFlagger struct {
	next     http.Handler
	name     string
	sameSite string
}

// New creates new instance of the plugin.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &CookieFlagger{
		name:     name,
		next:     next,
		sameSite: config.SameSite,
	}, nil
}

func (p *CookieFlagger) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	_sameSite := http.SameSiteDefaultMode //nolint

	switch p.sameSite {
	case "lax":
		_sameSite = http.SameSiteLaxMode
	case "strict":
		_sameSite = http.SameSiteStrictMode
	case "none":
		_sameSite = http.SameSiteNoneMode
	default:
		_sameSite = http.SameSiteDefaultMode
	}

	myWriter := &responseWriter{
		writer:   rw,
		sameSite: _sameSite,
	}

	p.next.ServeHTTP(myWriter, req)
}

type responseWriter struct {
	writer   http.ResponseWriter
	sameSite http.SameSite
}

func (r *responseWriter) Header() http.Header {
	return r.writer.Header()
}

func (r *responseWriter) Write(bytes []byte) (int, error) {
	return r.writer.Write(bytes)
}

func (r *responseWriter) WriteHeader(statusCode int) {
	// Get the cookies
	headers := r.writer.Header()
	req := http.Response{Header: headers}
	cookies := req.Cookies()

	// Delete set-cookie headers
	r.writer.Header().Del(setCookieHeader)

	// Add new cookie with modified path
	for _, cookie := range cookies {
		cookie.HttpOnly = true
		cookie.Secure = true
		cookie.SameSite = r.sameSite
		http.SetCookie(r, cookie)
	}

	r.writer.WriteHeader(statusCode)
}
