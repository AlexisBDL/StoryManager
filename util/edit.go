package util

import (
	"errors"
	"fmt"

	"github.com/AlexisBDL/StoryManager/spec"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/diff"
	"github.com/attic-labs/noms/go/types"
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

	oldCommitRef, oldCommitExists := ds.MaybeHeadRef()
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

	if oldCommitExists {
		fmt.Printf("New head #%v (was #%v)\n", ds.HeadRef().TargetHash().String(), oldCommitRef.TargetHash().String())
	} else {
		fmt.Printf("New head #%v\n", ds.HeadRef().TargetHash().String())
	}
	fmt.Printf("%s edited\nID : %s\n", title, ID)
}
