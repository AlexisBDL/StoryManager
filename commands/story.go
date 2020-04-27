package commands

import (
	"github.com/spf13/cobra"
)

//Lister ds pour afficher les stories dans l'annuaire

var storyCmd = &cobra.Command{
	Use:   "story",
	Short: "Create, delete or set a story",
}

func init() {
	RootCmd.AddCommand(storyCmd)
}
