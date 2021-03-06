package commands

import (
	"fmt"
	"strconv"

	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/types"
)

const (
	storyState  = ".value.State"
	storyTitle  = ".value.Title"
	storyTasks  = ".value.Tasks"
	stateOpen   = "Open"
	stateClose  = "Close"
	stateEc     = "Encours"
	stateTt     = "Test"
	stateTr     = "Terminé"
	commitStory = ".value"
)

func newStory(title string, author string, db datas.Database) types.Struct {
	fields := types.StructData{}

	fields["Title"] = types.String(title)
	fields["Description"] = types.String("")
	fields["Effort"] = types.Number(0)
	fields["State"] = types.String(stateOpen)
	fields["Author"] = types.String(author)
	fields["Tasks"] = types.NewList(db)

	return types.NewStruct("Story", fields)
}

func newTask(goal string, maker string, it uint64) types.Struct {
	fields := types.StructData{}

	fields["Goal"] = types.String(goal)
	fields["Maker"] = types.String(maker)
	fields["State"] = types.String("")
	fields["ID"] = types.Number(it)

	return types.NewStruct("Task", fields)
}

func getTitle(ID string) string {
	_, valTitle, err := cfg.GetPath(ID + storyTitle)
	d.PanicIfError(err)
	title, err := strconv.Unquote(types.EncodedValue(valTitle))
	d.PanicIfError(err)

	return title
}

func isOpenStory(ID string) bool {
	_, valState, _ := cfg.GetPath(ID + storyState)
	if valState == nil {
		d.CheckErrorNoUsage(fmt.Errorf("Story %s not found in my Stories", ID))
	}
	state, err := strconv.Unquote(types.EncodedValue(valState))
	d.PanicIfError(err)

	return state == stateClose
}
