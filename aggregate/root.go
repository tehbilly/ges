package aggregate

import (
	"fmt"
	"github.com/tehbilly/ges"
)

type ID interface {
	fmt.Stringer
}

type StringID string

func (id StringID) String() string {
	return string(id)
}

type Root interface {
	Applier

	AggregateID() ID

	flushEvents() []ges.Event
	recordEvent(applier Applier, events ...ges.Event) error
}

// BaseRoot is embedded into a type, which will complete the contract for the Root interface if the Applier interface
// and AggregateID methods are implemented.
type BaseRoot struct {
	events []ges.Event
}

func (b BaseRoot) flushEvents() []ges.Event {
	events := b.events
	b.events = nil
	return events
}

func (b BaseRoot) recordEvent(applier Applier, events ...ges.Event) error {
	for _, event := range events {
		if err := applier.Apply(event); err != nil {
			return fmt.Errorf("aggregate: failed to apply event: %w", err)
		}

		b.events = append(b.events, event)
	}

	return nil
}
