package vcs

import "time"

type Client interface {
	GetMergedPRList(owner string, repo string, from time.Time, to time.Time, base string) ([]PR, error)
	GetPRInfo(owner string, repo string, prNum int) (PR, error)
}

type PR struct {
	Number         int
	Creator        string
	CreatedAt      time.Time
	MergedAt       time.Time
	Base           string
	ChangedFiles   int
	ChangedLines   int
	ReviewComments int
	Commits        int
	Head           string
	FirstCommitAt  time.Time
	LastCommitAt   time.Time
	FirstCommentAt time.Time
	LastCommentAt  time.Time
}

func (pr *PR) PRLeadTime() time.Duration {
	return pr.MergedAt.Sub(pr.CreatedAt)
}

func (pr *PR) TimeToMerge() time.Duration {
	firstCommitToMerge := pr.MergedAt.Sub(pr.FirstCommitAt)
	createToMerge := pr.PRLeadTime()
	if firstCommitToMerge < createToMerge { // commits probably re-written during review
		return createToMerge
	}
	return firstCommitToMerge
}

func (pr *PR) TimeToReview() time.Duration {
	return pr.LastCommentAt.Sub(pr.FirstCommentAt)
}
func (pr *PR) TimeToFirstReview() time.Duration {
	if pr.FirstCommentAt.IsZero() {
		return 0
	}
	return pr.FirstCommentAt.Sub(pr.CreatedAt)
}
func (pr *PR) LastReviewToMerge() time.Duration {
	if pr.LastCommentAt.IsZero() || pr.LastCommentAt.After(pr.MergedAt) {
		return 0
	}
	return pr.MergedAt.Sub(pr.LastCommentAt)
}
