package commands

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/AlexisBDL/StoryManager/util"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/hash"
	"github.com/attic-labs/noms/go/types"
	"github.com/attic-labs/noms/go/util/datetime"

	"github.com/spf13/cobra"
)

var (
	datasetRe   = regexp.MustCompile("^" + datas.DatasetRe.String() + "$")
	choose      bool
	leftChoice  bool
	rightChoice bool
)

// Fusion de branch aillant le même parent
func runMergeStory(cmd *cobra.Command, args []string) error {
	ID1 := args[0]
	ID2 := args[1]

	if isOpenStory(ID1) {
		fmt.Printf("The story %s is close, you can't modify it\n", ID1)
		return nil
	}

	if isOpenStory(ID2) {
		fmt.Printf("The story %s is close, you can't modify it\n", ID1)
		return nil
	}

	_, valTitle, err := cfg.GetPath(ID1 + storyTitle)
	title, err := strconv.Unquote(types.EncodedValue(valTitle))
	d.PanicIfError(err)
	data := []byte(title[:4] + datetime.Now().String()[20:28] + randomString(10))
	newID := hash.New(data[:20]).String()

	// temp --> ID1
	util.SyncStory(ID1, "source", "temp", cfg, false)

	// temp --> ID2
	if util.SyncStory(ID2, "dest", "temp", cfg, true) {
		os.RemoveAll("temp")
		d.CheckErrorNoUsage(fmt.Errorf("Stories are already sync"))
		return nil
	}

	// merge
	tempDB, err := cfg.GetDatabase("temp")
	d.CheckError(err)

	if choose {
		util.MergeStory(tempDB, "source", "dest", "merge", user, "p")
	} else if leftChoice {
		util.MergeStory(tempDB, "source", "dest", "merge", user, "l")
	} else if rightChoice {
		util.MergeStory(tempDB, "source", "dest", "merge", user, "r")
	}

	// sync temp --> merge ==> MyStory --> newID
	util.SyncStory("temp::merge", newID, "Stories", cfg, false)

	dbU1, ds1, _ := cfg.GetDataset(ID1)
	dbU1.Delete(ds1)
	dbU1.Close()

	dbU2, ds2, _ := cfg.GetDataset(ID2)
	dbU2.Delete(ds2)
	dbU2.Close()

	tempDB.Close()
	os.RemoveAll("temp")

	fmt.Printf("New ID : %s\n", newID)

	return nil
}

var mergeStoryCmd = &cobra.Command{
	Use:   "merge <ID1> <ID2> [Flag]",
	Short: "Merge two stories that have similar ref. Choose the way to resolve conflicts whith flag. ID1 is left, ID2 is right",
	Args:  cobra.ExactArgs(2),
	RunE:  runMergeStory,
}

func init() {
	storyCmd.AddCommand(mergeStoryCmd)

	mergeStoryCmd.Flags().BoolVarP(&choose, "choose", "c", false, "Ask me to choose between values if a conflict append")
	mergeStoryCmd.Flags().Lookup("choose").NoOptDefVal = "true"

	mergeStoryCmd.Flags().BoolVarP(&leftChoice, "left", "l", false, "Force update if conflict with value in left")
	mergeStoryCmd.Flags().Lookup("left").NoOptDefVal = "true"

	mergeStoryCmd.Flags().BoolVarP(&rightChoice, "right", "r", false, "Force update if conflict with value in right")
	mergeStoryCmd.Flags().Lookup("right").NoOptDefVal = "true"
}

func resolveDatasets(db datas.Database, leftName, rightName, mergeName string) (leftDS, rightDS, merged datas.Dataset) {
	makeDS := func(dsName string) datas.Dataset {
		if !datasetRe.MatchString(dsName) {
			d.CheckErrorNoUsage(fmt.Errorf("Invalid dataset %s, must match %s", dsName, datas.DatasetRe.String()))
		}
		return db.GetDataset(dsName)
	}
	leftDS = makeDS(leftName)
	rightDS = makeDS(rightName)
	merged = makeDS(mergeName)
	return
}
