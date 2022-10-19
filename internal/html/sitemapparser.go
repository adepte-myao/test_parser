package html

import "github.com/adepte-myao/test_parser/internal/models"

type SitemapParser struct {
}

func NewSitemapParser() *SitemapParser {
	return &SitemapParser{}
}

func (parser *SitemapParser) ParseBasePage(html string) []models.Section {
	return nil
}

func (parser *SitemapParser) ParseSectionPage(html string) []models.CertArea {
	return nil
}

func (parser *SitemapParser) ParseCertAreaPage(html string) []models.Test {
	return nil
}

func (parser *SitemapParser) ParseTestPage(html string) []models.Link {
	return nil
}
