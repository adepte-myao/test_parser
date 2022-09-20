package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/adepte-myao/test_parser/internal/dto"
	"github.com/adepte-myao/test_parser/internal/tools"
	"github.com/sirupsen/logrus"
)

type LinksHandler struct {
	logger      *logrus.Logger
	testsClient *http.Client
}

func NewLinksHandler(logger *logrus.Logger) *LinksHandler {
	client, err := tools.NewTestsClient()
	if err != nil {
		panic(err)
	}

	return &LinksHandler{
		logger:      logger,
		testsClient: client,
	}
}

func (handler LinksHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Info("Find all links request received")

	var rd dto.LinksRequestData
	err := json.NewDecoder(r.Body).Decode(&rd)
	if err != nil {
		handler.logger.Error("[ERROR] Decoding failed, stop processing")

		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Invalid object body, must be dto.LinksRequestData"))
		return
	}

	url := "https://test24.ru/"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		handler.logger.Error("Can't make a proper request")

		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Request creating wasn't successful"))
	}
	req.Header.Set("User-Agent", "")
	req.Header.Add("Cookie", "tester=%D0%98%D0%BD%D0%BA%D0%BE%D0%B3%D0%BD%D0%B8%D1%82%D0%BE")

	resp, err := client.Do(req)
	if err != nil {
		handler.logger.Error("Can't receive response from given source, stop processing", err)

		rw.WriteHeader(http.StatusBadGateway)
		rw.Write([]byte("Response from given source wasn't received. Check your URL or try later"))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		handler.logger.Error("Status code is not OK, stop processing")

		rw.WriteHeader(http.StatusBadGateway)
		rw.Write([]byte("Response from given source is "))
		rw.Write([]byte(resp.Status))
		rw.Write([]byte("\n"))
		io.Copy(rw, resp.Body)

		for k, v := range req.Header {
			rw.Write([]byte(fmt.Sprint(k, " ", v, "\n")))
		}

		cookies := tools.GetCookies(handler.testsClient, rd.Link)
		for _, cookie := range cookies {
			rw.Write([]byte(cookie.String()))
			rw.Write([]byte("\n"))
		}
		return
	}

	rw.WriteHeader(http.StatusOK)
	io.Copy(rw, resp.Body)

	// bodyBytes, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	handler.logger.Error("Can't read response body")

	// 	rw.WriteHeader(http.StatusInternalServerError)
	// 	rw.Write([]byte("Error when reading response body"))
	// 	return
	// }

	// bodyString := string(bodyBytes)
	// entries := getAllHrefPartsFromStringifyBody(bodyString)

	// rw.WriteHeader(http.StatusOK)
	// for _, v := range entries {
	// 	ref := getReferenceFromHref(v)
	// 	rw.Write([]byte(ref))
	// 	rw.Write([]byte("\n"))
	// }
}

func getAllHrefPartsFromStringifyBody(str string) []string {
	reg := regexp.MustCompile(`href="[^"]*://[^"]*"`)
	return reg.FindAllString(str, -1)
}

func getReferenceFromHref(hrefMatch string) string {
	// hrefMatch is a string like `href="required-reference.org"`
	return hrefMatch[6 : len(hrefMatch)-1]
}
