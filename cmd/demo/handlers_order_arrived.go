package main

import (
	"context"
	"fmt"
	"github.com/tehbilly/ges/events"
)

type OrderArrivedEventHandler struct {
	orderStore *OrderStore
}

func (o *OrderArrivedEventHandler) HandlerName() string {
	return "OrderArrivedEventHandler"
}

func (o *OrderArrivedEventHandler) NewEvent() interface{} {
	return &events.OrderArrived{}
}

func (o *OrderArrivedEventHandler) Handle(ctx context.Context, event interface{}) error {
	oa, ok := event.(*events.OrderArrived)
	if !ok {
		return fmt.Errorf("unknown event type received: %T", event)
	}

	if err := o.orderStore.orderArrived(oa); err != nil {
		return err
	}

	return nil
}
