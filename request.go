package requests

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpRequest struct {
	URL              string
	Req              *http.Request
	Params           url.Values
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns "".
// To access multiple values of a key, access the map directly
// with CanonicalHeaderKey.
func (r *HttpRequest) Header(key string) string {
	return r.Req.Header.Get(key)
}

// Set sets the header entries associated with key to
// the single element value.  It replaces any existing
// values associated with key.
func (r *HttpRequest) SetHeader(key, value string) *HttpRequest {
	r.Req.Header.Set(key, value)
	return r
}

// Add adds the key, value pair to the header.
// It appends to any existing values associated with key.
func (r *HttpRequest) AddHeader(key, value string) *HttpRequest {
	r.Req.Header.Add(key, value)
	return r
}

// Del deletes the values associated with key.
func (r *HttpRequest) DelHeader(key string) *HttpRequest {
	r.Req.Header.Del(key)
	return r
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (r *HttpRequest) Param(key string) string {
	if r.Params == nil {
		return ""
	}
	v, ok := r.Params[key]
	if !ok || len(v) == 0 {
		return ""
	}
	return v[0]
}
	
func (r *HttpRequest) SetParam(key, value string) *HttpRequest {
	r.Params.Set(key, value)
	return r
}

func (r *HttpRequest) AddParam(key, value string) *HttpRequest {
	r.Params.Add(key, value)
	return r
}

func (r *HttpRequest) DelParam(key string) *HttpRequest {
	r.Params.Del(key)
	return r
}

// Cookies parses and returns the HTTP cookies sent with the request.
func (r *HttpRequest) Cookies() []*http.Cookie {
	return r.Req.Cookies()
}

// Cookie returns the named cookie provided in the request or
// ErrNoCookie if not found.
func (r *HttpRequest) Cookie(name string) (*http.Cookie, error) {
	return r.Req.Cookie(name)
}

// AddCookie does not attach more than one Cookie header field.  That
// means all cookies, if any, are written into the same line,
// separated by semicolon.
func (r *HttpRequest) AddCookie(c *http.Cookie) *HttpRequest {
	r.Req.AddCookie(c)
	return r
}

func (r *HttpRequest) Timeout(connectTimeout, readWriteTimeout time.Duration) *HttpRequest {
	r.ConnectTimeout = connectTimeout
	r.ReadWriteTimeout = readWriteTimeout
	return r
}

func (r *HttpRequest) SetBody(data interface{}, bodyType string) *HttpRequest {
	rd, ok := data.(io.Reader)
	if !ok && data != nil {
		switch v := data.(type) {
		case string:
			rd = bytes.NewBufferString(v)
		case []byte:
			rd = bytes.NewBuffer(v)
		}
	}
	rc, ok := rd.(io.ReadCloser)
	if !ok && rd != nil {
		rc = ioutil.NopCloser(rd)
	}
	r.Req.Body = rc

	if rd != nil {
		r.SetHeader("Content-Type", bodyType)
		switch v := rd.(type) {
		case *bytes.Buffer:
			r.Req.ContentLength = int64(v.Len())
		case *bytes.Reader:
			r.Req.ContentLength = int64(v.Len())
		case *strings.Reader:
			r.Req.ContentLength = int64(v.Len())
		}
	}

	return r
}

func (r *HttpRequest) encodeUrl() error {
	_, err := url.Parse(r.URL)
	if err != nil {
		return err
	}

	r.Req.URL.RawQuery = r.Params.Encode()
	
	if !r.Req.URL.IsAbs() {
		r.Req.URL.Scheme = "http"
	}
	return nil
}

func (r *HttpRequest) Send() (*HttpResponse, error) {
	err := r.encodeUrl()
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: TimeoutDialer(r.ConnectTimeout, r.ReadWriteTimeout),
		},
	}
	resp, err := client.Do(r.Req)
	if err != nil {
		return nil, err
	}
	return NewResponse(resp), nil
}

func NewRequest(method, rawurl string) *HttpRequest {
	// URL will be validated in send():
	req, _ := http.NewRequest(method, rawurl, nil)

	r := &HttpRequest{
		rawurl,
		req,
		make(url.Values),
		60 * time.Second,
		60 * time.Second}
	if r.Req.URL != nil {
		r.Params = r.Req.URL.Query()
	}
	return r
}

func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (net.Conn, error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}
