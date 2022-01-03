package main

import (
	"context"
	"fmt"
	"github.com/tehbilly/ges/events"
)

// OrderLeftEventHandler handles OrderLeft events
type OrderLeftEventHandler struct {
	orderStore *OrderStore
}

func (o *OrderLeftEventHandler) HandlerName() string {
	return "OrderLeftEventHandler"
}

func (o *OrderLeftEventHandler) NewEvent() interface{} {
	return &events.OrderLeft{}
}

func (o *OrderLeftEventHandler) Handle(ctx context.Context, event interface{}) error {
	ol, ok := event.(*events.OrderLeft)
	if !ok {
		return fmt.Errorf("unknown event type received: %T", event)
	}

	if err := o.orderStore.orderLeft(ol); err != nil {
		return err
	}

	return nil
}
