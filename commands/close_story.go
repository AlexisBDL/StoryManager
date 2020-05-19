package commands

import (
	"errors"
	"fmt"

	"github.com/AlexisBDL/StoryManager/util"

	"github.com/AlexisBDL/StoryManager/spec"

	"github.com/attic-labs/noms/go/datas"

	"github.com/attic-labs/noms/go/d"
	"github.com/spf13/cobra"
)

func runCloseStory(cmd *cobra.Command, args []string) error {
	ID := args[0]

	// Edit close
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

	field := []string{"State", stateClose}

	absPath := util.ApplyStructEdits(db, rootVal, basePath, field)

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

	message := ID + " was closed"
	meta, err := spec.CreateCommitMetaStruct(db, "", message, user, nil, nil)
	d.CheckErrorNoUsage(err)

	ds, err = db.Commit(ds, valPath, datas.CommitOptions{Meta: meta})
	d.CheckErrorNoUsage(err)

	if oldCommitExists {
		fmt.Printf("New head #%v (was #%v)\n", ds.HeadRef().TargetHash().String(), oldCommitRef.TargetHash().String())
	} else {
		fmt.Printf("New head #%v\n", ds.HeadRef().TargetHash().String())
	}
	fmt.Printf("%s closed\n", ID)

	return nil
}

var closeStoryCmd = &cobra.Command{
	Use:   "close",
	Short: "Close a story",
	Args:  cobra.ExactArgs(1),
	RunE:  runCloseStory,
}

func init() {
	storyCmd.AddCommand(closeStoryCmd)
}
