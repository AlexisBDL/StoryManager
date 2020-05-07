package commands

import (
	"fmt"
	"strconv"

	"github.com/attic-labs/noms/go/config"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/types"
	"github.com/attic-labs/noms/go/util/outputpager"

	"github.com/spf13/cobra"
)

var (
	open  bool
	close bool
)

func runListStory(cmd *cobra.Command, args []string) error {
	cfg := config.NewResolver() //config default db "Stories"
	db, err := cfg.GetDatabase("")
	d.CheckError(err)
	defer db.Close()

	var ls []string
	db.Datasets().IterAll(func(k, v types.Value) {
		ls = append(ls, fmt.Sprint(k))
	})

	pgr := outputpager.Start()
	defer pgr.Stop()

	var lsOpen []string
	var lsClose []string
	for _, v := range ls {
		_, stat, _ := cfg.GetPath(v + ".value.Stat")
		str, err := strconv.Unquote(types.EncodedValue(stat))
		d.PanicIfError(err)
		if str == "Open" {
			lsOpen = append(lsOpen, v)
		}
		if str == "Close" {
			lsClose = append(lsClose, v)
		}
	}

	if open {
		for _, v := range lsOpen {
			fmt.Println(v)
		}
	}
	if close {
		for _, v := range lsClose {
			fmt.Println(v)
		}
	}

	if !close && !open {
		for _, v := range ls {
			fmt.Println(v)
		}
	}

	return nil
}

var listStoryCmd = &cobra.Command{
	Use:   "list",
	Short: "Display stories in database Stories.",
	Args:  cobra.ExactArgs(0),
	RunE:  runListStory,
}

func init() {
	RootCmd.AddCommand(listStoryCmd)

	listStoryCmd.Flags().BoolVarP(&open, "open", "o", false, "Display open stories")
	listStoryCmd.Flags().Lookup("open").NoOptDefVal = "true"

	listStoryCmd.Flags().BoolVarP(&close, "close", "c", false, "Display close stories")
	listStoryCmd.Flags().Lookup("close").NoOptDefVal = "true"
}
