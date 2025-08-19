package services

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func IsValidURL(raw string) (*url.URL, error) {
	// Validación básica de longitud
	if len(raw) > 2048 {
		return nil, fmt.Errorf("url is too long (max 2048 characters)")
	}

	if len(raw) < 10 {
		return nil, fmt.Errorf("url is too short")
	}

	if containsDangerousChars(raw) {
		return nil, fmt.Errorf("url contains dangerous characters")
	}

	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid url format: %w", err)
	}

	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("url must have scheme and host")
	}

	if !isHttpURL(u) {
		return nil, fmt.Errorf("only http and https schemes are allowed")
	}

	if !isValidHost(u.Host) {
		return nil, fmt.Errorf("invalid host format")
	}

	if !isPublicIP(u) {
		return nil, fmt.Errorf("private/local urls are not allowed")
	}

	if !isAvailable(u) {
		return nil, fmt.Errorf("url is not reachable")
	}

	return u, nil
}

func containsDangerousChars(s string) bool {

	dangerousPattern := `[\x00-\x1F\x7F-\x9F<>"'&]`
	matched, _ := regexp.MatchString(dangerousPattern, s)
	return matched
}

func isValidHost(host string) bool {

	hostname := host
	if strings.Contains(host, ":") {
		hostname = host[:strings.LastIndex(host, ":")]
	}

	if len(hostname) > 253 {
		return false
	}
	validPattern := `^[a-zA-Z0-9.-]+$`
	matched, _ := regexp.MatchString(validPattern, hostname)
	return matched
}

func isHttpURL(u *url.URL) bool {
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	return true
}

func isPublicIP(u *url.URL) bool {
	hostname := u.Hostname()

	// Lista de hosts bloqueados
	blockedHosts := []string{
		"localhost", "127.0.0.1", "0.0.0.0", "::1",
		"169.254.", "224.", "239.", "255.255.255.255",
		"metadata.google.internal", "instance-data",
	}

	for _, blocked := range blockedHosts {
		if strings.Contains(strings.ToLower(hostname), blocked) {
			return false
		}
	}

	ips, err := net.LookupIP(hostname)
	if err != nil || len(ips) == 0 {
		return false
	}

	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() || ip.IsMulticast() {
			return false
		}

		// Verificar rangos adicionales peligrosos
		if ip.To4() != nil {
			ipv4 := ip.To4()
			// Bloquear 169.254.x.x (link-local)
			if ipv4[0] == 169 && ipv4[1] == 254 {
				return false
			}
			// Bloquear multicast (224.0.0.0 - 239.255.255.255)
			if ipv4[0] >= 224 && ipv4[0] <= 239 {
				return false
			}
		}
	}
	return true
}

func isAvailable(u *url.URL) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {

			if len(via) >= 3 {
				return fmt.Errorf("too many redirects")
			}

			if !isPublicIP(&url.URL{Host: req.URL.Host}) {
				return fmt.Errorf("redirect to private IP not allowed")
			}
			return nil
		},
	}

	req, err := http.NewRequest("HEAD", u.String(), nil)
	if err != nil {
		return false
	}

	req.Header.Set("User-Agent", "URL-Shortener-Bot/1.0")
	req.Header.Set("Accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 400
}
