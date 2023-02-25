package vcs

import (
	"fmt"
	"time"

	"github.com/montanaflynn/stats"
)

func averageDuration(xs []time.Duration) time.Duration {
	var total float64
	n := 0
	for _, v := range xs {
		if v.Nanoseconds() == 0 {
			continue
		}
		total += v.Seconds()
		n++
	}
	t := fmt.Sprintf("%fs", total/float64(n))
	d, _ := time.ParseDuration(t)
	return d
}

type KPICalculator struct {
	prs               []PR
	commits           []float64
	changes           []float64
	reviews           []float64
	timeToMerge       []time.Duration
	timeToReview      []time.Duration
	timeToFirstReview []time.Duration
	lastReviewToMerge []time.Duration
	pRLeadTime        []time.Duration
}

func NewKPICalculator(prs []PR) *KPICalculator {
	kpi := &KPICalculator{
		prs: prs,
	}
	kpi.calc()
	return kpi
}

func (kpi *KPICalculator) calc() {
	for _, pr := range kpi.prs {
		kpi.commits = append(kpi.commits, float64(pr.Commits))
		kpi.changes = append(kpi.changes, float64(pr.ChangedLines))
		kpi.timeToMerge = append(kpi.timeToMerge, pr.TimeToMerge())
		kpi.timeToReview = append(kpi.timeToReview, pr.TimeToReview())
		kpi.timeToFirstReview = append(kpi.timeToFirstReview, pr.TimeToFirstReview())
		kpi.lastReviewToMerge = append(kpi.lastReviewToMerge, pr.LastReviewToMerge())
		kpi.pRLeadTime = append(kpi.pRLeadTime, pr.PRLeadTime())
		kpi.reviews = append(kpi.reviews, float64(pr.ReviewComments))

	}

}

func (kpi *KPICalculator) CountPR() int {
	return len(kpi.prs)
}

func (kpi *KPICalculator) AvgCommits() float64 {
	avg, _ := stats.Mean(kpi.commits)
	return avg
}

func (kpi *KPICalculator) MedianCommits() float64 {
	m, _ := stats.Median(kpi.commits)
	return m
}

func (kpi *KPICalculator) AvgChangedLines() float64 {
	avg, _ := stats.Mean(kpi.changes)
	return avg
}

func (kpi *KPICalculator) MedianChangedLines() float64 {
	m, _ := stats.Median(kpi.changes)
	return m
}

func (kpi *KPICalculator) AvgTimeToMerge() time.Duration {
	return averageDuration(kpi.timeToMerge)
}

func (kpi *KPICalculator) AvgTimeToReview() time.Duration {
	return averageDuration(kpi.timeToReview)
}
func (kpi *KPICalculator) AvgTimeToFirstReview() time.Duration {
	return averageDuration(kpi.timeToFirstReview)
}
func (kpi *KPICalculator) AvgLastReviewToMerge() time.Duration {
	return averageDuration(kpi.lastReviewToMerge)
}
func (kpi *KPICalculator) AvgPRLeadTime() time.Duration {
	return averageDuration(kpi.pRLeadTime)
}

func (kpi *KPICalculator) AvgReviews() float64 {
	avg, _ := stats.Mean(kpi.reviews)
	return avg
}

func (kpi *KPICalculator) MedianReviews() float64 {
	m, _ := stats.Median(kpi.reviews)
	return m
}
