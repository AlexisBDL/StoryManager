package commands

import (
	"fmt"

	"github.com/attic-labs/noms/go/config"
	"github.com/attic-labs/noms/go/d"

	"github.com/spf13/cobra"
)

// Utile plus tard si notion de Login

func runInit(cmd *cobra.Command, args []string) error {
	cfg := config.NewResolver() //config default db "Stories"
	db, err := cfg.GetDatabase("")
	d.PanicIfError(err)
	defer db.Close()

	fmt.Printf("Stories created\n")

	return nil
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initilze stories. You need too create .nomsconfig",
	RunE:  runInit,
}

func init() {
	RootCmd.AddCommand(initCmd)
}
