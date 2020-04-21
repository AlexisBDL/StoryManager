package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/spec"

	"github.com/AlexisBDL/StoryManager/elements"
)

func runCreateStory(cmd *cobra.Command, args []string) error {
	sp, err := spec.ForDatabase("Stories")
	d.PanicIfError(err)
	applyStructEdits(sp, NewStory(args[0]), nil, args)
	fmt.Printf("%s created\n", args[0])

	return nil
}

func applyStructEdits(sp spec.Spec, rootVal types.Value, basePath types.Path, args []string) {
	if len(args)%2 != 0 {
		d.CheckError(fmt.Errorf("Must be an even number of key/value pairs"))
	}
	if rootVal == nil {
		d.CheckErrorNoUsage(fmt.Errorf("No value at: %s", sp.String()))
		return
	}
	db := sp.GetDatabase()
	patch := diff.Patch{}
	for i := 0; i < len(args); i += 2 {
		if !types.IsValidStructFieldName(args[i]) {
			d.CheckError(fmt.Errorf("Invalid field name: %s at position: %d", args[i], i))
		}
		nv, err := argumentToValue(args[i+1], db)
		if err != nil {
			d.CheckError(fmt.Errorf("Invalid field value: %s at position %d: %s", args[i+1], i+1, err))
		}
		patch = append(patch, diff.Difference{
			Path:       append(basePath, types.FieldPath{Name: args[i]}),
			ChangeType: types.DiffChangeModified,
			NewValue:   nv,
		})
	}
	appplyPatch(sp, rootVal, basePath, patch)
}

var createStoryCmd = &cobra.Command{
	Use:     "create <title>",
	Short:   "Create a new story.",
	Args:     cobra.MinimumNArgs(1),
	//PreRunE: loadRepoEnsureUser,
	RunE:    runCreateStory,
}

func init() {
	RootCmd.AddCommand(createStoryCmd)

	createStoryCmd.Flags().SortFlags = false
}
