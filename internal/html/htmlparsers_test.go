package html_test

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/adepte-myao/test_parser/internal/html"
	"github.com/adepte-myao/test_parser/internal/tools"
)

var (
	sitemapParser *html.SitemapParser
	basePage      string
	sectionPage   string
	certAreaPage  string
	testPage      string
)

func TestMain(m *testing.M) {
	sitemapParser = html.NewSitemapParser("https://tests24.ru/")

	resp, err := tools.DoProperRequest(http.MethodGet, "https://tests24.ru/")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	basePage = string(bodyBytes)

	resp, err = tools.DoProperRequest(http.MethodGet, "https://tests24.ru/?iter=1&s_group=1")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	sectionPage = string(bodyBytes)

	resp, err = tools.DoProperRequest(http.MethodGet, "https://tests24.ru/?iter=2&group=4")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	certAreaPage = string(bodyBytes)

	resp, err = tools.DoProperRequest(http.MethodGet, "https://tests24.ru/?iter=3&test=726")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	testPage = string(bodyBytes)

	os.Exit(m.Run())
}
