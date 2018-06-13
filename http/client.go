// Copyright (c) 2016-2018, Jan Cajthaml <jan.cajthaml@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/jancajthaml/rest-contract-test/model"
)

type dialContextFn func(ctx context.Context, network, address string) (net.Conn, error)

// DialContext implements our own dialer in order to set read and write idle timeouts.
func DialContext(rwtimeout, ctimeout time.Duration) dialContextFn {
	dialer := &net.Dialer{Timeout: ctimeout}
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		c, err := dialer.DialContext(ctx, network, addr)
		if err != nil {
			return nil, err
		}

		if rwtimeout > 0 {
			timeoutConn := &tcpConn{
				TCPConn: c.(*net.TCPConn),
				timeout: rwtimeout,
			}
			return timeoutConn, nil
		}

		return c, nil
	}
}

type tcpConn struct {
	*net.TCPConn
	timeout time.Duration
}

func (c *tcpConn) Read(b []byte) (int, error) {
	err := c.TCPConn.SetDeadline(time.Now().Add(c.timeout))
	if err != nil {
		return 0, err
	}
	return c.TCPConn.Read(b)
}

func (c *tcpConn) Write(b []byte) (int, error) {
	err := c.TCPConn.SetDeadline(time.Now().Add(c.timeout))
	if err != nil {
		return 0, err
	}
	return c.TCPConn.Write(b)
}

func (client *HttpClient) Call(endpoint *model.Endpoint) error {
	if endpoint == nil {
		return fmt.Errorf("no endpoint provided")
	}

	// fixme add defer recover error, don't panic here

	switch endpoint.Method {
	case "GET":
		_, code, err := client.Get(endpoint.URI, endpoint.Request.Headers)
		if err != nil {
			return err
		}
		fmt.Println("OK    |", code, *endpoint)
	default:
		fmt.Println("SKIP  |", *endpoint)
	}

	return nil
}

type HttpClient struct {
	client *http.Client
}

func NewHttpClient() *HttpClient {
	cookieJar, _ := cookiejar.New(nil)

	transport := &http.Transport{
		DialContext:           DialContext(500*time.Millisecond, 500*time.Millisecond),
		IdleConnTimeout:       100 * time.Millisecond,
		TLSHandshakeTimeout:   300 * time.Millisecond,
		ExpectContinueTimeout: 100 * time.Millisecond,
		ResponseHeaderTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		Proxy:                 http.ProxyFromEnvironment,
	}

	client := &http.Client{
		Jar:       cookieJar,
		Transport: transport,
	}

	return &HttpClient{
		client: client,
	}
}

func (c *HttpClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *HttpClient) Get(url string, headers map[string]string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, -1, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return contents, resp.StatusCode, nil
}
