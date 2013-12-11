package requests

import (
	"bufio"
	"net/http"
	"strings"
	"testing"
)

var newResponseTests = []struct {
	rawResp string
	status  int
}{
	{`HTTP/1.0 200 OK
Connection: close

Body here`,
		200},
	{`HTTP/1.1 200 OK

Body here`,
		200},
	{`HTTP/1.1 204 No Content

Body should not be read!`,
		204},
	{`HTTP/1.0 200 OK
Content-Length: 10
Connection: close

Body here`,
		200},
	{"HTTP/1.1 200 OK\r\n" +
		"Transfer-Encoding: chunked\r\n" +
		"\r\n" +
		"0a\r\n" +
		"Body here\n\r\n" +
		"09\r\n" +
		"continued\r\n" +
		"0\r\n" +
		"\r\n",
		200},
	{"HTTP/1.1 200 OK\r\n" +
		"Transfer-Encoding: chunked\r\n" +
		"Content-Length: 10\r\n" +
		"\r\n" +
		"0a\r\n" +
		"Body here\n\r\n" +
		"0\r\n" +
		"\r\n",
		200},
	{"HTTP/1.0 303 \r\n\r\n",
		303},
}

func TestNewResponse(t *testing.T) {
	for i, v := range newResponseTests {
		resp, err := http.ReadResponse(bufio.NewReader(strings.NewReader(v.rawResp)), &http.Request{Method: "GET"})
		if err != nil {
			t.Errorf("#%d: %v", i, err)
			continue
		}
		r := NewResponse(resp)
		if r.Status != v.status {
			t.Errorf("#%d.NewResponse(%s): Status = %d, want %d.", i, v.rawResp, r.Status, v.status)
		}
	}
}

var contentTests = []struct {
	rawResp string
	content string
}{
	{`HTTP/1.1 200 OK

Body`,
		"Body"},
}

func TestContent(t *testing.T) {
	for i, v := range contentTests {
		resp, err := http.ReadResponse(bufio.NewReader(strings.NewReader(v.rawResp)), &http.Request{Method: "GET"})
		if err != nil {
			t.Errorf("#%d: %v", i, err)
			continue
		}
		r := NewResponse(resp)
		body, err := r.Content()
		if err != nil {
			t.Errorf("#%d: %v", i, err)
			continue
		}
		if string(body) != v.content {
			t.Errorf("#%d.Content() => %v, want %v.", i, body, []byte(v.content))
		}
		// Call Content() again:
		body, err = r.Content()
		if err != nil {
			t.Errorf("#%d: %v", i, err)
			continue
		}
		if string(body) != v.content {
			t.Errorf("#%d.Content() again => %v, want %v.", i, body, []byte(v.content))
		}
	}
}

func TestText(t *testing.T) {
	for i, v := range contentTests {
		resp, err := http.ReadResponse(bufio.NewReader(strings.NewReader(v.rawResp)), &http.Request{Method: "GET"})
		if err != nil {
			t.Errorf("#%d: %v", i, err)
			continue
		}
		r := NewResponse(resp)
		body, err := r.Text()
		if err != nil {
			t.Errorf("#%d: %v", i, err)
			continue
		}
		if body != v.content {
			t.Errorf("#%d.Content() => %v, want %v.", i, body, v.content)
		}
		// Call Content() again:
		body, err = r.Text()
		if err != nil {
			t.Errorf("#%d: %v", i, err)
			continue
		}
		if body != v.content {
			t.Errorf("#%d.Content() again => %v, want %v.", i, body, v.content)
		}
	}
}
