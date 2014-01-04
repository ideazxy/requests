package requests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"testing"
)

func TestEncodeUrl(t *testing.T) {
	req := NewRequest("GET", "http://a.b?k=v")
	expected := url.Values{"k": {"v"}}
	if fmt.Sprintf("%v", req.Params) != fmt.Sprintf("%v", expected) {
		t.Fatalf("#0: Params => %v, expected: %v", req.Params, expected)
	}

	req.AddParam("k1", "v1")
	err := req.Prepare()
	if err != nil {
		t.Fatalf("#1: %v", err)
	}
	queries := req.Query()
	v, ok := queries["k1"]
	if !ok {
		t.Fatal("#1: k1 not found")
	}
	if fmt.Sprintf("%v", v) != "[v1]" {
		t.Fatal("#1: Req.URL.Query()[k1] => %v, want: [v1]", v)
	}

	req.AddParam("k1", "v11")
	err = req.Prepare()
	if err != nil {
		t.Fatalf("#2: %v", err)
	}
	queries = req.Query()
	v, ok = queries["k1"]
	if !ok {
		t.Fatal("#2: k1 not found")
	}
	if fmt.Sprintf("%v", v) != "[v1 v11]" {
		t.Fatalf("#2: Req.URL.Query()[k1] => %v, want: [v1, v11]", v)
	}

	req.SetParam("k1", "vv1")
	err = req.Prepare()
	if err != nil {
		t.Fatalf("#3: %v", err)
	}
	queries = req.Query()
	v, ok = queries["k1"]
	if !ok {
		t.Fatal("#3: k1 not found")
	}
	if fmt.Sprintf("%v", v) != "[vv1]" {
		t.Fatalf("#3: Req.URL.Query()[k1] => %v, want [vv1]", v)
	}

	req.DelParam("k1")
	err = req.Prepare()
	if err != nil {
		t.Fatalf("#4: %v", err)
	}
	queries = req.Query()
	v, ok = queries["k1"]
	if ok {
		t.Fatal("#4: k1 should not exist!")
	}

	req.DelParam("k")
	err = req.Prepare()
	if err != nil {
		t.Fatalf("#5: %v", err)
	}
	queries = req.Query()
	v, ok = queries["k"]
	if ok {
		t.Fatal("#5: k should not exist!")
	}

	req = NewRequest("HEAD", "http://example.com/path#fragment")
	req.AddParam("Key", "Value")
	err = req.Prepare()
	if err != nil {
		t.Fatalf("#6: %v", err)
	}
	expectedPath := "http://example.com/path?Key=Value#fragment"
	if req.Url() != expectedPath {
		t.Fatalf("#6: path() => %s, want: %s", req.Url(), expectedPath)
	}

	req = NewRequest("HEAD", "http://example.com/path?A=B#fragment")
	req.AddParam("Key", "Value")
	err = req.Prepare()
	if err != nil {
		t.Fatalf("#7: %v", err)
	}
	expectedPath = "http://example.com/path?A=B&Key=Value#fragment"
	if req.Url() != expectedPath {
		t.Fatalf("#7: path() => %s, want: %s", req.Url(), expectedPath)
	}
}

func TestBody(t *testing.T) {
	content := "string body."
	r := NewRequest("Get", "http://httpbin.org")

	r.SetBody(content, "text/plain")
	body, err := ioutil.ReadAll(r.Req.Body)
	if err != nil {
		t.Fatal("#1: Body => %v", err)
	}
	if string(body) != content {
		t.Errorf("#1: Body => %s, want %s.", string(body), content)
	}
	if r.Req.ContentLength != int64(len(content)) {
		t.Errorf("#1: ContentLength => %d, want %d.", r.Req.ContentLength, len(content))
	}

	content = "[]byte body"
	r.SetBody([]byte(content), "text/plain")
	body, err = ioutil.ReadAll(r.Req.Body)
	if err != nil {
		t.Fatal("#2: Body() => %v", err)
	}
	if string(body) != content {
		t.Errorf("#2: Body() => %s, want %s.", string(body), content)
	}
	if r.Req.ContentLength != int64(len(content)) {
		t.Errorf("#2: ContentLength => %d, want %d.", r.Req.ContentLength, len(content))
	}

	content = "buffer body"
	var f bytes.Buffer
	fmt.Fprint(&f, content)
	r.SetBody(&f, "text/plain")
	body, err = ioutil.ReadAll(r.Req.Body)
	if err != nil {
		t.Fatal("#3: Body() => %v", err)
	}
	if string(body) != content {
		t.Errorf("#3: Body() => %s, want %s.", string(body), content)
	}
	if r.Req.ContentLength != int64(len(content)) {
		t.Errorf("#3: ContentLength => %d, want %d.", r.Req.ContentLength, len(content))
	}
}
