package requests

import (
	"strings"
	"testing"
)

var TestHost = "http://httpbin.org"

type response struct {
	Data    string            `json:"data"`
	Args    map[string]string `json:"args"`
	Form    map[string]string `json:"form"`
	Headers map[string]string `json:"headers"`
	Cookies map[string]string `json:"cookies"`
	Url     string            `json:"url"`
	Origin  string            `json:"origin"`
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
	resp, err := Post(TestHost+"/post", body, "text/plain").Send()
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

func TestMixedCaseSchemeAcceptable(t *testing.T) {
	schemes := []string{"http://", "HTTP://", "hTtp://", "HttP://",
		"https://", "hTTps://", "HTTPS://", "Https://"}
	for i, v := range schemes {
		session := NewSession()
		resp, err := session.Get(v + "httpbin.org/get").Send()
		if err != nil {
			t.Fatalf("#%d: Error{%v}", i, err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("#%d: Failed for scheme {%s}", i, v)
		}
	}
}

func TestAllowRedirectGet(t *testing.T) {
	session := NewSession()
	resp, err := session.Get(TestHost + "/redirect/1").Send()
	if err != nil {
		t.Fatalf("Error{%v}", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Status: %d, want: 200", resp.StatusCode)
	}
}

func TestUserAgent(t *testing.T) {
	session := NewSession()
	session.UserAgent = "Mozilla/5.0"
	resp, err := session.Get(TestHost + "/user-agent").Send()
	if err != nil {
		t.Fatalf("Error{%v}", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Status: %d, want: 200", resp.StatusCode)
	}
	var respJson map[string]string
	err = resp.Json(&respJson)
	if err != nil {
		t.Fatalf("Json parse failed: %v", err)
	}
	value, ok := respJson["user-agent"]
	if !ok {
		t.Fatalf("Missing 'user-agent', data: %v", respJson)
	}
	if value != session.UserAgent {
		t.Fatalf("User-Agent: %s, want: %s", value, session.UserAgent)
	}
}

func TestCookies(t *testing.T) {
	session := NewSession()
	resp, err := session.Get(TestHost + "/cookies/set?foo=bar").Send()
	if err != nil {
		t.Fatalf("Error{%v}", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Status: %d, want: 200", resp.StatusCode)
	}
	cookie := session.Cookie("foo", nil)
	if cookie == nil {
		t.Fatalf("Cookie 'foo' not exist! cookies: %v", resp.Cookies())
	}
	if cookie.Value != "bar" {
		t.Fatalf("cookies['foo'] => %s, want 'bar'", cookie)
	}

	// Test cookie contained in next request:
	var body response
	err = resp.Json(&body)
	if err != nil {
		t.Fatalf("Json parse failed: %v", err)
	}
	v, ok := body.Cookies["foo"]
	if !ok {
		t.Fatal("Missing cookie 'foo' in Request.")
	}
	if v != "bar" {
		t.Fatalf("cookie['foo'] => %s, want: 'bar'", v)
	}
}
