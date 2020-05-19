package commands

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/AlexisBDL/StoryManager/config"
	"github.com/AlexisBDL/StoryManager/spec"
	"github.com/AlexisBDL/StoryManager/util"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/types"
	"github.com/attic-labs/noms/go/util/outputpager"

	"github.com/spf13/cobra"
)

const parallelism = 16

func runLogStory(cmd *cobra.Command, args []string) error {
	ID := args[0]               //dataset or value to display history for
	cfg := config.NewResolver() //config default db "Stories"

	o := util.NewOpts(ID)

	resolved := cfg.ResolvePathSpec(ID)
	sp, err := spec.ForPath(resolved)
	d.CheckErrorNoUsage(err)
	defer sp.Close()

	pinned, ok := sp.Pin()
	if !ok {
		fmt.Fprintf(os.Stderr, "Cannot resolve spec: %s\n", ID)
		return nil
	}
	defer pinned.Close()
	database := pinned.GetDatabase()

	absPath := pinned.Path
	path := absPath.Path
	if len(path) == 0 {
		path = types.MustParsePath(".value")
	}

	origCommit, ok := database.ReadValue(absPath.Hash).(types.Struct)
	if !ok || !datas.IsCommit(origCommit) {
		d.CheckError(fmt.Errorf("%s does not reference a Commit object", path))
	}

	iter := util.NewCommitIterator(database, origCommit)
	displayed := 0

	bytesChan := make(chan chan []byte, parallelism)

	var done = false

	go func() {
		for ln, ok := iter.Next(); !done && ok; ln, ok = iter.Next() {
			ch := make(chan []byte)
			bytesChan <- ch

			go func(ch chan []byte, node util.LogNode) {
				buff := &bytes.Buffer{}
				util.PrintCommit(node, path, buff, database, o)
				ch <- buff.Bytes()
			}(ch, ln)

			displayed++
		}
		//close(bytesChan)
	}()

	pgr := outputpager.Start()
	defer pgr.Stop()

	for ch := range bytesChan {
		commitBuff := <-ch
		_, err := io.Copy(pgr.Writer, bytes.NewReader(commitBuff))
		if err != nil {
			done = true
			for range bytesChan {
				// drain the output
			}
		}
	}

	return nil

}

var logStoryCmd = &cobra.Command{
	Use:   "log <title>",
	Short: "Log of a story",
	Args:  cobra.ExactArgs(1),
	RunE:  runLogStory,
}

func init() {
	RootCmd.AddCommand(logStoryCmd)
}
