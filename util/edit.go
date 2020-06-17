package util

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/AlexisBDL/StoryManager/config"
	"github.com/AlexisBDL/StoryManager/spec"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/diff"
	"github.com/attic-labs/noms/go/merge"
	"github.com/attic-labs/noms/go/types"
	"github.com/attic-labs/noms/go/util/profile"
	"github.com/attic-labs/noms/go/util/status"
	"github.com/attic-labs/noms/go/util/verbose"
	"github.com/dustin/go-humanize"
)

func ListInsert(db datas.Database, specStr string, pos int, fields []string) *spec.AbsolutePath {
	sp, err := spec.ForPath(specStr)
	d.PanicIfError(err)
	rootVal, basePath := splitPath(sp)
	return applyListInserts(sp, rootVal, basePath, uint64(pos), fields)
}

func ListDel(db datas.Database, specStr string, pos int) *spec.AbsolutePath {
	sp, err := spec.ForPath(specStr)
	d.PanicIfError(err)

	rootVal, basePath := splitPath(sp)
	patch := diff.Patch{}
	// TODO: if len-pos is large this will start to become problematic
	for i := pos; i < pos+1; i++ {
		patch = append(patch, diff.Difference{
			Path:       append(basePath, types.NewIndexPath(types.Number(i))),
			ChangeType: types.DiffChangeRemoved,
		})
	}

	return appplyPatch(db, rootVal, basePath, patch)
}

func ListAppend(specStr string, args []string) *spec.AbsolutePath {
	sp, err := spec.ForPath(specStr)
	d.PanicIfError(err)
	rootVal, basePath := splitPath(sp)
	listVal := basePath.Resolve(rootVal, sp.GetDatabase())
	if listVal == nil {
		d.CheckErrorNoUsage(fmt.Errorf("no value at path: %s", specStr))
	}
	if list, ok := listVal.(types.List); ok {
		return applyListInserts(sp, rootVal, basePath, list.Len(), args)
	} else {
		d.CheckErrorNoUsage(fmt.Errorf("value at %s is not list", specStr))
	}
	return nil
}

func ListGet(specStr string, id uint64) types.Value {
	sp, err := spec.ForPath(specStr)
	d.PanicIfError(err)
	rootVal, basePath := splitPath(sp)
	listVal := basePath.Resolve(rootVal, sp.GetDatabase())
	if listVal == nil {
		d.CheckErrorNoUsage(fmt.Errorf("no value at path: %s", specStr))
	}
	if list, ok := listVal.(types.List); ok {
		return list.Get(id)
	} else {
		d.CheckErrorNoUsage(fmt.Errorf("value at %s is not list", specStr))
	}
	return types.EmptyStruct
}

func ListGetBy(specStr string, field string, value string) {
	sp, err := spec.ForPath(specStr)
	d.PanicIfError(err)
	rootVal, basePath := splitPath(sp)
	listVal := basePath.Resolve(rootVal, sp.GetDatabase())
	if listVal == nil {
		d.CheckErrorNoUsage(fmt.Errorf("no value at path: %s", specStr))
	}
	if list, ok := listVal.(types.List); ok {
		list.IterAll(func(v types.Value, idx uint64) {
			val := v.(types.Struct)
			test := val.Get(field)
			valTest, _ := strconv.Unquote(types.EncodedValue(test))
			if valTest == value {
				fmt.Println(types.EncodedValue(val))
			}
		})
	} else {
		d.CheckErrorNoUsage(fmt.Errorf("value at %s is not list", specStr))
	}
}

func ListLen(specStr string) uint64 {
	sp, err := spec.ForPath(specStr)
	d.PanicIfError(err)
	rootVal, basePath := splitPath(sp)
	listVal := basePath.Resolve(rootVal, sp.GetDatabase())
	if listVal == nil {
		d.CheckErrorNoUsage(fmt.Errorf("no value at path: %s", specStr))
	}
	if list, ok := listVal.(types.List); ok {
		return list.Len()
	} else {
		d.CheckErrorNoUsage(fmt.Errorf("value at %s is not list", specStr))
	}
	return 0
}

func StoryEdit(db datas.Database, resolved string, fields []string) *spec.AbsolutePath {
	spS, err := spec.ForPath(resolved)
	d.PanicIfError(err)
	defer spS.Close()

	rootVal, basePath := splitPath(spS)

	return ApplyStructEdits(db, rootVal, basePath, fields)
}

func Commit(db datas.Database, ds datas.Dataset, valPath types.Value, ID string, msg string, user string, title string) {
	if valPath == nil {
		d.CheckErrorNoUsage(errors.New(fmt.Sprintf("Error resolving value path")))
	}

	_, oldCommitExists := ds.MaybeHeadRef()
	if oldCommitExists {
		head := ds.HeadValue()
		if head.Hash() == valPath.Hash() {
			d.CheckErrorNoUsage(errors.New(fmt.Sprintf("Commit aborted - allow-dupe is set to off and this commit would create a duplicate\n")))
		}
	}

	meta, err := spec.CreateCommitMetaStruct(db, "", msg, user, nil, nil)
	d.CheckErrorNoUsage(err)

	ds, err = db.Commit(ds, valPath, datas.CommitOptions{Meta: meta})
	d.CheckErrorNoUsage(err)

	fmt.Printf("%s edited\nID : %s\n", title, ID)
}

func MergeStory(db datas.Database, ds1, ds2, merge, user string) {

	leftDS, rightDS, mergeDS := resolveDatasets(db, ds1, ds2, merge)
	left, right, ancestor := getMergeCandidates(db, db, leftDS, rightDS)
	policy := decidePolicy("p")
	pc := newMergeProgressChan()
	merged, err := policy(left, right, ancestor, db, pc)
	d.CheckErrorNoUsage(err)
	close(pc)

	meta, err := spec.CreateCommitMetaStruct(db, "", "merge", user, nil, nil)
	d.CheckErrorNoUsage(err)

	r := db.WriteValue(datas.NewCommit(merged, types.NewSet(db, leftDS.HeadRef(), rightDS.HeadRef()), meta))
	_, err = db.SetHead(mergeDS, r)
	d.PanicIfError(err)
	db.Flush()
	fmt.Println(r.TargetHash())
}

func checkIfTrue(b bool, format string, args ...interface{}) {
	if b {
		d.CheckErrorNoUsage(fmt.Errorf(format, args...))
	}
}

var (
	datasetRe = regexp.MustCompile("^" + datas.DatasetRe.String() + "$")
)

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

func getMergeCandidates(dbL datas.Database, dbR datas.Database, leftDS, rightDS datas.Dataset) (left, right, ancestor types.Value) {
	leftRef, ok := leftDS.MaybeHeadRef()
	checkIfTrue(!ok, "Dataset %s has no data", leftDS.ID())
	rightRef, ok := rightDS.MaybeHeadRef()
	checkIfTrue(!ok, "Dataset %s has no data", rightDS.ID())
	ancestorCommit, ok := getCommonAncestor(leftRef, rightRef, dbL, dbR)
	checkIfTrue(!ok, "Datasets %s and %s have no common ancestor", leftDS.ID(), rightDS.ID())

	return leftDS.HeadValue(), rightDS.HeadValue(), ancestorCommit.Get(datas.ValueField)
}

func getCommonAncestor(r1, r2 types.Ref, vr1 types.ValueReader, vr2 types.ValueReader) (a types.Struct, found bool) {
	aRef, found := FindCommonAncestor(r1, r2, vr1, vr2)
	if !found {
		return
	}
	v := vr1.ReadValue(aRef.TargetHash())
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

func SyncStory(ID, dsName, dest string, cfg *config.Resolver, chanel bool) bool {
	sourceStore, sourceObj, err := cfg.GetPath(ID)
	d.CheckError(err)
	defer sourceStore.Close()

	if sourceObj == nil {
		d.CheckErrorNoUsage(fmt.Errorf("Story not found in my Stories: %s", ID))
	}

	sinkDB, sinkDataset, err := cfg.GetDataset(dest + "::" + dsName)
	d.CheckError(err)
	defer sinkDB.Close()

	start := time.Now()
	progressCh := make(chan datas.PullProgress)
	lastProgressCh := make(chan datas.PullProgress)

	go func() {
		var last datas.PullProgress

		for info := range progressCh {
			last = info
			if info.KnownCount == 1 {
				// It's better to print "up to date" than "0% (0/1); 100% (1/1)".
				continue
			}

			if status.WillPrint() {
				pct := 100.0 * float64(info.DoneCount) / float64(info.KnownCount)
				status.Printf("Syncing - %.2f%% (%s/s)", pct, bytesPerSec(info.ApproxWrittenBytes, start))
			}
		}
		lastProgressCh <- last
	}()

	sourceRef := types.NewRef(sourceObj)
	sinkRef, sinkExists := sinkDataset.MaybeHeadRef()
	nonFF := false
	err = d.Try(func() {
		defer profile.MaybeStartProfile().Stop()
		datas.Pull(sourceStore, sinkDB, sourceRef, progressCh)
		var err error
		sinkDataset, err = sinkDB.FastForward(sinkDataset, sourceRef)
		if err == datas.ErrMergeNeeded {
			sinkDataset, err = sinkDB.SetHead(sinkDataset, sourceRef)
			nonFF = true
		}
		d.PanicIfError(err)
	})

	if err != nil {
		log.Fatal(err)
	}

	close(progressCh)
	if chanel {
		if last := <-lastProgressCh; last.DoneCount > 0 {
			status.Printf("Done - Synced %s in %s (%s/s)",
				humanize.Bytes(last.ApproxWrittenBytes), since(start), bytesPerSec(last.ApproxWrittenBytes, start))
			status.Done()
		} else if !sinkExists {
			fmt.Printf("All chunks already exist at destination! Created new dataset %s.\n", dest)
			return true
		} else if nonFF && !sourceRef.Equals(sinkRef) {
			fmt.Printf("Abandoning %s; new head is %s\n", sinkRef.TargetHash(), sourceRef.TargetHash())
		} else {
			fmt.Printf("Dataset %s is already up to date.\n", dest)
		}
	}

	return false
}

func bytesPerSec(bytes uint64, start time.Time) string {
	bps := float64(bytes) / float64(time.Since(start).Seconds())
	return humanize.Bytes(uint64(bps))
}

func since(start time.Time) string {
	round := time.Second / 100
	now := time.Now().Round(round)
	return now.Sub(start.Round(round)).String()
}
