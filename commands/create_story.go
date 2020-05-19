package commands

import (
	"errors"
	"fmt"

	"github.com/AlexisBDL/StoryManager/config"
	"github.com/AlexisBDL/StoryManager/spec"
	"github.com/AlexisBDL/StoryManager/util"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/types"

	"github.com/spf13/cobra"
)

func runCreateStory(cmd *cobra.Command, args []string) error {
	ID := args[0]
	title := args[1]
	cfg := config.NewResolver() //config default db "Stories"
	user := cfg.GetUserString()
	db, err := cfg.GetDatabase("")
	d.PanicIfError(err)
	defer db.Close()

	// Check no duplicate story
	db.Datasets().IterAll(func(k, v types.Value) {
		if fmt.Sprint(k) == ID {
			d.CheckErrorNoUsage(errors.New(fmt.Sprintf("Error, ID : %s allready exist in database", ID)))
		}
	})

	// Create
	absPath := util.ApplyStructEdits(db, NewStory(title, user), nil, nil)

	// Commits
	value := absPath.Resolve(db)
	if value == nil {
		d.CheckErrorNoUsage(errors.New(fmt.Sprintf("Error resolving value: %s", absPath.String())))
	}

	ds := db.GetDataset(ID)

	oldCommitRef, oldCommitExists := ds.MaybeHeadRef()
	if oldCommitExists {
		fmt.Printf("Create aborted - %s allready exist (is #%s)\n", ID, oldCommitRef.TargetHash().String())
		return nil
	}

	meta, err := spec.CreateCommitMetaStruct(db, "", "Create new story : "+ID, user, nil, nil)
	d.CheckErrorNoUsage(err)

	ds, err = db.Commit(ds, value, datas.CommitOptions{Meta: meta})
	d.CheckErrorNoUsage(err)

	fmt.Printf("%s was created\n", ID)

	return nil
}

var createStoryCmd = &cobra.Command{
	Use:   "create <ID> <title>",
	Short: "Create a new story.",
	Args:  cobra.ExactArgs(2),
	RunE:  runCreateStory,
}

func init() {
	storyCmd.AddCommand(createStoryCmd)
}
