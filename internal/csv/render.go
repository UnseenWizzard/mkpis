package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/davidscholberg/go-durationfmt"
	"github.com/jmartin82/mkpis/pkg/vcs"
)

func Render(prs []vcs.PR, owner, repo string, from, to time.Time, includeCreator bool) error {
	f, err := os.Create("pr_report.csv")
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	header := []string{"Commits", "Size", "Time To First Review", "Review time", "Last Review To Merge", "Comments", "PR Lead Time", "Time To Merge"}
	err = w.Write(header)
	if err != nil {
		return err
	}
	for _, pr := range prs {
		err = w.Write([]string{
			strconv.Itoa(pr.Commits),
			strconv.Itoa(pr.ChangedLines),
			DurationFormater(pr.TimeToFirstReview()),
			DurationFormater(pr.TimeToReview()),
			DurationFormater(pr.LastReviewToMerge()),
			strconv.Itoa(pr.ReviewComments),
			DurationFormater(pr.PRLeadTime()),
			DurationFormater(pr.TimeToMerge()),
		})
		if err != nil {
			return err
		}
	}

	w.Flush()
	return nil
}

func DurationFormater(d time.Duration) string {

	if d.Microseconds() == 0 {
		return ""
	}

	t, err := durationfmt.Format(d, "%hh %mm")
	if err != nil {
		return "ERROR"
	}
	return t
}

func RenderSingle(pr vcs.PR) error {
	f, err := os.Create(fmt.Sprintf("pr_%d.csv", pr.Number))
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	header := []string{"Commits", "Size", "Time To First Review", "Review time", "Last Review To Merge", "Comments", "PR Lead Time", "Time To Merge"}
	err = w.Write(header)
	if err != nil {
		return err
	}

	err = w.Write([]string{
		strconv.Itoa(pr.Commits),
		strconv.Itoa(pr.ChangedLines),
		DurationFormater(pr.TimeToFirstReview()),
		DurationFormater(pr.TimeToReview()),
		DurationFormater(pr.LastReviewToMerge()),
		strconv.Itoa(pr.ReviewComments),
		DurationFormater(pr.PRLeadTime()),
		DurationFormater(pr.TimeToMerge()),
	})
	if err != nil {
		return err
	}

	w.Flush()
	return nil
}
