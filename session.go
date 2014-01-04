package requests

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type Session struct {
	Client    *Client
	UserAgent string
	lastReq   *HttpRequest
}

func NewSession() *Session {
	client := NewClient()
	client.SetTimeout(60*time.Second, 60*time.Second)
	// Jar will store cookies
	client.Jar, _ = cookiejar.New(nil)

	return &Session{
		client,
		"Go Requests",
		nil,
	}
}

func (s *Session) Request(method, url string) *HttpRequest {
	r := Request(method, url)
	r.SetHeader("User-Agent", s.UserAgent)
	r.Client = s.Client
	s.lastReq = r
	return r
}

func (s *Session) Cookies(url *url.URL) (cookies []*http.Cookie) {
	if url == nil {
		if s.lastReq == nil {
			return cookies
		}
		url = s.lastReq.Req.URL
	}
	return s.Client.Jar.Cookies(url)
}

func (s *Session) Cookie(key string, url *url.URL) (cookie *http.Cookie) {
	if cookies := s.Cookies(url); len(cookies) > 0 {
		for _, c := range cookies {
			if c.Name == key {
				return c
			}
		}
	}
	return
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
