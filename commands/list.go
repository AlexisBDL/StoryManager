package commands

import (
	"fmt"
	"strconv"

	"github.com/attic-labs/noms/go/config"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/types"

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

	var (
		lsOpen   []string
		lsClose  []string
		valState types.Value
		str      string
	)

	for _, v := range ls {
		_, valState, _ = cfg.GetPath(v + storyState)
		str, err = strconv.Unquote(types.EncodedValue(valState))
		d.PanicIfError(err)
		if str == strOpen {
			lsOpen = append(lsOpen, v)
		}
		if str == strClose {
			lsClose = append(lsClose, v)
		}
	}

	if open {
		for _, v := range lsOpen {
			fmt.Println(v)
		}
	} else if close {
		for _, v := range lsClose {
			fmt.Println(v)
		}
	} else if !close && !open {
		for _, v := range ls {
			fmt.Println(v)
		}
	}

	return nil
}

var listStoryCmd = &cobra.Command{
	Use:   "list",
	Short: "Display stories in database Stories.",
	Long:  "It's possible to filtred the displayed stories whith flags",
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
