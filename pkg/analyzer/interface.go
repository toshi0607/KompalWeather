package analyzer

import "context"

type Analyzer interface {
	Analyze(ctx context.Context) (*Result, error)
}
