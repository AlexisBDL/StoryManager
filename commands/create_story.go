package commands

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/AlexisBDL/StoryManager/util"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/hash"
	"github.com/attic-labs/noms/go/util/datetime"

	"github.com/spf13/cobra"
)

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(rand.Intn(99))
	}
	return string(bytes)
}

func runCreateStory(cmd *cobra.Command, args []string) error {
	title := args[0]

	if 4 < len(title) {
		d.CheckErrorNoUsage(fmt.Errorf("Title of story need to be more long than 4 characters"))
	}

	rand.Seed(time.Now().UTC().UnixNano())
	data := []byte(title[:4] + datetime.Now().String()[20:28] + randomString(10))
	ID := hash.New(data[:20]).String()

	db, err := cfg.GetDatabase("")
	d.PanicIfError(err)
	defer db.Close()

	// Create
	absPath := util.ApplyStructEdits(db, newStory(title, user, db), nil, nil)

	// Commits
	msg := "Create new story " + title + " with ID " + ID
	ds := db.GetDataset(ID)
	valPath := absPath.Resolve(db)

	util.Commit(db, ds, valPath, ID, msg, user, title)

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
