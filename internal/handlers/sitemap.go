package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/adepte-myao/test_parser/internal/html"
	"github.com/adepte-myao/test_parser/internal/models"
	"github.com/adepte-myao/test_parser/internal/storage"
	"github.com/adepte-myao/test_parser/internal/tools"
	"github.com/sirupsen/logrus"
)

type SitemapHandler struct {
	logger            *logrus.Logger
	sections          []models.Section
	sitemapRepository *storage.SitemapRepository
	sitemapParser     *html.SitemapParser
	baseLink          models.Link
}

func NewSitemapHandler(logger *logrus.Logger, baseLink string, store *storage.Store) *SitemapHandler {
	return &SitemapHandler{
		logger:            logger,
		sitemapRepository: storage.NewSitemapRepository(store),
		sitemapParser:     html.NewSitemapParser(),
		baseLink:          models.NewLink(baseLink),
		sections:          make([]models.Section, 0),
	}
}

func (handler *SitemapHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Info("Make sitemap request received")

	basePage, err := handler.getStringifySourceBody(handler.baseLink)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	}

	handler.sections, err = handler.sitemapParser.ParseBasePage(basePage)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	}

	handler.excludeArchive()

	for i := 0; i < len(handler.sections); i++ {
		err = handler.fillSectionLinks(&handler.sections[i])
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
		}
	}

	// Before filling the tables old values should be removed
	err = handler.sitemapRepository.TruncateAllSitemapTables()
	if err != nil {
		handler.logger.Error("Cannot truncate all sitemap tables: ", err.Error())
		return
	}

	handler.sitemapRepository.CreateFilledSections(handler.sections)
}

func (handler *SitemapHandler) getStringifySourceBody(link models.Link) (string, error) {
	stringLink := string(link)
	resp, err := tools.DoProperRequest(http.MethodGet, stringLink)
	if err != nil {
		handler.logger.Error("cannot do request: ", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code from source %s is not OK", link)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		handler.logger.Error("Can't read response body")
		return "", err
	}

	bodyString := string(bodyBytes)

	return bodyString, nil
}

func (handler *SitemapHandler) excludeArchive() {

}

func (handler *SitemapHandler) fillSectionLinks(section *models.Section) error {
	sectionPage, err := handler.getStringifySourceBody(section.Link)
	if err != nil {
		return err
	}

	section.CertAreas, err = handler.sitemapParser.ParseSectionPage(sectionPage)
	if err != nil {
		return err
	}

	for areaIndex := 0; areaIndex < len(section.CertAreas); areaIndex++ {
		err = handler.fillCertAreaLinks(&section.CertAreas[areaIndex])
		if err != nil {
			return err
		}
	}

	return nil
}

func (handler *SitemapHandler) fillCertAreaLinks(certArea *models.CertArea) error {
	certAreaPage, err := handler.getStringifySourceBody(certArea.Link)
	if err != nil {
		return err
	}

	certArea.Tests, err = handler.sitemapParser.ParseCertAreaPage(certAreaPage)
	if err != nil {
		return err
	}

	for testIndex := 0; testIndex < len(certArea.Tests); testIndex++ {
		err = handler.fillTestLinks(&certArea.Tests[testIndex])
		if err != nil {
			return err
		}
	}

	return nil
}

func (handler *SitemapHandler) fillTestLinks(test *models.Test) error {
	testPage, err := handler.getStringifySourceBody(test.Link)
	if err != nil {
		return err
	}

	test.TicketLinks, err = handler.sitemapParser.ParseTestPage(testPage)
	if err != nil {
		return err
	}

	return nil
}
