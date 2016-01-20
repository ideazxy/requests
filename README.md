#Requests [![Build Status](https://travis-ci.org/ideazxy/requests.svg?branch=master)](https://travis-ci.org/ideazxy/requests)

Simple HTTP library, do the same work as [Requests](https://github.com/kennethreitz/requests), but for Go.

## OVERVIEW

To GET or POST or PUT or DELETE...

    import "github.com/ideazxy/requests"
    
    getResp, err := requests.Get("http://httpbin.org/get").Send()
    
    postResp, err := requests.Post("http://httpbin.org/post", "string or []byte or io.Reader or io.ReadCloser", "text/plain").Send()
    
    putResp, err := requests.Put("http://httpbin.org/put", []byte("the same as post api"), "application/json").Send()
    
    delResp, err := requests.Delete("http://httpbin.org/delete").Send()
    
Do more settings in pipline style:

    resp, err := requests.Get("http://httpbin.org/get").SetHeader("User-Agent", "Chrome").AddParam("key", "value").AllowRedirects(false).Send()
    
Fetch data from reponse:

	// Get string:
    s, err := resp.Text()
    
    // Get []byte
    b, err := resp.Content()
    
    // Get Json struct
    var j struct{Msg `json:"message"`}
    err := resp.Json(&j)
    
    // Get Xml struct
    var x struct{Msg `xml:"message"`}
    err := resp.Xml(&x)

## INSTALLATION

To install:

    go get github.com/ideazxy/requests
    
## LICENSE

requests is licensed under the Apache Licence, Version 2.0  
<http://www.apache.org/licenses/LICENSE-2.0.html>.