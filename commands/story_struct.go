package commands

import (
	"github.com/attic-labs/noms/go/types"
)

const (
	storyState  = ".value.State"
	strOpen     = "Open"
	strClose    = "Close"
	commitStory = ".value"
)

func NewStory(title string) types.Struct {
	fields := types.StructData{}

	fields["Title"] = types.String(title)
	fields["Description"] = types.String("")
	fields["Effort"] = types.Number(0)
	fields["State"] = types.String(strOpen)

	return types.NewStruct("Story", fields)
}
