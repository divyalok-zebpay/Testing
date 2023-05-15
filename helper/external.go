package external

import (
	"brave/helper/httpclient"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"
)

type (
	HTTPMethod string
	Params     map[string]string
)

const (
	HttpMethodGet    = HTTPMethod(http.MethodGet)
	HttpMethodPost   = HTTPMethod(http.MethodPost)
	HTTPMethodDelete = HTTPMethod(http.MethodDelete)
)

const DefaultTimeOut = 60 * time.Second

type HTTPCallParams struct {
	Client  httpclient.HTTPClient
	Method  HTTPMethod
	URL     string
	Payload []byte
	Headers map[string]string
	Params  map[string]interface{}
}

func HTTPCall(params *HTTPCallParams) (int, []byte, error) {
	req, err := http.NewRequest(string(params.Method), params.URL, bytes.NewBuffer(params.Payload))
	if err != nil {
		return 0, []byte{}, err
	}
	if len(params.Params) > 0 {
		q := req.URL.Query()
		for k, v := range params.Params {
			q.Add(k, v.(string))
		}
		req.URL.RawQuery = q.Encode()
	}
	if len(params.Headers) > 0 {
		for k, v := range params.Headers {
			req.Header.Add(k, v)
		}
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := params.Client.Do(req)
	if err != nil {
		return http.StatusBadGateway, []byte{}, err
	}
	if resp == nil {
		return 500, []byte{}, err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, []byte{}, err
	}
	return resp.StatusCode, bodyBytes, nil
}

func ExternalAPICall(httpMethod HTTPMethod, url string, payload io.Reader, headers http.Header, params Params) (*http.Response, error) {
	client := &http.Client{
		Timeout: DefaultTimeOut,
	}

	req, err := http.NewRequest(string(httpMethod), url, payload)
	if err != nil {
		return nil, err
	}

	for i := range headers {
		for j := range headers[i] {
			req.Header.Add(i, headers[i][j])
		}
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "znew.user "+runtime.Version())

	if httpMethod == HttpMethodGet {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		// assign encoded query string to http request
		req.URL.RawQuery = q.Encode()
	}

	return client.Do(req)
}
