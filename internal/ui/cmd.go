package ui

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/davidscholberg/go-durationfmt"
	"github.com/jmartin82/mkpis/pkg/vcs"
	"github.com/olekukonko/tablewriter"
)

func AvgDurationFormater(d time.Duration) string {
	t, err := durationfmt.Format(d, "AVG: %dd %hh %mm")
	if err != nil {
		return "ERROR"
	}
	return t
}

func FullDurationFormater(avg, median time.Duration) string {
	aS, err := durationfmt.Format(avg, "%dd %hh %mm")
	if err != nil {
		aS = "ERROR"
	}
	mS, err := durationfmt.Format(median, "%dd %hh %mm")
	if err != nil {
		mS = "ERROR"
	}
	return fmt.Sprintf("AVG: %s\nMED: %s", aS, mS)
}

func DurationFormater(d time.Duration) string {

	if d.Microseconds() == 0 {
		return "--"
	}

	t, err := durationfmt.Format(d, "%hh %mm")
	if err != nil {
		return "ERROR"
	}
	return t
}

type CmdUI struct {
	client       vcs.Client
	owner        string
	repo         string
	develBranch  string
	masterBranch string
}

func NewCmdUI(client vcs.Client, owner, repo, develBranch, masterBranch string) *CmdUI {
	return &CmdUI{
		client:       client,
		owner:        owner,
		repo:         repo,
		develBranch:  develBranch,
		masterBranch: masterBranch,
	}
}

func (u CmdUI) Render(from, to time.Time, includeCreator bool) error {
	rfb, err := u.getFeatureBranchReport(from, to, includeCreator)
	if err != nil {
		return err
	}
	rrb, err := u.getReleaseBranchReport(from, to)
	if err != nil {
		return err
	}

	myFigure := figure.NewColorFigure("Printing the reports...", "standard", "white", true)
	myFigure.Blink(1000, 300, 300)

	fmt.Println("\033[2J") //clean previous ouput
	u.PrintPageHeader(from, to)
	u.PrintRepotHeader("Feature Branch Report")
	fmt.Println(rfb)
	u.PrintRepotHeader("Release Branch Report")
	fmt.Println(rrb)
	return nil
}

func (u CmdUI) PrintRepotHeader(text string) {
	figure.NewColorFigure(text, "small", "green", true).Print()
	fmt.Println("")
}

func (u CmdUI) PrintPageHeader(from time.Time, to time.Time) {
	figure.NewColorFigure("MKPIS", "standard", "red", true).Print()
	fLayout := "2006-02-01"
	fmt.Printf("\n Repo: %s/%s (%s-%s)", u.owner, u.repo, from.Format(fLayout), to.Format(fLayout))
	fmt.Println("")
}

func (u CmdUI) getFeatureBranchReport(from, to time.Time, includeCreator bool) (string, error) {
	prs, err := u.client.GetMergedPRList(u.owner, u.repo, from, to, u.develBranch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error gathering information: %s", err.Error())
		return "", err
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	header := []string{"PR"}
	if includeCreator {
		header = append(header, "Creator")
	}
	header = append(header, "Commits", "Size", "Time To First Review", "Review time", "Last Review To Merge", "Comments", "PR Lead Time", "Time To Merge")
	table.SetHeader(header)

	for _, pr := range prs {
		row := []string{strconv.Itoa(pr.Number)}
		if includeCreator {
			row = append(row, pr.Creator)
		}
		row = append(row,
			strconv.Itoa(pr.Commits),
			strconv.Itoa(pr.ChangedLines),
			DurationFormater(pr.TimeToFirstReview()),
			DurationFormater(pr.TimeToReview()),
			DurationFormater(pr.LastReviewToMerge()),
			strconv.Itoa(pr.ReviewComments),
			DurationFormater(pr.PRLeadTime()),
			DurationFormater(pr.TimeToMerge()),
		)

		table.Append(row)
	}

	kpi := vcs.NewKPICalculator(prs)
	footer := []string{fmt.Sprintf("Count: %d", kpi.CountPR())}
	if includeCreator {
		footer = append(footer, "-")
	}
	footer = append(footer,
		fmt.Sprintf("AVG: %.2f\nMED: %.2f", kpi.AvgCommits(), kpi.MedianCommits()),
		fmt.Sprintf("AVG: %.2f\nMED: %.2f", kpi.AvgChangedLines(), kpi.MedianChangedLines()),
		FullDurationFormater(kpi.AvgTimeToFirstReview(), kpi.MedianTimeToFirstReview()),
		FullDurationFormater(kpi.AvgTimeToReview(), kpi.MedianTimeToReview()),
		FullDurationFormater(kpi.AvgLastReviewToMerge(), kpi.MedianLastReviewToMerge()),
		fmt.Sprintf("AVG: %.2f\nMED: %.2f", kpi.AvgReviews(), kpi.MedianReviews()),
		FullDurationFormater(kpi.AvgPRLeadTime(), kpi.MedianPRLeadTime()),
		FullDurationFormater(kpi.AvgTimeToMerge(), kpi.MedianTimeToMerge()),
	)

	table.SetFooter(footer)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.Render() // Send output
	return tableString.String(), nil
}

func (u CmdUI) getReleaseBranchReport(from, to time.Time) (string, error) {
	prs, err := u.client.GetMergedPRList(u.owner, u.repo, from, to, u.masterBranch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error gathering information: %s", err.Error())
		return "", err
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"PR", "Commits", "Size", "PR Lead Time", "Time To Merge"})

	for _, pr := range prs {
		table.Append([]string{
			strconv.Itoa(pr.Number),
			strconv.Itoa(pr.Commits),
			strconv.Itoa(pr.ChangedLines),
			DurationFormater(pr.PRLeadTime()),
			DurationFormater(pr.TimeToMerge()),
		})

	}

	kpi := vcs.NewKPICalculator(prs)

	table.SetFooter([]string{
		fmt.Sprintf("Count: %d", kpi.CountPR()),
		fmt.Sprintf("AVG: %.2f", kpi.AvgCommits()),
		fmt.Sprintf("AVG: %.2f", kpi.AvgChangedLines()),
		AvgDurationFormater(kpi.AvgPRLeadTime()),
		AvgDurationFormater(kpi.AvgTimeToMerge()),
	}) // Add Footer
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.Render() // Send output
	return tableString.String(), nil
}
