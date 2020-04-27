package commands

import (
	"github.com/attic-labs/noms/go/config"
	"github.com/attic-labs/noms/go/d"

	"github.com/spf13/cobra"
)

// implementation compliqué
func runEditStory(cmd *cobra.Command, args []string) error {
	title := args[0]
	cfg := config.NewResolver() //config default db "Stories"
	db, ds, err := cfg.GetDataset(title)
	d.CheckError(err)
	defer db.Close()

	ds.Database().Close()

	return nil
}

var editStoryCmd = &cobra.Command{
	Use:   "edit <title>",
	Short: "Edit a story.",
	Args:  cobra.ExactArgs(1),
	RunE:  runEditStory,
}

func init() {
	storyCmd.AddCommand(editStoryCmd)
}
