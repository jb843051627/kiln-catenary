package report

import "fmt"

func PlainSummary(s RunSummary) string {
	return fmt.Sprintf("run=%s\nkiln=%s\nstatus=%s\nsamples=%d\nsafe=%d\nscore=%.3f", s.RunID, s.KilnID, s.Status, s.Samples, s.Safe, s.Score)
}
func Headers() []string { return []string{"run", "kiln", "status", "samples", "score"} }
func Row(s RunSummary) []string {
	return []string{s.RunID, s.KilnID, s.Status, fmt.Sprintf("%d", s.Samples), fmt.Sprintf("%.3f", s.Score)}
}
