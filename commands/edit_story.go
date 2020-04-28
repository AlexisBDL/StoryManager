package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/attic-labs/noms/go/config"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/diff"
	"github.com/attic-labs/noms/go/spec"
	"github.com/attic-labs/noms/go/types"

	"github.com/spf13/cobra"
)

func applyStructEditsSp(sp spec.Spec, rootVal types.Value, basePath types.Path, args []string) (path string) {
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
	return appplyPatchSp(sp, rootVal, basePath, patch)
}

func appplyPatchSp(sp spec.Spec, rootVal types.Value, basePath types.Path, patch diff.Patch) (path string) {
	db := sp.GetDatabase()
	baseVal := basePath.Resolve(rootVal, db)
	if baseVal == nil {
		d.CheckErrorNoUsage(fmt.Errorf("No value at: %s", sp.String()))
	}

	newRootVal := diff.Apply(rootVal, patch)
	d.Chk.NotNil(newRootVal)
	r := db.WriteValue(newRootVal)
	db.Flush()
	newAbsPath := spec.AbsolutePath{
		Hash: r.TargetHash(),
		Path: basePath,
	}
	return newAbsPath.String()
}

func splitPath(sp spec.Spec) (rootVal types.Value, basePath types.Path) {
	db := sp.GetDatabase()
	rootPath := sp.Path
	rootPath.Path = types.Path{}
	rootVal = rootPath.Resolve(db)
	if rootVal == nil {
		d.CheckError(fmt.Errorf("Invalid path: %s", sp.String()))
		return
	}
	basePath = sp.Path.Path
	return
}

func runEditStory(cmd *cobra.Command, args []string) error {
	title := args[0]
	key := args[1]
	value := args[2]

	// Edit
	str := "Stories::" + title + ".value " + key + " " + value
	sp, err := spec.ForPath(str)
	d.PanicIfError(err)

	rootVal, basePath := splitPath(sp)
	path := applyStructEditsSp(sp, rootVal, basePath, args)

	// Commit
	cfg := config.NewResolver() //config default db "Stories"
	fmt.Fprintf(os.Stdout, "%s\n", title)
	db, ds, err := cfg.GetDataset("::" + title)
	d.PanicIfError(err)
	defer db.Close()

	fmt.Fprintf(os.Stdout, "%s\n", path)
	absPath, err := spec.NewAbsolutePath(path)
	fmt.Fprintf(os.Stdout, "%s\n", absPath.String())
	d.CheckError(err)

	valPath := absPath.Resolve(db)
	if valPath == nil {
		d.CheckErrorNoUsage(errors.New(fmt.Sprintf("Error resolving value: %s", path)))
	}

	oldCommitRef, oldCommitExists := ds.MaybeHeadRef()
	if oldCommitExists {
		head := ds.HeadValue()
		if head.Hash() == valPath.Hash() {
			fmt.Fprintf(os.Stdout, "Commit aborted - allow-dupe is set to off and this commit would create a duplicate\n")
			return nil
		}
	}

	meta, err := spec.CreateCommitMetaStruct(db, "", "set value %s "+key+" in story : "+title, nil, nil)
	d.CheckErrorNoUsage(err)

	ds, err = db.Commit(ds, valPath, datas.CommitOptions{Meta: meta})
	d.CheckErrorNoUsage(err)

	if oldCommitExists {
		fmt.Fprintf(os.Stdout, "New head #%v (was #%v)\n", ds.HeadRef().TargetHash().String(), oldCommitRef.TargetHash().String())
	} else {
		fmt.Fprintf(os.Stdout, "New head #%v\n", ds.HeadRef().TargetHash().String())
	}
	fmt.Printf("%s edited\n", title)

	return nil
}

var editStoryCmd = &cobra.Command{
	Use:   "edit <title> <key> <value>",
	Short: "Edit a story.",
	Args:  cobra.ExactArgs(1),
	RunE:  runEditStory,
}

func init() {
	storyCmd.AddCommand(editStoryCmd)
}
