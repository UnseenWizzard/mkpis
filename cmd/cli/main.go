package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmartin82/mkpis/internal/config"
	"github.com/jmartin82/mkpis/internal/csv"
	"github.com/jmartin82/mkpis/internal/json"
	"github.com/jmartin82/mkpis/internal/ui"

	"github.com/jmartin82/mkpis/pkg/vcs"
	"github.com/jmartin82/mkpis/pkg/vcs/ghapi"
)

type renderer struct {
	renderSingle func(pr vcs.PR) error
	render       func(prs []vcs.PR, owner, repo string, from, to time.Time, includeCreator bool) error
}

func printError(err string) {
	fmt.Fprintf(os.Stderr, "Error: %s\n\n", err)

	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()

}

func main() {

	//default time range
	windowTime := 10 //days
	today := time.Now()
	nlw := today.AddDate(0, 0, -windowTime)
	tLayout := "2006-01-02"

	log.Println("Starting MKPIS Appplication")
	owner := flag.String("owner", "", "Owner of the repository")
	repo := flag.String("repo", "", "Repository name")
	base := flag.String("base", "master", "Base branch to check for PRs")
	pr := flag.Int("pr", -1, "Single PR to query. If set 'to'/'from' are ignored and single PR is fetched.")
	sfrom := flag.String("from", nlw.Format("2006-01-02"), "When the extraction starts")
	sto := flag.String("to", today.Format("2006-01-02"), "When the extraction ends")
	includeCreator := flag.Bool("include-creator", false, "If set, information about who created a PR is included")
	csv := flag.Bool("csv", false, "If set, output export as csv")
	json := flag.Bool("json", false, "If set, output export as json")
	flag.Parse()

	if len(os.Args) < 2 {
		printError("Invalid number of arguments")
		os.Exit(1)
	}

	if *owner == "" {
		printError("Invalid owner")
		os.Exit(1)
	}

	if *repo == "" {
		printError("Invalid repo")
		os.Exit(1)
	}

	from, err := time.Parse(tLayout, *sfrom)
	if err != nil {
		printError("Invalid `from` date")
		os.Exit(2)
	}

	to, err := time.Parse(tLayout, *sto)
	if err != nil {
		printError("Invalid `to` date")
		os.Exit(2)
	}

	if to.Before(from) {
		printError("`from` date is bigger than `to` date")
		os.Exit(2)
	}

	if config.Env.GitHubToken == "" {
		fmt.Fprintf(os.Stderr, "Error: GITHUB_TOKEN environment variable not found. (You can use .env file to define it)")
		os.Exit(3)
	}

	renderers := setupRenderers(*csv, *json)

	vchClient := ghapi.NewClient(config.Env.GitHubToken)

	if *pr > 0 {
		err = getSingle(*vchClient, *owner, *repo, *pr, renderers)
	} else {
		err = getAll(*vchClient, *owner, *repo, *base, from, to, *includeCreator, renderers)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering: %s\n", err.Error())
		os.Exit(4)
	}
	os.Exit(0)
}

func setupRenderers(renderCSV, renderJSON bool) []renderer {
	var renderers = []renderer{
		{
			ui.RenderSingle,
			ui.Render,
		},
	}
	if renderCSV {
		renderers = append(renderers,
			renderer{
				csv.RenderSingle,
				csv.Render,
			})
	}
	if renderJSON {
		renderers = append(renderers,
			renderer{
				json.RenderSingle,
				json.Render,
			})
	}
	return renderers
}

func getAll(client ghapi.Client, owner, repo, base string, from, to time.Time, includeCreator bool, renderers []renderer) error {
	prs, err := client.GetMergedPRList(owner, repo, from, to, base)
	if err != nil {
		return err
	}
	for _, r := range renderers {
		err = r.render(prs, owner, repo, from, to, includeCreator)
		if err != nil {
			return err
		}
	}
	return nil
}

func getSingle(client ghapi.Client, owner, repo string, prNum int, renderers []renderer) error {
	pr, err := client.GetPRInfo(owner, repo, prNum)
	if err != nil {
		return err
	}
	for _, r := range renderers {
		err = r.renderSingle(pr)
		if err != nil {
			return err
		}
	}
	return nil
}
