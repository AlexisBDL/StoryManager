package util

import (
	"errors"
	"fmt"

	"github.com/AlexisBDL/StoryManager/spec"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
)

func StoryEdit(db datas.Database, resolved string, fields []string) *spec.AbsolutePath {
	spS, err := spec.ForPath(resolved)
	d.PanicIfError(err)
	defer spS.Close()

	rootVal, basePath := SplitPath(spS)

	return ApplyStructEdits(db, rootVal, basePath, fields)
}

func MapEdit(db datas.Database, resolved string, fields []string) *spec.AbsolutePath {
	spT, err := spec.ForPath(resolved)
	d.PanicIfError(err)
	defer spT.Close()

	rootVal, basePath := SplitPath(spT)

	return ApplyMapEdits(db, rootVal, basePath, fields)
}

func Commit(db datas.Database, ds datas.Dataset, absPath *spec.AbsolutePath, ID string, msg string, user string, title string) {
	valPath := absPath.Resolve(db)
	if valPath == nil {
		d.CheckErrorNoUsage(errors.New(fmt.Sprintf("Error resolving value: %s", absPath.String())))
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
