package models

type Section struct {
	Name      string
	Link      Link
	CertAreas []CertArea
}

type CertArea struct {
	Name  string
	Link  Link
	Tests []Test
}

type Test struct {
	Name        string
	Link        Link
	TicketLinks []Link
}

func NewSection(name string, link Link) *Section {
	return &Section{
		Name:      name,
		Link:      link,
		CertAreas: make([]CertArea, 0),
	}
}

func NewCertArea(name string, link Link) *CertArea {
	return &CertArea{
		Name:  name,
		Link:  link,
		Tests: make([]Test, 0),
	}
}

func NewTest(name string, link Link) *Test {
	return &Test{
		Name:        name,
		Link:        link,
		TicketLinks: make([]Link, 0),
	}
}
