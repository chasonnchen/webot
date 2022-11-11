package wkteam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	httpClient = &HttpClient{}
)

type HttpClient struct {
}

func NewHttpClient() *HttpClient {
	return httpClient
}

type PostOptions struct {
	Header map[string]string
	//TODO cookie
}

func (h *HttpClient) PostJson(url string, data interface{}, options PostOptions) (body []byte, err error) {
	log.Printf("%s post json resquest is %+v", url, data)
	jsonStr, err := json.Marshal(data)
	if err != nil {
		log.Printf("post json request data json encode err %+v", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	// post json，所以固定加这个header.
	req.Header.Set("Content-Type", "application/json")

	if options.Header != nil {
		for field, value := range options.Header {
			if strings.ToLower(field) == "host" {
				req.Host = value
			}
			req.Header.Set(field, value)
		}
	}
	// TODO cookie

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       100,
			IdleConnTimeout:    10 * time.Second,
			DisableCompression: true,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http code not 200")
	}
	log.Printf("%s post json response is %s", url, string(body))

	return body, nil
}
