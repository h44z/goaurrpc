package aur

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

var aurTimeout = 30 * time.Second
var client *http.Client

func init() {
	dialer := &net.Dialer{
		Timeout: aurTimeout,
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}

		ips, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
		if err != nil {
			return nil, err
		}

		// Prefer IPv6
		for _, ip := range ips {
			if ip.To16() != nil && ip.To4() == nil {
				conn, err := dialer.DialContext(ctx, "tcp6", net.JoinHostPort(ip.String(), port))
				if err == nil {
					return conn, nil
				}
			}
		}

		// Fallback to IPv4
		for _, ip := range ips {
			if ip.To4() != nil {
				return dialer.DialContext(ctx, "tcp4", net.JoinHostPort(ip.String(), port))
			}
		}

		return nil, fmt.Errorf("no valid IPs for %s", host)
	}

	client = &http.Client{Transport: transport, Timeout: aurTimeout}
}

// DownloadPackageData downloads package data file from AUR; decompression happens automatically
func DownloadPackageData(address string, lastmod time.Time) ([]byte, time.Time, error) {
	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		return nil, lastmod, err
	}
	req.Header.Set("If-Modified-Since", lastmod.Format(http.TimeFormat))

	r, err := client.Do(req)
	if err != nil {
		return nil, lastmod, err
	}
	defer r.Body.Close()

	if r.StatusCode == http.StatusNotModified {
		_, _ = io.Copy(io.Discard, r.Body) // consume body
		return nil, lastmod, errors.New("not modified")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, lastmod, err
	}

	newmod, err := http.ParseTime(r.Header.Get("Last-Modified"))
	if err != nil {
		newmod = time.Now()
	}

	return body, newmod, nil
}
