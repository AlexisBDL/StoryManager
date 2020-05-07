package commands

import (
	"fmt"

	"github.com/attic-labs/noms/go/config"
	"github.com/attic-labs/noms/go/d"

	"github.com/spf13/cobra"
)

func runDeleteStory(cmd *cobra.Command, args []string) error {
	ID := args[0]
	cfg := config.NewResolver() //config default db "Stories"
	db, ds, err := cfg.GetDataset(ID)
	d.CheckError(err)
	defer db.Close()

	oldCommitRef, errBool := ds.MaybeHeadRef()
	if !errBool {
		d.CheckError(fmt.Errorf("Dataset %v not found", ds.ID()))
	}

	_, err = db.Delete(ds)
	d.CheckError(err)
	fmt.Printf("Deleted %v (was #%v)\n", ID, oldCommitRef.TargetHash().String())

	return nil
}

var deleteStoryCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "Delete a story.",
	Args:  cobra.ExactArgs(1),
	RunE:  runDeleteStory,
}

func init() {
	storyCmd.AddCommand(deleteStoryCmd)
}
