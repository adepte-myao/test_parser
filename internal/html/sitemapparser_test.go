package html_test

import (
	"regexp"
	"testing"
)

func TestSiteMapParser_ParseBaseBage(t *testing.T) {
	sections, err := sitemapParser.ParseBasePage(basePage)
	if err != nil {
		t.Fatal(err)
	}

	validatorReg := regexp.MustCompile(`[\<>A-Za-z"']*`)
	for _, section := range sections {
		nameValidationString := validatorReg.FindString(section.Name)
		linkValidationString := validatorReg.FindString(string(section.Link))

		if nameValidationString != "" {
			t.Fatal("forbidden symbol found in section.name")
		}
		if linkValidationString != "" {
			t.Fatal("forbidden symbol found in section.link")
		}
	}
}

func TestSiteMapParser_ParseSectionPage(t *testing.T) {
	certAreas, err := sitemapParser.ParseSectionPage(sectionPage)
	if err != nil {
		t.Fatal(err)
	}

	validatorReg := regexp.MustCompile(`[\<>A-Za-z"']*`)
	for _, cerArea := range certAreas {
		nameValidationString := validatorReg.FindString(cerArea.Name)
		linkValidationString := validatorReg.FindString(string(cerArea.Link))

		if nameValidationString != "" {
			t.Fatal("forbidden symbol found in certArea.name")
		}
		if linkValidationString != "" {
			t.Fatal("forbidden symbol found in certArea.link")
		}
	}
}

func TestSiteMapParser_ParseCertAreaPage(t *testing.T) {
	tests, err := sitemapParser.ParseCertAreaPage(certAreaPage)
	if err != nil {
		t.Fatal(err)
	}

	validatorReg := regexp.MustCompile(`[\<>A-Za-z"']*`)
	for _, test := range tests {
		nameValidationString := validatorReg.FindString(test.Name)
		linkValidationString := validatorReg.FindString(string(test.Link))

		if nameValidationString != "" {
			t.Fatal("forbidden symbol found in test.name")
		}
		if linkValidationString != "" {
			t.Fatal("forbidden symbol found in test.link")
		}
	}
}

func TestSiteMapParser_ParseTestPage(t *testing.T) {
	links, err := sitemapParser.ParseTestPage(testPage)
	if err != nil {
		t.Fatal(err)
	}

	validatorReg := regexp.MustCompile(`[\<>A-Za-z"']*`)
	for _, link := range links {
		linkValidationString := validatorReg.FindString(string(link))

		if linkValidationString != "" {
			t.Fatal("forbidden symbol found in certArea.link")
		}
	}
}
