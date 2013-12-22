package requests

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	// Transport specifies the mechanism by which individual
	// HTTP requests are made.
	// If nil, DefaultTransport is used.
	Transport http.RoundTripper
	// If CheckRedirect is nil, the Client uses its default policy,
	// which is to stop after 10 consecutive requests.
	CheckRedirect func(req *http.Request, via []*http.Request) error

	// Jar specifies the cookie jar.
	// If Jar is nil, cookies are not sent in requests and ignored
	// in responses.
	Jar http.CookieJar

	AllowRedirects bool
	RedirectMax    int
	History        []*http.Response
}

func NewClient() *Client {
	c := &Client{}
	c.RedirectMax = 10
	return c
}

// Copied from net/http/client.go
func (c *Client) send(req *http.Request) (*http.Response, error) {
	if c.Jar != nil {
		for _, cookie := range c.Jar.Cookies(req.URL) {
			req.AddCookie(cookie)
		}
	}
	resp, err := send(req, c.Transport)
	if err != nil {
		return nil, err
	}
	if c.Jar != nil {
		if rc := resp.Cookies(); len(rc) > 0 {
			c.Jar.SetCookies(req.URL, rc)
		}
	}
	return resp, err
}

// Copied from net/http/client.go
func (c *Client) Do(req *http.Request) (resp *http.Response, err error) {
	if req.Method == "GET" || req.Method == "HEAD" {
		return c.doFollowingRedirects(req, shouldRedirectGet)
	}
	if req.Method == "POST" || req.Method == "PUT" {
		return c.doFollowingRedirects(req, shouldRedirectPost)
	}
	return c.send(req)
}

// Copied from net/http/client.go
// Caller should close resp.Body when done reading from it.
func send(req *http.Request, t http.RoundTripper) (resp *http.Response, err error) {
	if t == nil {
		t = http.DefaultTransport
		if t == nil {
			err = errors.New("http: no Client.Transport or DefaultTransport")
			return
		}
	}

	if req.URL == nil {
		return nil, errors.New("http: nil Request.URL")
	}

	if req.RequestURI != "" {
		return nil, errors.New("http: Request.RequestURI can't be set in client requests.")
	}

	// Most the callers of send (Get, Post, et al) don't need
	// Headers, leaving it uninitialized.  We guarantee to the
	// Transport that this has been initialized, though.
	if req.Header == nil {
		req.Header = make(http.Header)
	}

	if u := req.URL.User; u != nil {
		req.Header.Set("Authorization", "Basic "+base64.URLEncoding.EncodeToString([]byte(u.String())))
	}
	resp, err = t.RoundTrip(req)
	if err != nil {
		if resp != nil {
			log.Printf("RoundTripper returned a response & error; ignoring response")
		}
		return nil, err
	}
	return resp, nil
}

// Copied from net/http/client.go
func shouldRedirectGet(statusCode int) bool {
	switch statusCode {
	case http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther, http.StatusTemporaryRedirect:
		return true
	}
	return false
}

// Copied from net/http/client.go
func shouldRedirectPost(statusCode int) bool {
	switch statusCode {
	case http.StatusFound, http.StatusSeeOther:
		return true
	}
	return false
}

func (c *Client) doFollowingRedirects(ireq *http.Request, shouldRedirect func(int) bool) (resp *http.Response, err error) {
	var base *url.URL
	redirectChecker := c.CheckRedirect
	if redirectChecker == nil {
		redirectChecker = c.defaultCheckRedirect()
	}
	var via []*http.Request
	c.History = nil

	if ireq.URL == nil {
		return nil, errors.New("http: nil Request.URL")
	}

	req := ireq
	urlStr := "" // next relative or absolute URL to fetch (after first request)
	redirectFailed := false
	for redirect := 0; ; redirect++ {
		if redirect != 0 {
			req = new(http.Request)
			req.Method = ireq.Method
			if ireq.Method == "POST" || ireq.Method == "PUT" {
				req.Method = "GET"
			}
			req.Header = make(http.Header)
			req.URL, err = base.Parse(urlStr)
			if err != nil {
				break
			}
			if len(via) > 0 {
				// Add the Referer header.
				lastReq := via[len(via)-1]
				if lastReq.URL.Scheme != "https" {
					req.Header.Set("Referer", lastReq.URL.String())
				}

				err = redirectChecker(req, via)
				if err != nil {
					redirectFailed = true
					break
				}
			}
		}

		urlStr = req.URL.String()
		if resp, err = c.send(req); err != nil {
			break
		}

		if shouldRedirect(resp.StatusCode) {
			if !c.AllowRedirects {
				return
			}
			resp.Body.Close()
			c.History = append(c.History, resp)
			if urlStr = resp.Header.Get("Location"); urlStr == "" {
				err = errors.New(fmt.Sprintf("%d response missing Location header", resp.StatusCode))
				break
			}
			base = req.URL
			via = append(via, req)
			continue
		}

		return
	}

	method := ireq.Method
	urlErr := &url.Error{
		Op:  method[0:1] + strings.ToLower(method[1:]),
		URL: urlStr,
		Err: err,
	}

	if redirectFailed {
		// Special case for Go 1 compatibility: return both the response
		// and an error if the CheckRedirect function failed.
		// See http://golang.org/issue/3795
		return resp, urlErr
	}

	if resp != nil {
		resp.Body.Close()
	}
	return nil, urlErr
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
