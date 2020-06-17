package commands

import (
	"fmt"

	"github.com/AlexisBDL/StoryManager/spec"
	"github.com/AlexisBDL/StoryManager/util"
	"github.com/attic-labs/noms/go/d"

	"github.com/spf13/cobra"
)

func runAddTaskStory(cmd *cobra.Command, args []string) error {
	ID := args[0]
	goal := args[1]
	maker := ""

	if len(args) == 2 {
		maker = user
	} else {
		maker = args[2]
	}

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
		fieldsList []string
	)

	resolvedList := cfg.ResolvePathSpec(ID) + storyTasks
	absPathTask := util.ApplyStructEdits(db, newTask(goal, maker, util.ListLen(resolvedList)), nil, nil)

	fieldsList = append(fieldsList, "@#"+absPathTask.Hash.String())
	absPathList := util.ListAppend(resolvedList, fieldsList)

	absPath, err := spec.NewAbsolutePath("#" + absPathList.Hash.String() + ".value")
	d.CheckError(err)

	absPathTask = util.ApplyStructEdits(db, newTask(goal, maker, util.ListLen(resolvedList)), nil, nil)

	// Commit
	title := getTitle(ID)
	msg := "Add task " + goal + " on story ID " + ID
	valPath := absPath.Resolve(db)

	util.Commit(db, ds, valPath, ID, msg, user, title)

	return nil
}

var addTaskStoryCmd = &cobra.Command{
	Use:   "Tadd <ID> <goal> <maker>",
	Short: "Add a task in story ID.",
	Args:  cobra.MinimumNArgs(2),
	RunE:  runAddTaskStory,
}

func init() {
	storyCmd.AddCommand(addTaskStoryCmd)
}
