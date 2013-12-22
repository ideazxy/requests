package requests

import (
	"testing"
	"strings"
)

var TestHost = "http://httpbin.org"

type response struct {
	Data string `json:"data"`
	Args map[string]string `json:"args"`
	Form map[string]string `json:"form"`
	Headers map[string]string `json:"headers"`
	Url string `json:"url"`
	Origin string `json:"origin"`
}

func TestGet(t *testing.T) {
	resp, err := Get(TestHost + "/get").Send()
	if err != nil {
		t.Fatalf("Get request failed: %v", err)
	}
	body, err := resp.Text()
	if err != nil {
		t.Fatalf("Text parsing failed: %v", err)
	}
	if !strings.Contains(body, `"url": "http://httpbin.org/get"`) {
		t.Errorf(`Get data: %s, expected to contain: "url": "http://httpbin.org/get"`, body)
	}
}

func TestPost(t *testing.T) {
	body := "content"
	resp, err := Post(TestHost + "/post", body, "text/plain").Send()
	if err != nil {
		t.Fatalf("Post request failed: %v", err)
	}
	var respJson response
	err = resp.Json(&respJson)
	if err != nil {
		t.Fatalf("Json parse failed: %v", err)
	}
	if respJson.Data != body {
		t.Errorf("Post data: %s, expected: %s", respJson.Data, body)
		t.Logf("Response: %v", respJson)
	}
}