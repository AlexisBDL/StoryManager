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

var (
	editTitle       string
	editEffort      int
	editDescription string
)

func runEditStory(cmd *cobra.Command, args []string) error {
	ID := args[0]

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
	resolved := cfg.ResolvePathSpec(ID) + commitStory
	sp, err := spec.ForPath(resolved)
	d.PanicIfError(err)
	defer sp.Close()

	pinned, ok := sp.Pin()
	if !ok {
		fmt.Printf("Cannot resolve spec: %s\n", ID)
		return nil
	}

	db := pinned.GetDatabase()
	ds := db.GetDataset(ID)

	rootVal, basePath := util.SplitPath(sp)

	var (
		absPath *spec.AbsolutePath
		change  string
		field   []string
	)
	if editDescription != "" {
		change += "description "
		field = append(field, "Description", editDescription)
	}
	if editEffort != -1 {
		change += "effort "
		field = append(field, "Effort", strconv.Itoa(editEffort))
	}
	if editTitle != "" {
		change += "title "
		field = append(field, "Title", editTitle)
	}

	absPath = util.ApplyStructEdits(db, rootVal, basePath, field)

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

	message := "Edit value " + change + "in story : " + ID
	meta, err := spec.CreateCommitMetaStruct(db, "", message, user, nil, nil)
	d.CheckErrorNoUsage(err)

	ds, err = db.Commit(ds, valPath, datas.CommitOptions{Meta: meta})
	d.CheckErrorNoUsage(err)

	if oldCommitExists {
		fmt.Printf("New head #%v (was #%v)\n", ds.HeadRef().TargetHash().String(), oldCommitRef.TargetHash().String())
	} else {
		fmt.Printf("New head #%v\n", ds.HeadRef().TargetHash().String())
	}
	fmt.Printf("%s edited\n", ID)

	return nil
}

var editStoryCmd = &cobra.Command{
	Use:   "edit <ID> [flag] <value>",
	Short: "Edit a field of story.",
	Args:  cobra.ExactArgs(1),
	RunE:  runEditStory,
}

func init() {
	storyCmd.AddCommand(editStoryCmd)

	editStoryCmd.Flags().IntVarP(&editEffort, "effort", "e", -1,
		"Provide an effort to evaluate the story",
	)
	editStoryCmd.Flags().StringVarP(&editDescription, "description", "d", "",
		"Provide a message to describe the story, use \"\" to add more than one word",
	)
	editStoryCmd.Flags().StringVarP(&editTitle, "title", "t", "",
		"Change the title of the story",
	)
}
