package commands

import (
	"github.com/attic-labs/noms/go/types"
)

const (
	storyState  = ".value.State"
	storyTitle  = ".value.Title"
	stateOpen   = "Open"
	stateClose  = "Close"
	commitStory = ".value"
)

func NewStory(title string, author string) types.Struct {
	fields := types.StructData{}

	fields["Title"] = types.String(title)
	fields["Description"] = types.String("")
	fields["Effort"] = types.Number(0)
	fields["State"] = types.String(stateOpen)
	fields["Author"] = types.String(author)

	return types.NewStruct("Story", fields)
}
