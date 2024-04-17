package main

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type Backend struct {
	URL *url.URL
}

type BackendList []*Backend

var backendClient = &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:    3 * time.Second,
		MaxIdleConns:           0,
		MaxIdleConnsPerHost:    4,
		MaxConnsPerHost:        0,
		IdleConnTimeout:        90 * time.Second,
		ResponseHeaderTimeout:  2 * time.Second,
		MaxResponseHeaderBytes: 2 * 1024,
		ForceAttemptHTTP2:      true,
	},
}

var (
	ErrNotAllowed = errors.New("not allowed")
)

func (b *Backend) String() string {
	return b.URL.String()
}

func (b *Backend) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

func (b *Backend) UnmarshalJSON(data []byte) error {
	var str string

	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	parsed, err := url.Parse(str)
	if err != nil {
		return err
	}
	b.URL = parsed

	return nil
}

func (b *Backend) Do(ctx context.Context, r *http.Request) error {
	req, _ := http.NewRequestWithContext(ctx, "GET", "", nil)
	req.URL = &url.URL{
		Scheme:   b.URL.Scheme,
		Opaque:   b.URL.Opaque,
		User:     b.URL.User,
		Host:     b.URL.Host,
		Path:     b.URL.Path,
		RawPath:  b.URL.RawPath,
		RawQuery: r.URL.RawQuery,
	}

	// Copy headers
	headers := []string{
		"X-Session-Id",
		"X-Real-Ip",
		"X-Real-Path",
		"X-Real-Origin",
		"X-Real-Ua",
	}

	for _, header := range headers {
		if v := r.Header.Get(header); v != "" {
			req.Header.Set(header, v)
		}
	}

	response, err := backendClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return ErrNotAllowed
	}

	return nil
}

func (bl BackendList) Check(r *http.Request) bool {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var (
		allow atomic.Bool
		wg    sync.WaitGroup
	)

	for _, backend := range bl {
		wg.Add(1)

		go func(ctx context.Context, b *Backend) {
			defer wg.Done()

			if err := b.Do(ctx, r); err == nil {
				allow.Store(true)
				cancel()
			}
		}(ctx, backend)
	}

	wg.Wait()

	return allow.Load()
}
