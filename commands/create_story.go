package commands

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/AlexisBDL/StoryManager/spec"
	"github.com/AlexisBDL/StoryManager/util"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/hash"
	"github.com/attic-labs/noms/go/util/datetime"

	"github.com/spf13/cobra"
)

func runCreateStory(cmd *cobra.Command, args []string) error {
	title := args[0]

	r := rand.New(rand.NewSource(99))
	data := []byte(title + datetime.Now().String()[20:28] + strconv.Itoa((r.Int())))
	ID := hash.New(data[:20]).String()

	db, err := cfg.GetDatabase("")
	d.PanicIfError(err)
	defer db.Close()

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

	meta, err := spec.CreateCommitMetaStruct(db, "", "Create new story : "+title, user, nil, nil)
	d.CheckErrorNoUsage(err)

	ds, err = db.Commit(ds, value, datas.CommitOptions{Meta: meta})
	d.CheckErrorNoUsage(err)

	fmt.Printf("%s was created\n", ID)

	return nil
}

var createStoryCmd = &cobra.Command{
	Use:   "create <ID> <title>",
	Short: "Create a new story.",
	Args:  cobra.ExactArgs(1),
	RunE:  runCreateStory,
}

func init() {
	storyCmd.AddCommand(createStoryCmd)
}
