package storage_test

import (
	"fmt"
	"net/url"
	"testing"
)

func TestURL(t *testing.T) {
	host, _ := url.ParseRequestURI("http://localhost")

	query := fmt.Sprintf("appId=%s&token=%s", "1234", "sdkjfhasdkljfhaskdjhfkjls")

	u := &url.URL{
		Scheme:   host.Scheme,
		Host:     host.Host,
		Path:     "/web/password/reset",
		RawQuery: query,
	}

	uu := &url.URL{Scheme: u.Scheme, Host: u.Host, Path: u.Path}
	fmt.Println(uu.String())
	fmt.Println(u.String())
}
