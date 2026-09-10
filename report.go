package main

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"text/template"
	"time"
)

// ranking is one row of the summary table: a backend reduced to two numbers.
type ranking struct {
	Name          string
	Language      string
	Description   string
	MedianLatency time.Duration
	TotalRPS      float64
	Errors        int
}

// reportData is everything a single-profile report needs.
type reportData struct {
	GeneratedAt string
	Host        string
	Docker      string
	Config      Config
	Profile     string
	ProfileName string

	Ranking   []ranking
	Endpoints []string
	Results   []Result
}

// writeReport renders one profile's report.
func writeReport(path string, cfg Config, profile string, results []Result) error {
	data := reportData{
		GeneratedAt: time.Now().Format("2006-01-02 15:04:05 MST"),
		Host:        hostDescription(),
		Docker:      dockerVersion(),
		Config:      cfg,
		Profile:     profile,
		ProfileName: profileName(profile),
		Endpoints:   endpoints,
		Results:     results,
		Ranking:     rank(results),
	}
	return render(path, reportTemplate, data)
}

// comparisonRow puts one backend's two profiles side by side.
type comparisonRow struct {
	Name     string
	Language string

	DevSetup  string
	ProdSetup string

	DevLatency  time.Duration
	ProdLatency time.Duration
	DevRPS      float64
	ProdRPS     float64

	// LatencySpeedup and ThroughputSpeedup are prod relative to dev: 2.0 means
	// production is twice as fast. Zero means one of the two never ran.
	LatencySpeedup    float64
	ThroughputSpeedup float64
}

// comparisonData is everything the dev-versus-prod report needs.
type comparisonData struct {
	GeneratedAt string
	Host        string
	Docker      string
	Config      Config

	Rows      []comparisonRow
	Endpoints []string
	Dev       []Result
	Prod      []Result
}

// writeComparison renders the report that sets the two profiles against each
// other, backend by backend.
func writeComparison(path string, cfg Config, dev, prod []Result) error {
	devRank := byName(rank(dev))
	prodRank := byName(rank(prod))

	var rows []comparisonRow
	for i, b := range backendsIn(dev) {
		row := comparisonRow{Name: b.Name, Language: b.Language}
		if i < len(dev) {
			row.DevSetup = dev[i].Deployment.Description
		}
		if i < len(prod) {
			row.ProdSetup = prod[i].Deployment.Description
		}

		d, hasDev := devRank[b.Name]
		p, hasProd := prodRank[b.Name]
		if hasDev {
			row.DevLatency, row.DevRPS = d.MedianLatency, d.TotalRPS
		}
		if hasProd {
			row.ProdLatency, row.ProdRPS = p.MedianLatency, p.TotalRPS
		}
		if row.DevLatency > 0 && row.ProdLatency > 0 {
			row.LatencySpeedup = float64(row.DevLatency) / float64(row.ProdLatency)
		}
		if row.DevRPS > 0 && row.ProdRPS > 0 {
			row.ThroughputSpeedup = row.ProdRPS / row.DevRPS
		}
		rows = append(rows, row)
	}

	// Biggest production win first: that is the interesting end of the table.
	sort.SliceStable(rows, func(i, j int) bool {
		return rows[i].ThroughputSpeedup > rows[j].ThroughputSpeedup
	})

	data := comparisonData{
		GeneratedAt: time.Now().Format("2006-01-02 15:04:05 MST"),
		Host:        hostDescription(),
		Docker:      dockerVersion(),
		Config:      cfg,
		Rows:        rows,
		Endpoints:   endpoints,
		Dev:         dev,
		Prod:        prod,
	}
	return render(path, comparisonTemplate, data)
}

func render(path, tmplText string, data any) error {
	tmpl, err := template.New("report").Funcs(reportFuncs).Parse(tmplText)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// rank reduces each backend to a median-of-medians latency and a summed
// throughput, then sorts fastest first.
func rank(results []Result) []ranking {
	var rows []ranking

	for _, r := range results {
		row := ranking{Name: r.Backend.Name, Language: r.Backend.Language, Description: r.Deployment.Description}
		if r.Failed {
			rows = append(rows, row)
			continue
		}

		medians := make([]time.Duration, 0, len(r.Endpoints))
		for _, er := range r.Endpoints {
			medians = append(medians, er.Sequential.Median)
			row.TotalRPS += er.Concurrent.RPS
			row.Errors += er.Sequential.Errors + er.Concurrent.Errors
		}
		row.MedianLatency = medianOf(medians)
		rows = append(rows, row)
	}

	sort.SliceStable(rows, func(i, j int) bool {
		// Backends that never ran sort last.
		if (rows[i].MedianLatency == 0) != (rows[j].MedianLatency == 0) {
			return rows[j].MedianLatency == 0
		}
		return rows[i].MedianLatency < rows[j].MedianLatency
	})
	return rows
}

func byName(rows []ranking) map[string]ranking {
	out := make(map[string]ranking, len(rows))
	for _, r := range rows {
		out[r.Name] = r
	}
	return out
}

// backendsIn lists the backends a result set covers, in the order they ran.
func backendsIn(results []Result) []Backend {
	out := make([]Backend, 0, len(results))
	for _, r := range results {
		out = append(out, r.Backend)
	}
	return out
}

func hostDescription() string {
	return fmt.Sprintf("%s/%s, %d CPUs", runtime.GOOS, runtime.GOARCH, runtime.NumCPU())
}

func profileName(profile string) string {
	if profile == ProfileProd {
		return "production"
	}
	return "development"
}

// findEndpoint locates one endpoint's result within a backend's results.
func findEndpoint(r Result, endpoint string) (EndpointResult, bool) {
	for _, er := range r.Endpoints {
		if er.Endpoint == endpoint {
			return er, true
		}
	}
	return EndpointResult{}, false
}

// formatDuration prints a duration with a consistent unit and precision, since
// Go's default switches between microseconds and milliseconds and makes columns
// hard to compare.
func formatDuration(d time.Duration) string {
	if d == 0 {
		return "-"
	}
	return fmt.Sprintf("%.2f ms", float64(d.Microseconds())/1000)
}

func formatRPS(v float64) string {
	if v == 0 {
		return "-"
	}
	return fmt.Sprintf("%.1f", v)
}

// formatSpeedup states a ratio as a multiplier, and marks a regression rather
// than printing a confusing fraction.
func formatSpeedup(v float64) string {
	switch {
	case v == 0:
		return "-"
	case v >= 1:
		return fmt.Sprintf("%.2fx faster", v)
	default:
		return fmt.Sprintf("%.2fx slower", 1/v)
	}
}

var reportFuncs = template.FuncMap{
	"ms":      formatDuration,
	"rps":     formatRPS,
	"speedup": formatSpeedup,
	"bytes": func(n int64) string {
		if n < 1024 {
			return fmt.Sprintf("%d B", n)
		}
		return fmt.Sprintf("%.1f KB", float64(n)/1024)
	},
	"inc": func(i int) int { return i + 1 },
	"cell": func(r Result, endpoint string) string {
		er, ok := findEndpoint(r, endpoint)
		if !ok {
			return "-"
		}
		return formatDuration(er.Sequential.Median)
	},
	"rpsCell": func(r Result, endpoint string) string {
		er, ok := findEndpoint(r, endpoint)
		if !ok {
			return "-"
		}
		return formatRPS(er.Concurrent.RPS)
	},
	// endpointSpeedup compares one endpoint's throughput across two result
	// sets, matching backends by position.
	"endpointSpeedup": func(dev, prod []Result, index int, endpoint string) string {
		if index >= len(dev) || index >= len(prod) {
			return "-"
		}
		d, okDev := findEndpoint(dev[index], endpoint)
		p, okProd := findEndpoint(prod[index], endpoint)
		if !okDev || !okProd || d.Concurrent.RPS == 0 || p.Concurrent.RPS == 0 {
			return "-"
		}
		return formatSpeedup(p.Concurrent.RPS / d.Concurrent.RPS)
	},
	"join": strings.Join,
}
