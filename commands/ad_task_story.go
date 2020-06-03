package commands

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/AlexisBDL/StoryManager/spec"
	"github.com/AlexisBDL/StoryManager/util"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/types"

	"github.com/spf13/cobra"
)

func runAddTaskStory(cmd *cobra.Command, args []string) error {
	ID := args[0]

	db, ds, err := cfg.GetDataset(ID)
	d.PanicIfError(err)
	defer db.Close()

	// Test Open
	_, valState, _ := cfg.GetPath(ID + storyState)
	if valState == nil {
		d.CheckErrorNoUsage(fmt.Errorf("Story %s not found in my Stories", ID))
	}
	state, err := strconv.Unquote(types.EncodedValue(valState))
	d.PanicIfError(err)
	if state == stateClose {
		fmt.Printf("The story %s is close, you con't modify it\n", ID)
		return nil
	}

	// Edit
	var (
		fieldsT []string
	)

	fieldsT = append(fieldsT, "1", "name")
	resolved := cfg.ResolvePathSpec(ID) + storyTasks
	absPathT := util.MapEdit(db, resolved, fieldsT)

	absPath, err := spec.NewAbsolutePath("#" + absPathT.Hash.String() + ".value")
	d.CheckError(err)

	// Commit
	valPath := absPath.Resolve(db)
	if valPath == nil {
		d.CheckErrorNoUsage(errors.New(fmt.Sprintf("Error resolving value: %s", absPath.String())))
	}

	oldCommitRef, oldCommitExists := ds.MaybeHeadRef()
	if oldCommitExists {
		head := ds.HeadValue()
		if head.Hash() == valPath.Hash() {
			fmt.Printf("Commit aborted - allow-dupe is set to off and this commit would create a duplicate\n")
			return nil
		}
	}

	_, valTitle, err := cfg.GetPath(ID + storyTitle)
	d.PanicIfError(err)
	title, err := strconv.Unquote(types.EncodedValue(valTitle))
	d.PanicIfError(err)

	message := "Add task in story " + title + " with ID " + ID
	meta, err := spec.CreateCommitMetaStruct(db, "", message, user, nil, nil)
	d.CheckErrorNoUsage(err)

	ds, err = db.Commit(ds, valPath, datas.CommitOptions{Meta: meta})
	d.CheckErrorNoUsage(err)

	if oldCommitExists {
		fmt.Printf("New head #%v (was #%v)\n", ds.HeadRef().TargetHash().String(), oldCommitRef.TargetHash().String())
	} else {
		fmt.Printf("New head #%v\n", ds.HeadRef().TargetHash().String())
	}
	fmt.Printf("%s edited --> add task\nID : %s\n", title, ID)

	return nil
}

var addTaskStoryCmd = &cobra.Command{
	Use:   "add <ID>",
	Short: "Add a task in story.",
	Args:  cobra.ExactArgs(1),
	RunE:  runAddTaskStory,
}

func init() {
	storyCmd.AddCommand(addTaskStoryCmd)
}
