package elements

import "fmt"

type story struct {
	title string
	description string
	effort int
}

func NewStory(title string) story{
	s := story{title, "", 0}
	return s
}