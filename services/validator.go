package services

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

func IsValidURL(raw string) (bool, error) {

	if len(raw) > 2048 {
		return false, fmt.Errorf("url is too long")
	}

	u, err := url.ParseRequestURI(raw)

	if err != nil {
		return false, fmt.Errorf("url is null: %w", err)
	}

	if u.Scheme == "" || u.Host == "" {
		return false, fmt.Errorf("url not has scheme or host")
	}

	if !isHttpURL(u) {
		return false, fmt.Errorf("url is not http or https")
	}

	if !isPublicIP(u) {
		return false, fmt.Errorf("url is not public IP")
	}

	if !isAvailable(u) {
		return false, fmt.Errorf("url is not available: %w", err)
	}

	return true, nil
}

func isHttpURL(u *url.URL) bool {
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	return true
}

func isPublicIP(u *url.URL) bool {

	ips, err := net.LookupIP(u.Hostname())
	if err != nil || len(ips) == 0 {
		return false
	}

	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() {
			return false
		}
	}
	return true
}

func isAvailable(u *url.URL) bool {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Head(u.String())
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 400
}
