package httpClient

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

const METHOD_GET = "GET"
const METHOD_POST = "POST"
const METHOD_PUT = "PUT"
const METHOD_DELETE = "DELETE"

type HttpClient struct {
	url       string
	timeout   time.Duration
	multiplex bool
	proxy     string
	staticIP  string
}

func NewHttpClient(url string, timeout time.Duration, multiplex bool) *HttpClient {
	client := new(HttpClient)
	client.url = url
	client.timeout = timeout
	client.multiplex = multiplex
	return client
}

func (this *HttpClient) WithProxy(proxy string) *HttpClient {
	this.proxy = proxy
	return this
}
func (this *HttpClient) WithStaticIP(ip string) *HttpClient {
	this.staticIP = ip
	return this
}

func (this *HttpClient) createHttpClient() *http.Client {
	client := &http.Client{Timeout: this.timeout}

	transport := new(http.Transport)
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	if len(this.proxy) > 0 {
		transport.Proxy = func(_ *http.Request) (*url.URL, error) {
			return url.Parse(this.proxy)
		}
	}
	if len(this.staticIP) > 0 {
		transport.DialContext = func(ctx context.Context, netw, addr string) (net.Conn, error) {
			//本地地址  ipaddr是本地外网IP
			lAddr, err := net.ResolveTCPAddr(netw, this.staticIP+":0")
			if err != nil {
				return nil, err
			}
			//被请求的地址
			rAddr, err := net.ResolveTCPAddr(netw, addr)
			if err != nil {
				return nil, err
			}
			conn, err := net.DialTCP(netw, lAddr, rAddr)
			if err != nil {
				return nil, err
			}
			return conn, nil
		}
	}
	client.Transport = transport
	return client
}

func (this *HttpClient) Get() HttpResponse {
	client := this.createHttpClient()
	resp, err := client.Get(this.url)
	if err != nil {
		return NewErrorHttpResponse(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		_ = resp.Body.Close()
		return NewErrorHttpResponse(err)
	}

	_ = resp.Body.Close()
	return NewHttpResponse(err, resp.StatusCode, resp.Status, body, resp.Header)
}

func (this *HttpClient) PostJson(gson []byte) HttpResponse {
	header := make(map[string]string)
	header["Content-Type"] = "application/json; charset=utf-8"
	return this.PostJsonWithHeader(header, gson)
}

func (this *HttpClient) PostJsonWithHeader(header map[string]string, gson []byte) HttpResponse {
	client := this.createHttpClient()
	var reqs *http.Request
	var err error
	if len(gson) == 0 {
		reqs, err = http.NewRequest(METHOD_POST, this.url, nil)
	} else {
		reqs, err = http.NewRequest(METHOD_POST, this.url, bytes.NewReader(gson))
	}

	if err != nil {
		return NewErrorHttpResponse(err)
	}

	if this.multiplex {
		reqs.Close = false
	} else {
		reqs.Close = true
	}

	if header != nil {
		for key, val := range header {
			reqs.Header.Add(key, val)
		}
	}

	resp, err := client.Do(reqs)
	if err != nil {
		return NewErrorHttpResponse(err)
	}

	resBuff, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return NewHttpResponse(err, resp.StatusCode, resp.Status, resBuff, resp.Header)
}
