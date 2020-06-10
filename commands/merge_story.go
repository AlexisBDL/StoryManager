package commands

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/merge"
	"github.com/attic-labs/noms/go/types"
	"github.com/attic-labs/noms/go/util/status"
	"github.com/attic-labs/noms/go/util/verbose"

	"github.com/spf13/cobra"
)

var (
	datasetRe = regexp.MustCompile("^" + datas.DatasetRe.String() + "$")
)

// Fusion de branch aillant le même parent
func runMergeStory(cmd *cobra.Command, args []string) error {
	db, err := cfg.GetDatabase("")
	d.CheckError(err)
	defer db.Close()

	leftDS, rightDS := resolveDatasets(db, args[0], args[1])
	left, right, ancestor := getMergeCandidates(db, leftDS, rightDS)
	policy := decidePolicy("p")
	pc := newMergeProgressChan()
	merged, err := policy(left, right, ancestor, db, pc)
	d.CheckErrorNoUsage(err)
	close(pc)

	r := db.WriteValue(datas.NewCommit(merged, types.NewSet(db, leftDS.HeadRef(), rightDS.HeadRef()), types.EmptyStruct))
	db.Flush()
	fmt.Println(r.TargetHash())
	return nil
}

var mergeStoryCmd = &cobra.Command{
	Use:   "merge <ID1> <ID2> ",
	Short: "Merge two stories that have similar ref.",
	Args:  cobra.ExactArgs(2),
	RunE:  runMergeStory,
}

func init() {
	storyCmd.AddCommand(mergeStoryCmd)
}

func checkIfTrue(b bool, format string, args ...interface{}) {
	if b {
		d.CheckErrorNoUsage(fmt.Errorf(format, args...))
	}
}

func resolveDatasets(db datas.Database, leftName, rightName string) (leftDS, rightDS datas.Dataset) {
	makeDS := func(dsName string) datas.Dataset {
		if !datasetRe.MatchString(dsName) {
			d.CheckErrorNoUsage(fmt.Errorf("Invalid dataset %s, must match %s", dsName, datas.DatasetRe.String()))
		}
		return db.GetDataset(dsName)
	}
	leftDS = makeDS(leftName)
	rightDS = makeDS(rightName)
	return
}

func getMergeCandidates(db datas.Database, leftDS, rightDS datas.Dataset) (left, right, ancestor types.Value) {
	leftRef, ok := leftDS.MaybeHeadRef()
	checkIfTrue(!ok, "Dataset %s has no data", leftDS.ID())
	rightRef, ok := rightDS.MaybeHeadRef()
	checkIfTrue(!ok, "Dataset %s has no data", rightDS.ID())
	ancestorCommit, ok := getCommonAncestor(leftRef, rightRef, db)
	checkIfTrue(!ok, "Datasets %s and %s have no common ancestor", leftDS.ID(), rightDS.ID())

	return leftDS.HeadValue(), rightDS.HeadValue(), ancestorCommit.Get(datas.ValueField)
}

func getCommonAncestor(r1, r2 types.Ref, vr types.ValueReader) (a types.Struct, found bool) {
	aRef, found := datas.FindCommonAncestor(r1, r2, vr)
	if !found {
		return
	}
	v := vr.ReadValue(aRef.TargetHash())
	if v == nil {
		panic(aRef.TargetHash().String() + " not found")
	}
	if !datas.IsCommit(v) {
		panic("Not a commit: " + types.EncodedValueMaxLines(v, 10) + "  ...")
	}
	return v.(types.Struct), true
}

func newMergeProgressChan() chan struct{} {
	pc := make(chan struct{}, 128)
	go func() {
		count := 0
		for range pc {
			if !verbose.Quiet() {
				count++
				status.Printf("Applied %d changes...", count)
			}
		}
	}()
	return pc
}

func decidePolicy(policy string) merge.Policy {
	var resolve merge.ResolveFunc
	switch policy {
	case "n", "N":
		resolve = merge.None
	case "l", "L":
		resolve = merge.Ours
	case "r", "R":
		resolve = merge.Theirs
	case "p", "P":
		resolve = func(aType, bType types.DiffChangeType, a, b types.Value, path types.Path) (change types.DiffChangeType, merged types.Value, ok bool) {
			return cliResolve(os.Stdin, os.Stdout, aType, bType, a, b, path)
		}
	default:
		d.CheckErrorNoUsage(fmt.Errorf("Unsupported merge policy: %s. Choices are n, l, r and a.", policy))
	}
	return merge.NewThreeWay(resolve)
}

func cliResolve(in io.Reader, out io.Writer, aType, bType types.DiffChangeType, a, b types.Value, path types.Path) (change types.DiffChangeType, merged types.Value, ok bool) {
	stringer := func(v types.Value) (s string, success bool) {
		switch v := v.(type) {
		case types.Bool, types.Number, types.String:
			return fmt.Sprintf("%v", v), true
		}
		return "", false
	}
	left, lOk := stringer(a)
	right, rOk := stringer(b)
	if !lOk || !rOk {
		return change, merged, false
	}

	// TODO: Handle removes as well.
	fmt.Fprintf(out, "\nConflict at: %s\n", path.String())
	fmt.Fprintf(out, "Left:  %s\nRight: %s\n\n", left, right)
	var choice rune
	for {
		fmt.Fprintln(out, "Enter 'l' to accept the left value, 'r' to accept the right value")
		_, err := fmt.Fscanf(in, "%c\n", &choice)
		d.PanicIfError(err)
		switch choice {
		case 'l', 'L':
			return aType, a, true
		case 'r', 'R':
			return bType, b, true
		}
	}
}
