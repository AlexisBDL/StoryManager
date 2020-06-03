package commands

import (
	"math/rand"
	"strconv"

	"github.com/AlexisBDL/StoryManager/util"
	"github.com/attic-labs/noms/go/d"
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
	absPath := util.ApplyStructEdits(db, newStory(title, user, db), nil, nil)

	// Commits
	msg := "Create new story " + title + " with ID " + ID
	ds := db.GetDataset(ID)

	util.Commit(db, ds, absPath, ID, msg, user, title)

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
