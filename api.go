package requests

import (
	"net/url"
)

var DefaultUserAgent = "GoRequests"

func Request(method, url string) *HttpRequest {
	var r = NewRequest(method, url)
	r.SetHeader("User-Agent", DefaultUserAgent)
	r.SetHeader("Accept", "*/*")
	return r
}

func Get(url string) *HttpRequest {
	return Request("GET", url)
}

func Post(url string, data interface{}, bodyType string) *HttpRequest {
	return Request("POST", url).SetBody(data, bodyType)
}

func PostForm(u string, data url.Values) *HttpRequest {
	return Request("POST", u).SetBody(data.Encode(), "application/x-www-form-urlencoded")
}

func Put(url string, data interface{}, bodyType string) *HttpRequest {
	return Request("PUT", url).SetBody(data, bodyType)
}

func Head(url string) *HttpRequest {
	return Request("HEAD", url)
}

func Options(url string) *HttpRequest {
	return Request("OPTIONS", url)
}

func Patch(url string, data interface{}) *HttpRequest {
	return Request("PATCH", url)
}

func Delete(url string) *HttpRequest {
	return Request("DELETE", url)
}
