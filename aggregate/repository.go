package aggregate

import (
	"context"
	"errors"
	"github.com/tehbilly/ges/store"
)

type Repository struct {
	aggregateType Type
	store         store.Store
}

func NewRepository(aggregateType Type, store store.Store) *Repository {
	return &Repository{
		aggregateType: aggregateType,
		store:         store,
	}
}

func (r Repository) Add(ctx context.Context, root Root) error {
	events := root.flushEvents()
	if len(events) == 0 {
		return nil
	}

	streamID := store.StreamID{
		Type: r.aggregateType.Name(),
		ID:   root.AggregateID().String(),
	}

	if err := r.store.Write(ctx, streamID, events...); err != nil {
		return err
	}

	return nil
}

func (r Repository) Get(ctx context.Context, id ID) (Root, error) {
	instance := r.aggregateType.new()
	empty := true

	all, err := r.store.ReadAll(ctx, store.StreamID{
		Type: r.aggregateType.Name(),
		ID:   id.String(),
	})
	if err != nil {
		return nil, err
	}

	for _, event := range all {
		empty = false
		if err := instance.Apply(event.Event); err != nil {
			return nil, err
		}
	}

	if empty {
		return nil, errors.New("unable to find '" + r.aggregateType.Name() + "' instance with id: " + id.String())
	}

	return instance, nil
}
