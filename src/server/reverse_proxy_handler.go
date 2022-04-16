package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
			return a + b[1:]
	case !aslash && !bslash:
			return a + "/" + b
	}
	return a + b
}

// 구조화된 대상 URL (API 서버) 가 입력 파라미터로 들어오며, 해당 URL에 대해 reverse proxy를 리턴.
// 대상이 되는 API 서버의 기본 경로가 "/reverse" 이고 들어오는 요청이 "/api" 라면, 결과적으로 API 서버에 대한 요청은 "/reverse/api" 가 됨
func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
			if targetQuery == "" || req.URL.RawQuery == "" {
					req.URL.RawQuery = targetQuery + req.URL.RawQuery
			} else {
					req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
			}
			if _, ok := req.Header["User-Agent"]; !ok {
					req.Header.Set("User-Agent", "")
			}
	}

	modifyResp := func(resp *http.Response) error{
			var oldData,newData []byte
			oldData,err := ioutil.ReadAll(resp.Body)
			if err != nil{
					return err
			}
			if resp.StatusCode == 200 {
					newData = []byte("[INFO] " + string(oldData))

			}else{
					newData = []byte("[ERROR] " + string(oldData))
			}

			resp.Body = ioutil.NopCloser(bytes.NewBuffer(newData))
			resp.ContentLength = int64(len(newData))
			resp.Header.Set("Content-Length",fmt.Sprint(len(newData)))
			return nil
	}
	return &httputil.ReverseProxy{Director: director,ModifyResponse:modifyResp}
}

func ReverseProxyHandler(paramPort string) {
	apiServer := "http://localhost:8989/api" // 실제 API 서버 URL 주소. API 소스상에서 정의된 baseUrl이 반드시 포함되어야 함.
	targetUrl, err := url.Parse(apiServer) // 구조화된(Parse) API 서버 URL 주소를  targetUrl 변수에 load
	if err != nil {
		log.Fatal(err)
	}

	proxy := NewSingleHostReverseProxy(targetUrl)
	log.Println("Reverse proxy server serves at : http://localhost:" + paramPort)
	var portString string  = ":" + paramPort
	if err := http.ListenAndServe(portString, proxy); err != nil {
		log.Fatal("Start server failed, err : ", err)
	}
}
