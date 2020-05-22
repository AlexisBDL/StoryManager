package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/types"
	"github.com/attic-labs/noms/go/util/profile"
	"github.com/attic-labs/noms/go/util/status"
	humanize "github.com/dustin/go-humanize"

	"github.com/spf13/cobra"
)

func runSyncStory(cmd *cobra.Command, args []string) error {
	ID := args[0]
	dest := args[1]
	sourceStore, sourceObj, err := cfg.GetPath(ID)
	d.CheckError(err)
	defer sourceStore.Close()

	if sourceObj == nil {
		d.CheckErrorNoUsage(fmt.Errorf("Story not found in my Stories: %s", ID))
	}

	sinkDB, sinkDataset, err := cfg.GetDataset(dest + "::" + ID)
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

	//close(progressCh)
	if last := <-lastProgressCh; last.DoneCount > 0 {
		status.Printf("Done - Synced %s in %s (%s/s)",
			humanize.Bytes(last.ApproxWrittenBytes), since(start), bytesPerSec(last.ApproxWrittenBytes, start))
		status.Done()
	} else if !sinkExists {
		fmt.Printf("All chunks already exist at destination! Created new dataset in %s.\n", dest)
	} else if nonFF && !sourceRef.Equals(sinkRef) {
		fmt.Printf("Abandoning %s; new head is %s\n", sinkRef.TargetHash(), sourceRef.TargetHash())
	} else {
		fmt.Printf("Dataset %s is already up to date.\n", dest)
	}

	return nil
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

var syncStoryCmd = &cobra.Command{
	Use:   "sync <ID> <destination>",
	Short: "Syncronize the story <ID> with the databases <destination>.",
	Args:  cobra.ExactArgs(2),
	RunE:  runSyncStory,
}

func init() {
	storyCmd.AddCommand(syncStoryCmd)

}
