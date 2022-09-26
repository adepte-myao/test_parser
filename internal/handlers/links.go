package handlers

import (
	"io"
	"net/http"
	"regexp"

	"github.com/sirupsen/logrus"
)

type LinksHandler struct {
	logger   *logrus.Logger
	baseLink string
}

func NewLinksHandler(logger *logrus.Logger, baseLink string) *LinksHandler {
	return &LinksHandler{
		logger:   logger,
		baseLink: baseLink,
	}
}

func (handler LinksHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Info("Find all links request received")

	resp, err := doRequest(http.MethodGet, handler.baseLink)
	if err != nil {
		handler.logger.Error("cannot do request: ", err)

		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Response from given source wasn't received. Check your URL or try later"))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		handler.logger.Error("Status code is not OK, stop processing")
		sendErrorResponse(rw, resp)
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		handler.logger.Error("Can't read response body")

		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error when reading response body"))
		return
	}

	bodyString := string(bodyBytes)
	entries := getAllHrefPartsFromStringifyBody(bodyString)
	references := getReferencesFromHref(entries)
	absoluteReferences := getProperReferences(handler.baseLink, references)

	rw.WriteHeader(http.StatusOK)
	for _, ref := range absoluteReferences {
		rw.Write([]byte(ref))
		rw.Write([]byte("\n"))
	}
}

func getAllHrefPartsFromStringifyBody(str string) []string {
	reg := regexp.MustCompile(`href="[^"]*/[^"]*"`)
	return reg.FindAllString(str, -1)
}

func getReferencesFromHref(hrefMatches []string) []string {
	out := make([]string, 0)
	for _, hrefMatch := range hrefMatches {
		// hrefMatch is a string like `href="required-reference.org"`
		cleanRef := hrefMatch[6 : len(hrefMatch)-1]

		out = append(out, cleanRef)
	}

	return out
}

func getProperReferences(base string, allRefs []string) []string {
	regFindRelativesRefs := regexp.MustCompile(`/?iter=.*`)
	out := make([]string, 0)
	for _, ref := range allRefs {
		result := regFindRelativesRefs.FindString(ref)

		if result != "" {
			out = append(out, base+result)
		}
	}

	return out
}

func doRequest(method string, url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "")
	req.Header.Add("Cookie", "tester=%D0%98%D0%BD%D0%BA%D0%BE%D0%B3%D0%BD%D0%B8%D1%82%D0%BE")

	return client.Do(req)
}

func sendErrorResponse(rw http.ResponseWriter, resp *http.Response) {
	rw.WriteHeader(http.StatusBadGateway)
	rw.Write([]byte("Response from given source is "))
	rw.Write([]byte(resp.Status))
	rw.Write([]byte("\n"))
	io.Copy(rw, resp.Body)
}
