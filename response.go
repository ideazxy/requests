package requests

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type HttpResponse struct {
	resp     *http.Response
	Status   int
	content  []byte
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
