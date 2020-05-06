package commands

import (
	"errors"
	"fmt"

	"../util"

	"github.com/attic-labs/noms/go/config"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/spec"
	"github.com/attic-labs/noms/go/types"

	"github.com/spf13/cobra"
)

func runCreateStory(cmd *cobra.Command, args []string) error {
	title := args[0]
	cfg := config.NewResolver() //config default db "Stories"
	db, ds, err := cfg.GetDataset(title)
	d.PanicIfError(err)
	defer db.Close()

	// Create
	var composition = []string{"description", " ", "effort", "0"}
	absPath := util.ApplyStructEdits(db, types.NewStruct(title, nil), nil, composition)

	// Commits
	value := absPath.Resolve(db)
	if value == nil {
		d.CheckErrorNoUsage(errors.New(fmt.Sprintf("Error resolving value: %s", absPath.String())))
	}

	oldCommitRef, oldCommitExists := ds.MaybeHeadRef()
	if oldCommitExists {
		fmt.Printf("Create aborted - %s allready exist (is #%s)\n", title, oldCommitRef.TargetHash().String())
		return nil
	}

	meta, err := spec.CreateCommitMetaStruct(db, "", "Create new story : "+title, nil, nil)
	d.CheckErrorNoUsage(err)

	ds, err = db.Commit(ds, value, datas.CommitOptions{Meta: meta})
	d.CheckErrorNoUsage(err)

	fmt.Printf("%s was created\n", title)

	return nil
}

var createStoryCmd = &cobra.Command{
	Use:   "create <title>",
	Short: "Create a new story.",
	Args:  cobra.ExactArgs(1),
	RunE:  runCreateStory,
}

func init() {
	storyCmd.AddCommand(createStoryCmd)
}
