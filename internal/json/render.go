package json

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/davidscholberg/go-durationfmt"
	"github.com/jmartin82/mkpis/pkg/vcs"
)

type PRList struct {
	PRs []PR `json:"prs"`
}

type PR struct {
	Commits           int    `json:"commits"`
	Size              int    `json:"size"`
	TimeToFirstReview string `json:"timeToFirstReview"`
	ReviewTime        string `json:"reviewTime"`
	LastReviewToMerge string `json:"lastReviewToMerge"`
	Comments          int    `json:"comments"`
	PRLeadTime        string `json:"prLeadTime"`
	TimeToMerge       string `json:"timeToMerge"`
}

func Render(prs []vcs.PR, owner, repo string, from, to time.Time, includeCreator bool) error {
	f, err := os.Create("pr_report.json")
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)

	jsonPRs := make([]PR, len(prs))

	for i, pr := range prs {
		jsonPRs[i] = PR{
			pr.Commits,
			pr.ChangedLines,
			DurationFormater(pr.TimeToFirstReview()),
			DurationFormater(pr.TimeToReview()),
			DurationFormater(pr.LastReviewToMerge()),
			pr.ReviewComments,
			DurationFormater(pr.PRLeadTime()),
			DurationFormater(pr.TimeToMerge()),
		}
	}

	b, err := json.MarshalIndent(PRList{jsonPRs}, "", "  ")
	if err != nil {
		return err
	}

	w.Write(b)
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
	f, err := os.Create(fmt.Sprintf("pr_%d.json", pr.Number))
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)

	jsonPR := PR{
		pr.Commits,
		pr.ChangedLines,
		DurationFormater(pr.TimeToFirstReview()),
		DurationFormater(pr.TimeToReview()),
		DurationFormater(pr.LastReviewToMerge()),
		pr.ReviewComments,
		DurationFormater(pr.PRLeadTime()),
		DurationFormater(pr.TimeToMerge()),
	}

	b, err := json.MarshalIndent(jsonPR, "", "  ")
	if err != nil {
		return err
	}

	w.Write(b)
	w.Flush()
	return nil
}
