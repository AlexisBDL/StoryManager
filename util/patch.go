package util

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"

	"github.com/AlexisBDL/StoryManager/spec"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/diff"
	"github.com/attic-labs/noms/go/hash"
	"github.com/attic-labs/noms/go/types"
)

const (
	ParentsField = "parents"
	ValueField   = "value"
	MetaField    = "meta"
	commitName   = "Commit"
)

func applyListInserts(sp spec.Spec, rootVal types.Value, basePath types.Path, pos uint64, args []string) *spec.AbsolutePath {
	if rootVal == nil {
		d.CheckErrorNoUsage(fmt.Errorf("No value at: %s", sp.String()))
	}
	db := sp.GetDatabase()
	patch := diff.Patch{}
	for i := 0; i < len(args); i++ {
		vv, err := argumentToValue(args[i], db)
		if err != nil {
			d.CheckError(fmt.Errorf("Invalid value: %s at position %d: %s", args[i], i, err))
		}
		patch = append(patch, diff.Difference{
			Path:       append(basePath, types.NewIndexPath(types.Number(pos+uint64(i)))),
			ChangeType: types.DiffChangeAdded,
			NewValue:   vv,
		})
	}
	return appplyPatch(db, rootVal, basePath, patch)
}

func ApplyStructEdits(db datas.Database, rootVal types.Value, basePath types.Path, args []string) *spec.AbsolutePath {
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
	return appplyPatch(db, rootVal, basePath, patch)
}

func appplyPatch(db datas.Database, rootVal types.Value, basePath types.Path, patch diff.Patch) *spec.AbsolutePath {
	newRootVal := diff.Apply(rootVal, patch)
	d.Chk.NotNil(newRootVal)
	r := db.WriteValue(newRootVal)
	db.Flush()
	newAbsPath := spec.AbsolutePath{
		Hash: r.TargetHash(),
		Path: basePath,
	}
	return &newAbsPath
}

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

func parseKeyPart(args []string, i int) (res types.PathPart) {
	idx, h, rem, err := types.ParsePathIndex(args[i])
	if rem != "" {
		d.CheckError(fmt.Errorf("Invalid key: %s at position %d", args[i], i))
	}
	if err != nil {
		d.CheckError(fmt.Errorf("Invalid key: %s at position %d: %s", args[i], i, err))
	}
	if idx != nil {
		res = types.NewIndexPath(idx)
	} else {
		res = types.NewHashIndexPath(h)
	}
	return
}

// FindCommonAncestor returns the most recent common ancestor of c1 and c2, if
// one exists, setting ok to true. If there is no common ancestor, ok is set
// to false.
func FindCommonAncestor(c1, c2 types.Ref, vr1 types.ValueReader, vr2 types.ValueReader) (a types.Ref, ok bool) {
	if !datas.IsRefOfCommitType(types.TypeOf(c1)) {
		d.Panic("FindCommonAncestor() called on %s", types.TypeOf(c1).Describe())
	}
	if !datas.IsRefOfCommitType(types.TypeOf(c2)) {
		d.Panic("FindCommonAncestor() called on %s", types.TypeOf(c2).Describe())
	}

	c1Q, c2Q := &types.RefByHeight{c1}, &types.RefByHeight{c2}
	for !c1Q.Empty() && !c2Q.Empty() {
		c1Ht, c2Ht := c1Q.MaxHeight(), c2Q.MaxHeight()
		if c1Ht == c2Ht {
			c1Parents, c2Parents := c1Q.PopRefsOfHeight(c1Ht), c2Q.PopRefsOfHeight(c2Ht)
			if common, ok := findCommonRef(c1Parents, c2Parents); ok {
				return common, true
			}
			parentsToQueue(c1Parents, c1Q, vr1)
			parentsToQueue(c2Parents, c2Q, vr2)
		} else if c1Ht > c2Ht {
			parentsToQueue(c1Q.PopRefsOfHeight(c1Ht), c1Q, vr1)
		} else {
			parentsToQueue(c2Q.PopRefsOfHeight(c2Ht), c2Q, vr2)
		}
	}
	return
}

func parentsToQueue(refs types.RefSlice, q *types.RefByHeight, vr types.ValueReader) {
	for _, r := range refs {
		c := r.TargetValue(vr).(types.Struct)
		p := c.Get(ParentsField).(types.Set)
		p.IterAll(func(v types.Value) {
			q.PushBack(v.(types.Ref))
		})
	}
	sort.Sort(q)
}

func findCommonRef(a, b types.RefSlice) (types.Ref, bool) {
	toRefSet := func(s types.RefSlice) map[hash.Hash]types.Ref {
		out := map[hash.Hash]types.Ref{}
		for _, r := range s {
			out[r.TargetHash()] = r
		}
		return out
	}

	aSet, bSet := toRefSet(a), toRefSet(b)
	for s, r := range aSet {
		if _, present := bSet[s]; present {
			return r, true
		}
	}
	return types.Ref{}, false
}
