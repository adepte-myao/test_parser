package tools

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

func NewTestsClient() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: 5 * time.Second,
	}

	return client, nil
}

func GetCookies(client *http.Client, uri string) []*http.Cookie {
	parserURI, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}
	cookies := client.Jar.Cookies(parserURI)
	return cookies
}
