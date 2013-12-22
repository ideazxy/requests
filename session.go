package requests

import (
	"net/http/cookiejar"
	"net/url"
	"time"
)

type Session struct {
	Client    *Client
	userAgent string
}

func NewSession() *Session {
	client := NewClient()
	client.SetTimeout(60*time.Second, 60*time.Second)
	// Jar will store cookies
	client.Jar, _ = cookiejar.New(nil)

	return &Session{
		client,
		"Go Requests",
	}
}

func (s *Session) Request(method, url string) *HttpRequest {
	r := Request(method, url)
	r.Client = s.Client
	return r
}

func (s *Session) Get(url string) *HttpRequest {
	return s.Request("GET", url)
}

func (s *Session) Post(url string, data interface{}, bodyType string) *HttpRequest {
	return s.Request("POST", url).SetBody(data, bodyType)
}

func (s *Session) PostForm(u string, data url.Values) *HttpRequest {
	return s.Request("POST", u).SetBody(data.Encode(), "application/x-www-form-urlencoded")
}

func (s *Session) Put(url string, data interface{}, bodyType string) *HttpRequest {
	return s.Request("PUT", url).SetBody(data, bodyType)
}

func (s *Session) Head(url string) *HttpRequest {
	return s.Request("HEAD", url)
}

func (s *Session) Options(url string) *HttpRequest {
	return s.Request("OPTIONS", url)
}

func (s *Session) Patch(url string, data interface{}) *HttpRequest {
	return s.Request("PATCH", url)
}

func (s *Session) Delete(url string) *HttpRequest {
	return s.Request("DELETE", url)
}
