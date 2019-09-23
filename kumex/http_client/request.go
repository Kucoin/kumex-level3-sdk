package http_client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Client struct {
	baseUrl       string
	skipVerifyTls bool
	timeout       time.Duration
	signer        *KcSigner
}

func NewClient(baseUrl string, skipVerifyTls bool, timeout time.Duration) *Client {
	var signer *KcSigner
	if os.Getenv("API_KEY") != "" {
		signer = NewKcSigner(os.Getenv("API_KEY"), os.Getenv("API_SECRET"), os.Getenv("API_PASSPHRASE"))
	}
	return &Client{
		baseUrl:       strings.TrimRight(baseUrl, "/"),
		skipVerifyTls: skipVerifyTls,
		timeout:       timeout,
		signer:        signer,
	}
}

func (c *Client) Request(method string, uri string, params map[string]string) (*Response, error) {
	query := make(url.Values)
	var body []byte
	switch method {
	case http.MethodGet, http.MethodDelete:
		for key, value := range params {
			query.Add(key, value)
		}
	default:
		if params == nil {
			break
		}
		b, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		body = b
	}

	if len(query) > 0 {
		if strings.Contains(uri, "?") {
			uri += "&" + query.Encode()
		} else {
			uri += "?" + query.Encode()
		}
	}

	req, err := http.NewRequest(method, c.baseUrl+uri, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "KuMEX-Go-SDK/1.0")

	//sign
	if c.signer != nil {
		var b bytes.Buffer
		b.WriteString(method)
		b.WriteString(uri)
		b.Write(body)
		h := c.signer.Headers(b.String())
		//log.Println(base.ToJsonString(h))
		for k, v := range h {
			//log.Printf("%s : %s", k, v)
			req.Header.Set(k, v)
		}
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: c.skipVerifyTls},
		},
		Timeout: c.timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return NewResponse(req, resp), nil
}
