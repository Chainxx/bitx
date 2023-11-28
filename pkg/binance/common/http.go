package common

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func client(timeout int) *http.Client {
	return &http.Client{
		//Transport: tr,
		Timeout: time.Duration(timeout) * time.Second,
	}
}

func request(method, url, body string, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	
	return req, nil
}

func Do(method, url, body string, headers map[string]string) ([]byte, int, error) {
	cli := client(5)
	req, err := request(method, url, body, headers)
	
	if err != nil {
		return nil, 0, err
	}
	
	resp, err := cli.Do(req)
	
	if err != nil {
		return nil, 0, err
	}
	
	statusCode := resp.StatusCode
	respBody := resp.Body
	defer resp.Body.Close()
	
	respByte, err := ioutil.ReadAll(respBody)
	if err != nil {
		return nil, 0, err
	}
	
	if strings.HasPrefix(fmt.Sprint(statusCode), "4") || strings.HasPrefix(fmt.Sprint(statusCode), "5") {
		return nil, statusCode, errors.New(string(respByte))
	}
	
	return respByte, statusCode, nil
}
