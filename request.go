package requests

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpRequest struct {
	Client *Client
	rawUrl string
	Req    *http.Request
	Params url.Values
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
	if r.Client == nil {
		r.Client = NewClient()
	}
	r.Client.SetTimeout(connectTimeout, readWriteTimeout)
	return r
}

func (r *HttpRequest) AllowRedirects(allow bool) *HttpRequest {
	if r.Client == nil {
		r.Client = NewClient()
	}
	if !allow {
		r.Client.ForbidRedirects()
	}
	return r
}

func (r *HttpRequest) SetRedirectMax(count int) *HttpRequest {
	if r.Client == nil {
		r.Client = NewClient()
	}
	r.Client.RedirectMax = count
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

func (r *HttpRequest) Query() url.Values {
	return r.Req.URL.Query()
}

func (r *HttpRequest) Path() string {
	return r.Req.URL.Path
}

func (r *HttpRequest) Url() string {
	return r.Req.URL.String()
}

func (r *HttpRequest) encodeUrl() error {
	_, err := url.Parse(r.rawUrl)
	if err != nil {
		return err
	}

	r.Req.URL.RawQuery = r.Params.Encode()

	if !r.Req.URL.IsAbs() {
		r.Req.URL.Scheme = "http"
	}
	return nil
}

func (r *HttpRequest) Prepare() error {
	return r.encodeUrl()
}

func (r *HttpRequest) Send() (*HttpResponse, error) {
	err := r.Prepare()
	if err != nil {
		return nil, err
	}

	if r.Client == nil {
		r.Client = NewClient()
	}
	return r.Client.Do(r.Req)
}

func NewRequest(method, rawurl string) *HttpRequest {
	// URL will be validated in send():
	req, _ := http.NewRequest(strings.ToUpper(method), rawurl, nil)

	r := &HttpRequest{
		nil,
		rawurl,
		req,
		make(url.Values)}
	if r.Req.URL != nil {
		r.Params = r.Req.URL.Query()
	}
	return r
}
