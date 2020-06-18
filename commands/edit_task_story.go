package commands

import (
	"fmt"
	"strconv"

	"github.com/AlexisBDL/StoryManager/spec"
	"github.com/AlexisBDL/StoryManager/util"
	"github.com/attic-labs/noms/go/d"

	"github.com/spf13/cobra"
)

var (
	editGoal  string
	editMaker string
	editState string
)

func runEditTaskStory(cmd *cobra.Command, args []string) error {
	ID := args[0]
	IDT := args[1]

	IDX, err := strconv.Atoi(IDT)
	d.PanicIfError(err)

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
	if editGoal != "" {
		change += "goal "
		fields = append(fields, "Goal", editGoal)
	}
	if editMaker != "" {
		change += "maker "
		fields = append(fields, "Maker", editMaker)
	}
	if editState != "" {
		change += "state "
		fields = append(fields, "State", editState)
	}

	var (
		fieldsList []string
	)

	resolvedList := cfg.ResolvePathSpec(ID) + storyTasks
	absPathTask := util.ApplyStructEdits(db, util.ListGet(resolvedList, uint64(IDX)), nil, fields)

	fieldsList = append(fieldsList, "@#"+absPathTask.Hash.String())
	absPathDelT := util.ListDel(db, resolvedList, IDX)
	resolvedListAfterDel := cfg.ResolvePathSpec(absPathDelT.String())
	absPathInsT := util.ListInsert(db, resolvedListAfterDel, IDX, fieldsList)

	absPath, err := spec.NewAbsolutePath("#" + absPathInsT.Hash.String() + ".value")
	d.CheckError(err)

	absPathTask = util.ApplyStructEdits(db, util.ListGet(resolvedList, uint64(IDX)), nil, fields) //Remake -> Flush

	// Commit
	title := getTitle(ID)
	msgCommit := "Edit task : " + IDT
	msgCli := "Task " + IDT + " edited in story " + title + "\nID : " + ID
	valPath := absPath.Resolve(db)

	util.Commit(db, ds, valPath, ID, msgCommit, msgCli, user)

	return nil
}

var editTaskStoryCmd = &cobra.Command{
	Use:   "Tedit <ID> <IDTask> [flag] <value>",
	Short: "Edit a task IDTask in story ID.",
	Args:  cobra.ExactArgs(2),
	RunE:  runEditTaskStory,
}

func init() {
	storyCmd.AddCommand(editTaskStoryCmd)

	editTaskStoryCmd.Flags().StringVarP(&editGoal, "goal", "g", "",
		"Provide a goal of a task, use \"\" to add more than one word",
	)
	editTaskStoryCmd.Flags().StringVarP(&editMaker, "maker", "m", "",
		"Change the maker of the task",
	)
	editTaskStoryCmd.Flags().StringVarP(&editState, "state", "s", "",
		"Change the state of the task",
	)
}
