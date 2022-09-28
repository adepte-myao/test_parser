package handlers

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/sirupsen/logrus"
)

type LinksHandler struct {
	logger    *logrus.Logger
	baseLink  string
	links     []string
	testLinks []string
}

func NewLinksHandler(logger *logrus.Logger, baseLink string) *LinksHandler {
	links := make([]string, 0)
	links = append(links, baseLink)

	testLinks := make([]string, 0)

	return &LinksHandler{
		logger:    logger,
		baseLink:  baseLink,
		links:     links,
		testLinks: testLinks,
	}
}

func (handler *LinksHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Info("Find all links request received")

	handler.findAllSourceLinks()
	regTestLink := regexp.MustCompile(`bil=`)
	for _, link := range handler.links {
		if found := regTestLink.FindAllString(link, -1); len(found) == 1 {
			handler.testLinks = append(handler.testLinks, link)
		}
	}

	rw.WriteHeader(http.StatusOK)
	for _, link := range handler.testLinks {
		rw.Write([]byte(link))
		rw.Write([]byte("\n"))
	}

	handler.logger.Info("Find all links: processing finished")
}

func (handler *LinksHandler) findAllSourceLinks() {
	currentLinkIndex := -1
	for currentLinkIndex < len(handler.links)-1 {
		if currentLinkIndex == 100000 {
			break
		}
		currentLinkIndex++

		handler.logger.Info("Processing link:", handler.links[currentLinkIndex])

		pageLinks, err := handler.processLink(handler.links[currentLinkIndex])
		if err != nil {
			handler.logger.Warn("Link: ", handler.links[currentLinkIndex], " error when processing")
			continue
		}

		for _, link := range pageLinks {
			if handler.linksContain(link) {
				continue
			}

			handler.links = append(handler.links, link)
		}
	}
}

func (handler *LinksHandler) linksContain(link string) bool {
	for _, existingLink := range handler.links {
		if existingLink == link {
			return true
		}
	}
	return false
}

func (handler *LinksHandler) processLink(link string) ([]string, error) {
	resp, err := doRequest(http.MethodGet, link)
	if err != nil {
		handler.logger.Error("cannot do request: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code from source %s is not OK", link)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		handler.logger.Error("Can't read response body")
		return nil, err
	}

	bodyString := string(bodyBytes)
	entries := getAllHrefPartsFromStringifyBody(bodyString)
	references := getReferencesFromHref(entries)
	absoluteReferences := getProperReferences(handler.baseLink, references)

	return absoluteReferences, nil
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
	regFindRelativesRefs := regexp.MustCompile(`/\?iter=.*`)

	out := make([]string, 0)
	for _, ref := range allRefs {
		result := regFindRelativesRefs.FindString(ref)

		if result != "" {
			// base format: https://blabla/
			// result format: /ddd/asd/asd
			// to get proper reference it's required to exclude one /
			out = append(out, base+(result[1:]))
		}
	}

	regExcludeSemiCol := regexp.MustCompile(`amp|;`)
	for i := 0; i < len(out); i++ {
		out[i] = regExcludeSemiCol.ReplaceAllString(out[i], "")
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
