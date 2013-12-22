package requests

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

type RedirectForbiddenError struct {
}

func (e *RedirectForbiddenError) Error() string {
	return "no redirects allowed"
}

func NewRedirectForbiddenError() *RedirectForbiddenError {
	return &RedirectForbiddenError{}
}

type Client struct {
	http.Client
	forbidRedirects bool
	RedirectMax    int
	History        []*http.Response
}

func NewClient() *Client {
	c := &Client{}
	c.RedirectMax = 10
	return c
}

func (c *Client) Do(req *http.Request) (*HttpResponse, error) {
	if c.Client.CheckRedirect == nil {
		c.Client.CheckRedirect = c.defaultCheckRedirect()
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		if urlErr, ok := err.(*url.Error); ok {
			_, ok := urlErr.Err.(*RedirectForbiddenError)
			if ok {
				response := NewResponse(resp)
				response.RedirectForbidden = true
				return response, nil
			}
		}
		return nil, err
	}
	return NewResponse(resp), nil
}

func (c *Client) ForbidRedirects() {
	c.forbidRedirects = true
}

func (c *Client) SetTimeout(connectTimeout, readWriteTimeout time.Duration) {
	c.Transport = &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(netw, addr, connectTimeout)
			if err != nil {
				return nil, err
			}
			conn.SetDeadline(time.Now().Add(readWriteTimeout))
			return conn, nil
		},
	}
}

func (c *Client) defaultCheckRedirect() func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if c.forbidRedirects {
			return NewRedirectForbiddenError()
		}
		if len(via) >= c.RedirectMax {
			return errors.New(fmt.Sprintf("stopped after %d redirects", c.RedirectMax))
		}
		return nil
	}
}

func getTimeoutDialer(connectTimeout, readWriteTimeout time.Duration) func(net, addr string) (net.Conn, error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, connectTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(readWriteTimeout))
		return conn, nil
	}
}
