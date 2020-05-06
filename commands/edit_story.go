package commands

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/attic-labs/noms/go/config"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/spec"

	"github.com/spf13/cobra"
)

var (
	editTitle       string
	editEffort      int
	editDescription string
)

func runEditStory(cmd *cobra.Command, args []string) error {
	title := args[0]
	cfg := config.NewResolver() //config default db "Stories"

	// Edit
	resolved := cfg.ResolvePathSpec(title) + ".value" //.value acces story
	sp, err := spec.ForPath(resolved)
	d.PanicIfError(err)
	defer sp.Close()

	pinned, ok := sp.Pin()
	if !ok {
		fmt.Fprintf(os.Stderr, "Cannot resolve spec: %s\n", title)
		return nil
	}

	db := pinned.GetDatabase()
	ds := db.GetDataset(title)

	rootVal, basePath := SplitPath(sp)
	var absPath *spec.AbsolutePath
	var change string
	var field []string
	if editDescription != "" {
		change += "description "
		field = append(field, "Description", editDescription)
		fmt.Println("desc")
	}
	if editEffort != -1 {
		change += "effort "
		field = append(field, "Effort", strconv.Itoa(editEffort))
		fmt.Println("eff")
	}
	if editTitle != "" {
		change += "title "
		field = append(field, "Title", editTitle)
		fmt.Println("title")
	}

	fmt.Println(field)
	absPath = ApplyStructEdits(db, rootVal, basePath, field)

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

	meta, err := spec.CreateCommitMetaStruct(db, "", "Edit value "+change+"in story : "+title, nil, nil)
	d.CheckErrorNoUsage(err)

	ds, err = db.Commit(ds, valPath, datas.CommitOptions{Meta: meta})
	d.CheckErrorNoUsage(err)

	if oldCommitExists {
		fmt.Printf("New head #%v (was #%v)\n", ds.HeadRef().TargetHash().String(), oldCommitRef.TargetHash().String())
	} else {
		fmt.Printf("New head #%v\n", ds.HeadRef().TargetHash().String())
	}
	fmt.Printf("%s edited\n", title)

	return nil
}

var editStoryCmd = &cobra.Command{
	Use:   "edit <title> [flag] <value>",
	Short: "Edit a field of story.",
	Args:  cobra.ExactArgs(1),
	RunE:  runEditStory,
}

func init() {
	storyCmd.AddCommand(editStoryCmd)

	editStoryCmd.Flags().IntVarP(&editEffort, "effort", "e", -1,
		"Provide an effort to evaluate the story",
	)
	editStoryCmd.Flags().StringVarP(&editDescription, "description", "d", "",
		"Provide a message to describe the story, use \"\" to add more than one word",
	)
	editStoryCmd.Flags().StringVarP(&editTitle, "title", "t", "",
		"Change the title of the story",
	)
}
