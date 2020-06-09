package commands

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/AlexisBDL/StoryManager/spec"
	"github.com/AlexisBDL/StoryManager/util"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/types"

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
	_, valState, _ := cfg.GetPath(ID + storyState)
	if valState == nil {
		d.CheckErrorNoUsage(fmt.Errorf("Story %s not found in my Stories", ID))
	}
	state, err := strconv.Unquote(types.EncodedValue(valState))
	d.PanicIfError(err)
	if state == stateClose {
		fmt.Printf("The story %s is close, you con't modify it\n", ID)
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

	absPathTask = util.ApplyStructEdits(db, util.ListGet(resolvedList, uint64(IDX)), nil, fields)

	// Commit
	valPath := absPath.Resolve(db)
	if valPath == nil {
		d.CheckErrorNoUsage(errors.New(fmt.Sprintf("Error resolving value: %s", absPath.String())))
	}

	oldCommitRef, oldCommitExists := ds.MaybeHeadRef()
	if oldCommitExists {
		head := ds.HeadValue()
		if head.Hash() == valPath.Hash() {
			fmt.Printf("Commit aborted - allow-dupe is set to off and this commit would create a duplicate\n")
			return nil
		}
	}

	_, valTitle, err := cfg.GetPath(ID + storyTitle)
	d.PanicIfError(err)
	title, err := strconv.Unquote(types.EncodedValue(valTitle))
	d.PanicIfError(err)

	message := "Add task in story " + title + " with ID " + ID
	meta, err := spec.CreateCommitMetaStruct(db, "", message, user, nil, nil)
	d.CheckErrorNoUsage(err)

	ds, err = db.Commit(ds, valPath, datas.CommitOptions{Meta: meta})
	d.CheckErrorNoUsage(err)

	if oldCommitExists {
		fmt.Printf("New head #%v (was #%v)\n", ds.HeadRef().TargetHash().String(), oldCommitRef.TargetHash().String())
	} else {
		fmt.Printf("New head #%v\n", ds.HeadRef().TargetHash().String())
	}
	fmt.Printf("%s edited --> add task\nID : %s\n", title, ID)

	return nil
}

var editTaskStoryCmd = &cobra.Command{
	Use:   "Tedit <ID> <IDTask>",
	Short: "Edit a task in story.",
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
