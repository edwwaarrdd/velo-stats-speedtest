package main

import (
	"sort"
	"time"
)

type Stats struct {
	Count  int
	Errors int

	Min    time.Duration
	Mean   time.Duration
	Median time.Duration
	P95    time.Duration
	P99    time.Duration
	Max    time.Duration

	// For the sequential pass RPS is simply the inverse of the mean; for the
	// concurrent pass it is the number that matters.
	Elapsed time.Duration
	RPS     float64
}

// A copy is sorted, so the caller's slice keeps its original order.
func summarise(durations []time.Duration, errors int, elapsed time.Duration) Stats {
	s := Stats{Count: len(durations), Errors: errors, Elapsed: elapsed}
	if len(durations) == 0 {
		return s
	}

	sorted := make([]time.Duration, len(durations))
	copy(sorted, durations)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	var total time.Duration
	for _, d := range sorted {
		total += d
	}

	s.Min = sorted[0]
	s.Max = sorted[len(sorted)-1]
	s.Mean = total / time.Duration(len(sorted))
	s.Median = percentile(sorted, 0.50)
	s.P95 = percentile(sorted, 0.95)
	s.P99 = percentile(sorted, 0.99)

	if elapsed > 0 {
		s.RPS = float64(len(durations)) / elapsed.Seconds()
	}
	return s
}

// Nearest-rank, over an already sorted slice.
func percentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	index := int(p * float64(len(sorted)))
	if index >= len(sorted) {
		index = len(sorted) - 1
	}
	return sorted[index]
}

// A median of already-computed medians is how a backend's overall latency is
// summarised across its endpoints.
func medianOf(values []time.Duration) time.Duration {
	if len(values) == 0 {
		return 0
	}
	sorted := make([]time.Duration, len(values))
	copy(sorted, values)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	return sorted[len(sorted)/2]
}
