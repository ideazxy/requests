package requests

import (
	"net/http"
	"io/ioutil"
)

type HttpResponse struct {
	resp *http.Response
	Status int
	content []byte
	consumed bool
}

func NewResponse(resp *http.Response) *HttpResponse {
	var r HttpResponse
	r.resp = resp
	r.Status = resp.StatusCode
	return &r
}

func (r *HttpResponse) Content() ([]byte, error) {
	if r.consumed {
		return r.content, nil
	}
	defer r.resp.Body.Close()
	data, err := ioutil.ReadAll(r.resp.Body)
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
