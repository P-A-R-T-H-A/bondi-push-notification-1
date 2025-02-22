package service

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/url"
	"time"
)

func httpGet(apiUrl string, headers map[string]string, paramQuery url.Values, timeout time.Duration) ([]byte, http.Header, string, int, error) {
	return baseHttpGet(apiUrl, "GET", headers, paramQuery, nil, "", "", timeout)
}

func httpGetWithAuth(apiUrl string, headers map[string]string, paramQuery url.Values, basicAuthUser string, basicAuthPass string, timeout time.Duration) ([]byte, http.Header, string, int, error) {
	return baseHttpGet(apiUrl, "GET", headers, paramQuery, nil, basicAuthUser, basicAuthPass, timeout)
}

func httpPost(apiUrl string, headers map[string]string, bodyData []byte, timeout time.Duration) ([]byte, http.Header, string, int, error) {
	return baseHttpGet(apiUrl, "POST", headers, nil, bodyData, "", "", timeout)
}

func httpPut(apiUrl string, headers map[string]string, bodyData []byte, timeout time.Duration) ([]byte, http.Header, string, int, error) {
	return baseHttpGet(apiUrl, "PUT", headers, nil, bodyData, "", "", timeout)
}

func httpDelete(apiUrl string, headers map[string]string, paramQuery url.Values, timeout time.Duration) ([]byte, http.Header, string, int, error) {
	return baseHttpGet(apiUrl, "DELETE", headers, paramQuery, nil, "", "", timeout)
}

func baseHttpGet(apiUrl string, method string, headers map[string]string, paramQuery url.Values, bodyData []byte,
	basicAuthUser string, basicAuthPass string, timeout time.Duration) ([]byte, http.Header, string, int, error) {

	var req *http.Request
	var err error
	if bodyData != nil {
		req, err = http.NewRequest(method, apiUrl, bytes.NewBuffer(bodyData))
	} else {
		req, err = http.NewRequest(method, apiUrl, nil)
	}
	if req != nil {
		req.Close = true
	}
	if err == nil {
		client := &http.Client{Timeout: time.Second * timeout}

		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept-Encoding", "gzip,deflate,sdch")

		if headers != nil && len(headers) > 0 {
			for key, value := range headers {
				req.Header.Set(key, value)
			}
		}

		if len(basicAuthUser) > 0 && len(basicAuthPass) > 0 {
			req.SetBasicAuth(basicAuthUser, basicAuthPass)
		}

		if paramQuery != nil && len(paramQuery) > 0 {
			req.URL.RawQuery = paramQuery.Encode()
		}

		resp, err := client.Do(req)
		var responseData []byte
		var responseHeaders http.Header
		responseStatus := 0

		if err == nil {
			responseStatus = resp.StatusCode
			responseHeaders = resp.Header

			if resp.Header.Get("Content-Encoding") == "gzip" {
				var reader *gzip.Reader
				reader, err = gzip.NewReader(resp.Body)
				if err == nil {
					responseData, err = io.ReadAll(reader)
				}
			} else {
				responseData, err = io.ReadAll(resp.Body)
			}

			defer func() {
				err = resp.Body.Close()
				if err != nil {
					return
				}
			}()
		}
		return responseData, responseHeaders, req.URL.String(), responseStatus, err
	} else {
		return nil, nil, apiUrl, 0, err
	}
}
