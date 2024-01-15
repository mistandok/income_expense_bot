package events

import "context"

type Fetcher interface {
	Fetch(ctx context.Context, limit int) ([]Event, error)
}

type Processor interface {
	Process(ctx context.Context, e Event) error
}
