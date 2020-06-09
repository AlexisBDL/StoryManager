package commands

import (
	"fmt"
	"strconv"

	"github.com/AlexisBDL/StoryManager/util"
	"github.com/attic-labs/noms/go/d"

	"github.com/spf13/cobra"
)

var (
	editTitle       string
	editEffort      int
	editDescription string
)

func runEditStory(cmd *cobra.Command, args []string) error {
	ID := args[0]

	db, ds, err := cfg.GetDataset(ID)
	d.PanicIfError(err)
	defer db.Close()

	// Test Open
	if isOpenStory(ID) {
		fmt.Printf("The story %s is close, you can't modify it\n", ID)
		return nil
	}

	// Edit
	var (
		change string
		fields []string
	)
	if editDescription != "" {
		change += "description "
		fields = append(fields, "Description", editDescription)
	}
	if editEffort != -1 {
		change += "effort "
		fields = append(fields, "Effort", strconv.Itoa(editEffort))
	}
	if editTitle != "" {
		change += "title "
		fields = append(fields, "Title", editTitle)
	}

	resolved := cfg.ResolvePathSpec(ID) + commitStory
	absPath := util.StoryEdit(db, resolved, fields)

	// Commit
	title := getTitle(ID)
	msg := "Edit value " + change + "in story " + title + " with ID " + ID
	valPath := absPath.Resolve(db)

	util.Commit(db, ds, valPath, ID, msg, user, title)

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
