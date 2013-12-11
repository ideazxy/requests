package requests

import ()

var DefaultUserAgent = "GoRequests"

func Request(method, url string, data interface{}) *HttpRequest {
	var r = NewRequest(method, url, data)
	r.Header("User-Agent", DefaultUserAgent)
	r.Header("Accept", "*/*")
	return r
}

func Get(url string) *HttpRequest {
	return Request("GET", url, nil)
}

func Post(url string, data interface{}) *HttpRequest {
	return Request("POST", url, data)
}

func Put(url string, data interface{}) *HttpRequest {
	return Request("PUT", url, data)
}

func Head(url string) *HttpRequest {
	return Request("HEAD", url, nil)
}

func Options(url string) *HttpRequest {
	return Request("OPTIONS", url, nil)
}

func Patch(url string, data interface{}) *HttpRequest {
	return Request("PATCH", url, data)
}

func Delete(url string) *HttpRequest {
	return Request("DELETE", url, nil)
}
