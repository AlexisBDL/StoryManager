package util

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/diff"
	"github.com/attic-labs/noms/go/spec"
	"github.com/attic-labs/noms/go/types"
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

func applyStructEdits(db datas.Database, rootVal types.Value, basePath types.Path, args []string) *spec.AbsolutePath {
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
