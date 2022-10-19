package models

type Link string

func NewLink(url string) Link {
	return Link(url)
}
