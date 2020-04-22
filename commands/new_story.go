package commands

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/diff"
	"github.com/attic-labs/noms/go/spec"
	"github.com/attic-labs/noms/go/types"
	"github.com/spf13/cobra"
)

func argumentToValue(arg string, db datas.Database) (types.Value, error) {
	if arg == "" {
		return types.String(""), nil
	}
	if arg == "true" {
		return types.Bool(true), nil
	}
	if arg == "false" {
		return types.Bool(false), nil
	}
	if arg[0] == '"' {
		buf := bytes.Buffer{}
		for i := 1; i < len(arg); i++ {
			c := arg[i]
			if c == '"' {
				if i != len(arg)-1 {
					break
				}
				return types.String(buf.String()), nil
			}
			if c == '\\' {
				i++
				c = arg[i]
				if c != '\\' && c != '"' {
					return nil, fmt.Errorf("Invalid string argument: %s: Only '\\' and '\"' can be escaped", arg)
				}
			}
			buf.WriteByte(c)
		}
		return nil, fmt.Errorf("Invalid string argument: %s", arg)
	}
	if arg[0] == '@' {
		p, err := spec.NewAbsolutePath(arg[1:])
		d.PanicIfError(err)
		return p.Resolve(db), nil
	}
	if n, err := strconv.ParseFloat(arg, 64); err == nil {
		return types.Number(n), nil
	}

	return types.String(arg), nil
}

func applyStructEdits(sp spec.Spec, rootVal types.Value, basePath types.Path, args []string) {
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

func appplyPatch(sp spec.Spec, rootVal types.Value, basePath types.Path, patch diff.Patch) {
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
	fmt.Println(newAbsPath.String())
}

func runCreateStory(cmd *cobra.Command, args []string) error {
	sp, err := spec.ForDatabase("Stories")
	d.PanicIfError(err)
	var composition = []string{"description", " ", "effort", "0"}
	applyStructEdits(sp, types.NewStruct(args[0], nil), nil, composition)
	fmt.Printf("%s created\n", args[0])

	return nil
}

var createStoryCmd = &cobra.Command{
	Use:   "create <title>",
	Short: "Create a new story.",
	Args:  cobra.MinimumNArgs(1),
	//PreRunE: loadRepoEnsureUser,
	RunE: runCreateStory,
}

func init() {
	RootCmd.AddCommand(createStoryCmd)

	createStoryCmd.Flags().SortFlags = false
}
