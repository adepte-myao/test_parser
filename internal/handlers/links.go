package handlers

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/adepte-myao/test_parser/internal/models"
	"github.com/adepte-myao/test_parser/internal/storage"
	"github.com/adepte-myao/test_parser/internal/tools"
	"github.com/sirupsen/logrus"
)

type LinksHandler struct {
	logger          *logrus.Logger
	linksRepository *storage.LinkRepository
	baseLink        string
	links           []string
	testLinks       []models.Link
}

func NewLinksHandler(logger *logrus.Logger, baseLink string, store *storage.Store) *LinksHandler {
	linksRepo := storage.NewLinksRepository(store)

	links := make([]string, 0)
	links = append(links, baseLink)

	testLinks := make([]models.Link, 0)

	return &LinksHandler{
		logger:          logger,
		linksRepository: linksRepo,
		baseLink:        baseLink,
		links:           links,
		testLinks:       testLinks,
	}
}

func (handler *LinksHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Info("Find all links request received")

	handler.findAllSourceLinks()
	regTestLink := regexp.MustCompile(`bil=`)
	for _, link := range handler.links {
		if found := regTestLink.FindAllString(link, -1); len(found) == 1 {
			handler.testLinks = append(handler.testLinks, (models.Link)(link))
		}
	}

	handler.linksRepository.DeleteAll()
	handler.linksRepository.CreateRange(handler.testLinks)

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(fmt.Sprintf("Added %d links", len(handler.testLinks))))

	handler.testLinks = make([]models.Link, 0)
	handler.links = make([]string, 0)

	handler.logger.Info("Find all links: processing finished")
}

func (handler *LinksHandler) findAllSourceLinks() {
	currentLinkIndex := -1
	for currentLinkIndex < len(handler.links)-1 {
		if currentLinkIndex == 1000000 {
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
	resp, err := tools.DoProperRequest(http.MethodGet, link)
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
