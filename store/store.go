package store

import (
	"context"
	"github.com/tehbilly/ges"
)

type Event struct {
	ges.Event

	// Stream is the stream of Events that this Event belongs to
	Stream StreamID

	// TODO: Add Version

	// SequenceNumber is the index of this event in the Store, used for ordered streaming
	SequenceNumber int64
}

type StreamID struct {
	// Type typically represents an AggregateType
	Type string
	// ID typically references an entity's ID
	ID string
}

type Writer interface {
	Write(ctx context.Context, id StreamID, events ...ges.Event) error
}

type Reader interface {
	//Read(ctx context.Context, id StreamID, events chan<-Event) error

	ReadAll(ctx context.Context, id StreamID) ([]Event, error)
}

type Store interface {
	Writer
	Reader
}
