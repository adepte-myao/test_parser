package html

import (
	"fmt"
	"regexp"

	"github.com/adepte-myao/test_parser/internal/models"
)

type SitemapParser struct {
	sectionReg     *regexp.Regexp
	sectionNameReg *regexp.Regexp
	sectionLinkReg *regexp.Regexp

	certAreaReg     *regexp.Regexp
	certAreaNameReg *regexp.Regexp
	certAreaLinkReg *regexp.Regexp

	testReg     *regexp.Regexp
	testNameReg *regexp.Regexp
	testLinkReg *regexp.Regexp

	ticketReg     *regexp.Regexp
	ticketLinkReg *regexp.Regexp
	baseLink      string
}

func NewSitemapParser(baseLink string) *SitemapParser {
	return &SitemapParser{
		sectionReg:     regexp.MustCompile(`<div class="card shadow col-sm-12 col-md-.*? col-lg-4">[[:print:][:cntrl:]А-Яа-я№«»]*?</div>`),
		sectionNameReg: regexp.MustCompile(`<h2 class="my-0 font-weight-normal">[А-Яа-я ()]+<\/h2>`),
		sectionLinkReg: regexp.MustCompile(`<a href="/\?iter=[0-9]*?&[s_]*group=[0-9]*?">`),

		certAreaReg:     regexp.MustCompile(`href="/\?iter=[0-9]+&[[:print:]]*?group=[0-9]+"><h4 class="font-weight-normal">[[:print:][:cntrl:]А-Яа-я№«»]*?</h4>`),
		certAreaNameReg: regexp.MustCompile(`normal">[[:print:][:cntrl:]А-Яа-я№«»]*?<`),
		certAreaLinkReg: regexp.MustCompile(`/\?iter=[0-9]+&[[:print:]]*?group=[0-9]+`),

		testReg:     regexp.MustCompile(`href="/\?iter=[0-9]+&[[:print:]]*?test=[0-9]+" ><h4 class="font-weight-normal">[[:print:][:cntrl:]А-Яа-я№«»]*?<`),
		testNameReg: regexp.MustCompile(`normal">[[:print:][:cntrl:]А-Яа-я№«»]*?<`),
		testLinkReg: regexp.MustCompile(`/\?iter=[0-9]+&[[:print:]]*?test=[0-9]+"`),

		ticketReg:     regexp.MustCompile(`col-lg-2"> <div class="card flex-shrink-1 shadow">[[:print:][:cntrl:]А-Яа-я№«»]*?</div>`),
		ticketLinkReg: regexp.MustCompile(`/\?iter=[0-9]+&[[:print:]]*?bil=[0-9]+&[[:print:]]*?test=[0-9]+`),
		baseLink:      baseLink,
	}
}

func (parser *SitemapParser) ParseBasePage(html string) ([]models.Section, error) {
	sectionsStrings := parser.sectionReg.FindAllString(html, -1)
	if len(sectionsStrings) == 0 {
		return nil, fmt.Errorf("no sections found")
	}

	sections := make([]models.Section, 0)
	for _, sectionString := range sectionsStrings {
		sectionName := parser.sectionNameReg.FindString(sectionString)
		sectionName = sectionName[36 : len(sectionName)-5]

		sectionLink := parser.sectionLinkReg.FindString(sectionString)
		sectionLink = sectionLink[9 : len(sectionLink)-2]

		section := models.NewSection(
			sectionName,
			models.Link(sectionLink),
		)

		sections = append(sections, *section)
	}

	return sections, nil
}

func (parser *SitemapParser) ParseSectionPage(html string) ([]models.CertArea, error) {
	certAreaStrings := parser.certAreaReg.FindAllString(html, -1)
	if len(certAreaStrings) == 0 {
		return nil, fmt.Errorf("no cert areas found")
	}

	certAreas := make([]models.CertArea, 0)
	for _, certAreaString := range certAreaStrings {
		certAreaName := parser.certAreaNameReg.FindString(certAreaString)
		certAreaName = certAreaName[8 : len(certAreaName)-1]

		certAreaLink := parser.certAreaLinkReg.FindString(certAreaString)

		certArea := models.NewCertArea(
			certAreaName,
			models.Link(certAreaLink),
		)

		certAreas = append(certAreas, *certArea)
	}

	return certAreas, nil
}

func (parser *SitemapParser) ParseCertAreaPage(html string) ([]models.Test, error) {
	testStrings := parser.testReg.FindAllString(html, -1)
	if len(testStrings) == 0 {
		return nil, fmt.Errorf("no cert areas found")
	}

	tests := make([]models.Test, 0)
	for _, testString := range testStrings {
		testName := parser.testNameReg.FindString(testString)
		testName = testName[8 : len(testName)-1]

		testLink := parser.testLinkReg.FindString(testString)

		test := models.NewTest(
			testName,
			models.Link(testLink),
		)

		tests = append(tests, *test)
	}

	return tests, nil
}

func (parser *SitemapParser) ParseTestPage(html string) ([]models.Link, error) {
	ticketStrings := parser.ticketReg.FindAllString(html, -1)
	if len(ticketStrings) == 0 {
		return nil, fmt.Errorf("no cert areas found")
	}

	links := make([]models.Link, 0)
	for _, ticketString := range ticketStrings {
		ticketLink := parser.ticketLinkReg.FindString(ticketString)

		// base format: https://blabla/
		// result format: /ddd/asd/asd
		// to get proper reference it's required to exclude one /
		link := models.Link(parser.baseLink + ticketLink[1:])

		links = append(links, link)
	}

	return links, nil
}
