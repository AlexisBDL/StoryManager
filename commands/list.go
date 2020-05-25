package commands

import (
	"fmt"
	"strconv"

	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/types"

	"github.com/spf13/cobra"
)

var (
	isOpen   bool
	isClose  bool
	dbTarget string
)

func runListStory(cmd *cobra.Command, args []string) error {

	if dbTarget != "" {
		find, err := cfg.FindDatabase(dbTarget)
		if !find {
			fmt.Printf("Database not found at %s\n", dbTarget)
			d.PanicIfError(err)
			return nil
		}
	}

	db, err := cfg.GetDatabase(dbTarget)
	d.CheckError(err)
	defer db.Close()

	var (
		valState types.Value
		state    string
		valTitle types.Value
		title    string
		ID       string
	)

	ls := make(map[string]string)
	lsOpen := make(map[string]string)
	lsClose := make(map[string]string)

	db.Datasets().IterAll(func(k, v types.Value) {
		ID = fmt.Sprint(k)
		_, valTitle, err = cfg.GetPath(ID + storyTitle)
		title, err = strconv.Unquote(types.EncodedValue(valTitle))
		d.PanicIfError(err)
		ls[ID] = title
	})

	for k := range ls {
		_, valState, err = cfg.GetPath(k + storyState)
		state, err = strconv.Unquote(types.EncodedValue(valState))
		d.PanicIfError(err)
		if state == stateOpen {
			lsOpen[k] = ls[k]
		}
		if state == stateClose {
			lsClose[k] = ls[k]
		}
	}

	if isOpen {
		for k, v := range lsOpen {
			fmt.Println(k + "\t\t" + v)
		}
	}
	if isClose {
		for k, v := range lsClose {
			fmt.Println(k + "\t\t" + v)
		}
	}
	if !isClose && !isOpen {
		for k, v := range ls {
			fmt.Println(k + "\t\t" + v)
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

	listStoryCmd.Flags().BoolVarP(&isOpen, "open", "o", false, "Display open stories")
	listStoryCmd.Flags().Lookup("open").NoOptDefVal = "true"

	listStoryCmd.Flags().BoolVarP(&isClose, "close", "c", false, "Display close stories")
	listStoryCmd.Flags().Lookup("close").NoOptDefVal = "true"

	listStoryCmd.Flags().StringVarP(&dbTarget, "db", "d", "", "Display stories in other database path")
	listStoryCmd.Flags().Lookup("db").NoOptDefVal = ""
}
