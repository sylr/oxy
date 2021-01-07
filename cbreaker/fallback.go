package cbreaker

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"abstraction.fr/oxy/v2/utils"
)

// Response response model
type Response struct {
	StatusCode  int
	ContentType string
	Body        []byte
}

// ResponseFallback fallback response handler
type ResponseFallback struct {
	r   Response
	log utils.Logger
}

// NewResponseFallbackWithLogger creates a new ResponseFallback
func NewResponseFallbackWithLogger(r Response, l utils.Logger) (*ResponseFallback, error) {
	if r.StatusCode == 0 {
		return nil, fmt.Errorf("response code should not be 0")
	}
	return &ResponseFallback{r: r, log: l}, nil
}

// NewResponseFallback creates a new ResponseFallback
func NewResponseFallback(r Response) (*ResponseFallback, error) {
	return NewResponseFallbackWithLogger(r, &utils.DefaultLogger{})
}

func (f *ResponseFallback) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if f.r.ContentType != "" {
		w.Header().Set("Content-Type", f.r.ContentType)
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(f.r.Body)))
	w.WriteHeader(f.r.StatusCode)
	_, err := w.Write(f.r.Body)
	if err != nil {
		f.log.Errorf("fallback/response: failed to write response, err: %v", err)
	}
}

// Redirect redirect model
type Redirect struct {
	URL          string
	PreservePath bool
}

// RedirectFallback fallback redirect handler
type RedirectFallback struct {
	r   Redirect
	u   *url.URL
	log utils.Logger
}

// NewRedirectFallbackWithLogger creates a new RedirectFallback
func NewRedirectFallbackWithLogger(r Redirect, l utils.Logger) (*RedirectFallback, error) {
	u, err := url.ParseRequestURI(r.URL)
	if err != nil {
		return nil, err
	}
	return &RedirectFallback{r: r, u: u, log: l}, nil
}

// NewRedirectFallback creates a new RedirectFallback
func NewRedirectFallback(r Redirect) (*RedirectFallback, error) {
	return NewRedirectFallbackWithLogger(r, &utils.DefaultLogger{})
}

func (f *RedirectFallback) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	location := f.u.String()
	if f.r.PreservePath {
		location += req.URL.Path
	}

	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusFound)
	_, err := w.Write([]byte(http.StatusText(http.StatusFound)))
	if err != nil {
		f.log.Errorf("fallback/redirect: failed to write response, err: %v", err)
	}
}
