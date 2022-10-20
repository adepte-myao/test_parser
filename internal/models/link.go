package models

type Link string

func NewLink(url string) Link {
	return Link(url)
}

type TestLink struct {
	TestId int
	Link   Link
}
