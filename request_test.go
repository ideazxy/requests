package requests

import (
	"bytes"
	"testing"
)

func TestHeader(t *testing.T) {
	var buf bytes.Buffer
	req := NewRequest("GET", "http://a.b", nil)
	req.Header("Content-Type", "image/png")
	req.Req.Header.Write(&buf)
	want := "Content-Type: image/png\r\n"
	if buf.String() != want {
		t.Fatalf("#1 : %s, want: %s", buf.String(), want)
	}
	buf.Reset()

	req.Header("Content-Type", "image/jpeg")
	req.Req.Header.Write(&buf)
	want = "Content-Type: image/jpeg\r\n"
	if buf.String() != want {
		t.Fatalf("#2 : %s, want: %s", buf.String(), want)
	}
	buf.Reset()

	req.Header("key", "value")
	req.Req.Header.Write(&buf)
	want = "Content-Type: image/jpeg\r\nKey: value\r\n"
	if buf.String() != want {
		t.Fatalf("#3 : %s, want: %s.", buf.String(), want)
	}
	buf.Reset()
}

func TestEncodeUrl(t *testing.T) {
	req := NewRequest("GET", "http://a.b", nil)
	req.Param("k1", "v1")
	err := req.encodeUrl()
	if err != nil {
		t.Fatalf("#1: %v", err)
	}
	queries := req.Req.URL.Query()
	v, ok := queries["k1"]
	if !ok {
		t.Fatal("#1: k1 not found")
	}
	if v[0] != "v1" {
		t.Fatal("#1: k1 => %s, want v1", v)
	}

	req = NewRequest("GET", "http://a.b?k0=v0", nil)
	req.Param("k1", "v1")
	err = req.encodeUrl()
	if err != nil {
		t.Fatalf("#2: %v", err)
	}
	queries = req.Req.URL.Query()
	v, ok = queries["k1"]
	if !ok {
		t.Fatal("#2: k1 not found")
	}
	if v[0] != "v1" {
		t.Fatal("#2: k1 => %s, want v1", v)
	}
	v, ok = queries["k0"]
	if !ok {
		t.Fatal("#2: k0 not found")
	}
	if v[0] != "v0" {
		t.Fatal("#2: k0 => %s, want v0", v)
	}

	req = NewRequest("GET", "http://a.b?k0=v0", nil)
	req.Param("k1", "vv1")
	err = req.encodeUrl()
	if err != nil {
		t.Fatalf("#3: %v", err)
	}
	queries = req.Req.URL.Query()
	v, ok = queries["k1"]
	if !ok {
		t.Fatal("#3: k1 not found")
	}
	if v[0] != "vv1" {
		t.Fatal("#3: k1 => %s, want vv1", v)
	}
}
