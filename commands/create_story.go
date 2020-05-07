package commands

import (
	"errors"
	"fmt"

	"github.com/attic-labs/noms/go/config"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/spec"

	"github.com/spf13/cobra"
)

func runCreateStory(cmd *cobra.Command, args []string) error {
	ID := args[0]
	cfg := config.NewResolver() //config default db "Stories"
	db, ds, err := cfg.GetDataset(ID)
	d.PanicIfError(err)
	defer db.Close()

	// Create
	absPath := ApplyStructEdits(db, NewStory(ID), nil, nil)

	// Commits
	value := absPath.Resolve(db)
	if value == nil {
		d.CheckErrorNoUsage(errors.New(fmt.Sprintf("Error resolving value: %s", absPath.String())))
	}

	oldCommitRef, oldCommitExists := ds.MaybeHeadRef()
	if oldCommitExists {
		fmt.Printf("Create aborted - %s allready exist (is #%s)\n", ID, oldCommitRef.TargetHash().String())
		return nil
	}

	meta, err := spec.CreateCommitMetaStruct(db, "", "Create new story : "+ID, nil, nil)
	d.CheckErrorNoUsage(err)

	ds, err = db.Commit(ds, value, datas.CommitOptions{Meta: meta})
	d.CheckErrorNoUsage(err)

	fmt.Printf("%s was created\n", ID)

	return nil
}

var createStoryCmd = &cobra.Command{
	Use:   "create <ID>",
	Short: "Create a new story.",
	Args:  cobra.ExactArgs(1),
	RunE:  runCreateStory,
}

func init() {
	storyCmd.AddCommand(createStoryCmd)
}
