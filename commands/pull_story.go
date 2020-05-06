package commands

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/attic-labs/noms/go/config"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/types"
	"github.com/attic-labs/noms/go/util/profile"
	"github.com/attic-labs/noms/go/util/status"
	humanize "github.com/dustin/go-humanize"

	"github.com/spf13/cobra"
)

// creer package util pour plus d'organisation
func runPullStory(cmd *cobra.Command, args []string) error {
	cfg := config.NewResolver()
	title := args[0]
	src := args[1]
	myStore, myObj, err := cfg.GetPath(title)
	d.CheckError(err)
	defer myStore.Close()

	if myObj == nil {
		d.CheckErrorNoUsage(fmt.Errorf("Object %s not found in my Stories", title))
	}

	dbLocal, valueLocal, err := cfg.GetPath(title + ".meta.date")
	d.PanicIfError(err)
	defer dbLocal.Close()
	dbSrc, valueSrc, err := cfg.GetPath(src + "::" + title + ".meta.date")
	d.PanicIfError(err)
	defer dbSrc.Close()

	layout := time.RFC3339
	timeValueLocal, err := strconv.Unquote(types.EncodedValue(valueLocal))
	d.PanicIfError(err)
	timeValueSrc, err := strconv.Unquote(types.EncodedValue(valueSrc))
	d.PanicIfError(err)

	tLocal, err := time.Parse(layout, timeValueLocal)
	tSrc, err := time.Parse(layout, timeValueSrc)
	d.PanicIfError(err)

	if tLocal.After(tSrc) {
		d.CheckErrorNoUsage(fmt.Errorf("Your story %s is more recent. No changes", title))
	}

	sinkDB, sinkDataset, err := cfg.GetDataset(src + "::" + title)
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

	sourceRef := types.NewRef(myObj)
	sinkRef, sinkExists := sinkDataset.MaybeHeadRef()
	nonFF := false
	err = d.Try(func() {
		defer profile.MaybeStartProfile().Stop()
		datas.Pull(myStore, sinkDB, sourceRef, progressCh)

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
	if last := <-lastProgressCh; last.DoneCount > 0 {
		status.Printf("Done - Synced %s in %s (%s/s)",
			humanize.Bytes(last.ApproxWrittenBytes), since(start), bytesPerSec(last.ApproxWrittenBytes, start))
		status.Done()
	} else if !sinkExists {
		fmt.Printf("All chunks already exist at destination! Created new dataset in.\n")
	} else if nonFF && !sourceRef.Equals(sinkRef) {
		fmt.Printf("Abandoning %s; new head is %s\n", sinkRef.TargetHash(), sourceRef.TargetHash())
	} else {
		fmt.Printf("Story is already up to date in my Stories.\n")
	}

	return nil
}

var pullStoryCmd = &cobra.Command{
	Use:   "pull <title> <source>",
	Short: "Pull the story <title> from the databases <source>.",
	Args:  cobra.ExactArgs(2),
	RunE:  runPullStory,
}

func init() {
	storyCmd.AddCommand(pullStoryCmd)

}
