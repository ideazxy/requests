package requests

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HttpResponse struct {
	Resp              *http.Response
	StatusCode        int
	RedirectForbidden bool
	cookies           map[string]*http.Cookie
	content           []byte
	consumed          bool
}

func NewResponse(resp *http.Response) *HttpResponse {
	var r HttpResponse
	r.Resp = resp
	r.StatusCode = resp.StatusCode

	r.cookies = make(map[string]*http.Cookie)
	cookies := r.Resp.Cookies()
	if cookies != nil {
		for _, v := range cookies {
			r.cookies[v.Name] = v
		}
	}

	return &r
}

func (r *HttpResponse) Content() ([]byte, error) {
	if r.consumed {
		return r.content, nil
	}
	defer r.Resp.Body.Close()
	data, err := ioutil.ReadAll(r.Resp.Body)
	if err != nil {
		return nil, err
	}
	r.consumed = true
	r.content = data
	return r.content, nil
}

func (r *HttpResponse) Text() (string, error) {
	b, err := r.Content()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (r *HttpResponse) Json(v interface{}) error {
	b, err := r.Content()
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	return nil
}

func (r *HttpResponse) Xml(v interface{}) error {
	b, err := r.Content()
	if err != nil {
		return err
	}
	err = xml.Unmarshal(b, v)
	if err != nil {
		return err
	}
	return nil
}

// Cookies parses and returns the cookies set in the Set-Cookie headers.
func (r *HttpResponse) Cookies() []*http.Cookie {
	return r.Resp.Cookies()
}

func (r *HttpResponse) Cookie(name string) *http.Cookie {
	return r.cookies[name]
}

// Location returns the URL of the response's "Location" header,
// if present.  Relative redirects are resolved relative to
// the Response's Request.  ErrNoLocation is returned if no
// Location header is present.
func (r *HttpResponse) Location() (*url.URL, error) {
	return r.Resp.Location()
}
