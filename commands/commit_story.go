package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/attic-labs/noms/go/config"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/spec"

	"github.com/spf13/cobra"
)

func runCommitStory(cmd *cobra.Command, args []string) error {
	cfg := config.NewResolver()
	path := args[0]    //path ou hash
	dataSet := args[1] //story
	message := args[3]
	db, ds, err := cfg.GetDataset("Stories::" + dataSet)
	d.CheckError(err)
	defer db.Close()

	absPath, err := spec.NewAbsolutePath(path)
	d.CheckError(err)

	value := absPath.Resolve(db)
	if value == nil {
		d.CheckErrorNoUsage(errors.New(fmt.Sprintf("Error resolving value: %s", args[0])))
	}

	oldCommitRef, oldCommitExists := ds.MaybeHeadRef()
	if oldCommitExists {
		head := ds.HeadValue()
		if head.Hash() == value.Hash() {
			fmt.Fprintf(os.Stdout, "Commit aborted - allow-dupe is set to off and this commit would create a duplicate\n")
			return nil
		}
	}

	meta, err := spec.CreateCommitMetaStruct(db, "", message, nil, nil)
	d.CheckErrorNoUsage(err)

	ds, err = db.Commit(ds, value, datas.CommitOptions{Meta: meta})
	d.CheckErrorNoUsage(err)

	if oldCommitExists {
		fmt.Fprintf(os.Stdout, "New head #%v (was #%v)\n", ds.HeadRef().TargetHash().String(), oldCommitRef.TargetHash().String())
	} else {
		fmt.Fprintf(os.Stdout, "New head #%v\n", ds.HeadRef().TargetHash().String())
	}

	return nil
}

var commitStoryCmd = &cobra.Command{
	Use:   "commit <#path> <story> <message>",
	Short: "Commit a story.",
	Args:  cobra.MinimumNArgs(3),
	RunE:  runCommitStory,
}

func init() {
	RootCmd.AddCommand(commitStoryCmd)
}
