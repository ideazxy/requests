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
	params           map[string]string
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
}

func (r *HttpRequest) Header(key, value string) *HttpRequest {
	r.Req.Header.Set(key, value)
	return r
}

func (r *HttpRequest) Param(key, value string) *HttpRequest {
	r.params[key] = value
	return r
}

func (r *HttpRequest) Timeout(connectTimeout, readWriteTimeout time.Duration) *HttpRequest {
	r.ConnectTimeout = connectTimeout
	r.ReadWriteTimeout = readWriteTimeout
	return r
}

func (r *HttpRequest) Body(data interface{}) *HttpRequest {
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
	newUrl, err := url.Parse(r.URL)
	if err != nil {
		return err
	}
	if r.params != nil && len(r.params) > 0 {
		params := newUrl.Query()
		for k, v := range r.params {
			params.Add(k, v)
		}
		newUrl.RawQuery = params.Encode()
	}
	if !newUrl.IsAbs() {
		newUrl.Scheme = "http"
	}
	r.Req.URL = newUrl
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

func NewRequest(method, url string, data interface{}) *HttpRequest {
	// URL will be validated in send():
	req, _ := http.NewRequest(method, url, nil)

	r := &HttpRequest{
		url,
		req,
		make(map[string]string),
		60 * time.Second,
		60 * time.Second}
	r.Body(data)
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
