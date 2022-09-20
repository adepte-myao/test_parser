package tools_test

import (
	"testing"

	"github.com/adepte-myao/test_parser/internal/tools"
)

func TestNewRequest(t *testing.T) {
	req, err := tools.NewRequest("GET", "http://someuri.com/what")
	if err != nil {
		t.Fatal(err)
	}

	expectedUserAgent := ""
	// expectedUserAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36"
	if req.Header.Get("User-Agent") != expectedUserAgent {
		t.Fatal("user agent is not valid")
	}

	expectedHost := "https://tests24.ru/"
	if req.Host != expectedHost {
		t.Fatal("host is not valid")
	}
}
