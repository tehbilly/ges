package store

import (
	"context"
	"fmt"
	"github.com/tehbilly/ges"
	"sync"
)

type inMemStore struct {
	m sync.RWMutex

	events            []Event
	byType            map[string][]int
	byTypeAndInstance map[string]map[string][]int
}

func NewInMemoryStore() Store {
	return &inMemStore{
		byType:            make(map[string][]int),
		byTypeAndInstance: make(map[string]map[string][]int),
	}
}

func (s *inMemStore) Write(ctx context.Context, id StreamID, events ...ges.Event) error {
	if len(events) == 0 {
		return nil
	}

	s.m.Lock()
	defer s.m.Unlock()

	s.ensureMaps(id.Type)

	nextOffset := int64(len(s.events))
	newEvents := make([]Event, 0, len(events))
	newIndexes := make([]int, 0, len(events))

	for i, event := range events {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context done: %w", ctx.Err())
		default:
			// Nothing to worry about
		}

		nextIndex := int(nextOffset) + i

		// SequenceNumber should start at 1, that's why there's `+ 1` below
		newEvents = append(newEvents, Event{
			Stream:         id,
			SequenceNumber: int64(nextIndex) + 1,
			Event:          event,
		})
		newIndexes = append(newIndexes, nextIndex)
	}

	s.events = append(s.events, newEvents...)
	s.byType[id.Type] = append(s.byType[id.Type], newIndexes...)
	s.byTypeAndInstance[id.Type][id.ID] = append(s.byTypeAndInstance[id.Type][id.ID], newIndexes...)

	return nil
}

func (s *inMemStore) ReadAll(ctx context.Context, id StreamID) ([]Event, error) {
	if m, ok := s.byTypeAndInstance[id.Type]; !ok || m == nil {
		// Type doesn't exist
		return nil, nil
	}

	idxs, ok := s.byTypeAndInstance[id.Type][id.ID]
	if !ok {
		return nil, nil
	}

	var result []Event

	for _, idx := range idxs {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context is done: %w", ctx.Err())
		default:
			result = append(result, s.events[idx])
		}
	}

	return result, nil
}

func (s *inMemStore) ensureMaps(t string) {
	if v, ok := s.byTypeAndInstance[t]; !ok || v == nil {
		s.byTypeAndInstance[t] = make(map[string][]int)
	}
}
