package analyzer

import "context"

// Analyzer is an interface of an analyzer
type Analyzer interface {
	Analyze(ctx context.Context) (*Result, error)
}
