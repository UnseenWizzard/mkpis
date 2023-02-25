package vcs

import (
	"time"

	"github.com/montanaflynn/stats"
)

type KPICalculator struct {
	prs               []PR
	commits           []float64
	changes           []float64
	reviews           []float64
	timeToMerge       []float64
	timeToReview      []float64
	timeToFirstReview []float64
	lastReviewToMerge []float64
	pRLeadTime        []float64
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
		kpi.timeToMerge = append(kpi.timeToMerge, float64(pr.TimeToMerge()))
		kpi.timeToReview = append(kpi.timeToReview, float64(pr.TimeToReview()))
		kpi.timeToFirstReview = append(kpi.timeToFirstReview, float64(pr.TimeToFirstReview()))
		kpi.lastReviewToMerge = append(kpi.lastReviewToMerge, float64(pr.LastReviewToMerge()))
		kpi.pRLeadTime = append(kpi.pRLeadTime, float64(pr.PRLeadTime()))
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
	return timeStatsWithoutZeroDurations(kpi.timeToMerge, stats.Mean)
}

func (kpi *KPICalculator) MedianTimeToMerge() time.Duration {
	return timeStatsWithoutZeroDurations(kpi.timeToMerge, stats.Median)
}

func (kpi *KPICalculator) AvgTimeToReview() time.Duration {
	return timeStatsWithoutZeroDurations(kpi.timeToReview, stats.Mean)
}

func (kpi *KPICalculator) MedianTimeToReview() time.Duration {
	return timeStatsWithoutZeroDurations(kpi.timeToReview, stats.Median)
}

func (kpi *KPICalculator) AvgTimeToFirstReview() time.Duration {
	return timeStatsWithoutZeroDurations(kpi.timeToFirstReview, stats.Mean)
}

func (kpi *KPICalculator) MedianTimeToFirstReview() time.Duration {
	return timeStatsWithoutZeroDurations(kpi.timeToFirstReview, stats.Median)
}

func (kpi *KPICalculator) AvgLastReviewToMerge() time.Duration {
	return timeStatsWithoutZeroDurations(kpi.lastReviewToMerge, stats.Mean)
}

func (kpi *KPICalculator) MedianLastReviewToMerge() time.Duration {
	return timeStatsWithoutZeroDurations(kpi.lastReviewToMerge, stats.Median)
}

func (kpi *KPICalculator) AvgPRLeadTime() time.Duration {
	return timeStatsWithoutZeroDurations(kpi.pRLeadTime, stats.Mean)
}

func (kpi *KPICalculator) MedianPRLeadTime() time.Duration {
	return timeStatsWithoutZeroDurations(kpi.pRLeadTime, stats.Median)
}

func (kpi *KPICalculator) AvgReviews() float64 {
	avg, _ := stats.Mean(kpi.reviews)
	return avg
}

func (kpi *KPICalculator) MedianReviews() float64 {
	m, _ := stats.Median(kpi.reviews)
	return m
}

func timeStatsWithoutZeroDurations(durs []float64, statFunc func(stats.Float64Data) (float64, error)) time.Duration {
	filtered := make([]float64, 0, len(durs))
	for _, d := range durs {
		if d > 0.0 {
			filtered = append(filtered, d)
		}
	}
	v, _ := statFunc(filtered)
	return time.Duration(v)
}
