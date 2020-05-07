package commands

import (
	"github.com/attic-labs/noms/go/types"
)

func NewStory(title string) types.Struct {
	fields := types.StructData{}

	fields["Title"] = types.String(title)
	fields["Description"] = types.String("")
	fields["Effort"] = types.Number(0)
	fields["Stat"] = types.String("Open")

	return types.NewStruct("Story", fields)
}
