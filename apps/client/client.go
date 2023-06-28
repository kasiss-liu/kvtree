package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DataItem struct {
	Type  string      `json:"type"`
	Ver   string      `json:"ver"`
	Value interface{} `json:"value"`
}

func (d *DataItem) String() string {
	if d.Type == "string" {
		return d.Value.(string)
	}
	return ""
}

type HttpClient struct {
	Addr      string
	Auth      string
	Namespace string
	client    *http.Client
}

func NewBmccHttpClient(addr, auth, namespace string) *HttpClient {
	hc := &HttpClient{
		Addr:      addr,
		Auth:      auth,
		Namespace: namespace,
	}
	hc.client = &http.Client{
		Timeout: 10 * time.Second,
	}

	return hc
}

func (hc *HttpClient) Get(key string, ver ...string) (*DataItem, error) {
	version := ""
	if len(ver) > 0 {
		version = ver[0]
	}

	url := fmt.Sprintf("%s/auth/get?namespace=%s&key=%s&ver=%s", hc.Addr, hc.Namespace, key, version)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", hc.Auth)
	if err != nil {
		return nil, err
	}
	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 401:
		return nil, fmt.Errorf("UnAuthorization code 401")
	case 500:
		return nil, fmt.Errorf("erver internal error 500")
	case 404:
		return nil, fmt.Errorf("get key %s error: %d", key, 404)
	}
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	d := &DataItem{}
	err = json.Unmarshal(bs, &d)
	if err != nil {
		return nil, fmt.Errorf("get key %s error: %w", key, err)
	}
	return d, nil
}

func (hc *HttpClient) GetValString(key string, ver ...string) (string, error) {
	d, err := hc.Get(key, ver...)
	if err != nil {
		return "", err
	}
	return d.String(), nil
}

func (hc *HttpClient) GetValBytes(key string, ver ...string) ([]byte, error) {
	d, err := hc.Get(key, ver...)
	if err != nil {
		return nil, err
	}
	return []byte(d.String()), nil
}
