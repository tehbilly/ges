package main

import (
	"context"
	"fmt"
	"github.com/tehbilly/ges/events"
)

type OrderDeliveredEventHandler struct {
	orderStore *OrderStore
}

func (o *OrderDeliveredEventHandler) HandlerName() string {
	return "OrderDeliveredEventHandler"
}

func (o *OrderDeliveredEventHandler) NewEvent() interface{} {
	return &events.OrderDelivered{}
}

func (o *OrderDeliveredEventHandler) Handle(ctx context.Context, event interface{}) error {
	od, ok := event.(*events.OrderDelivered)
	if !ok {
		return fmt.Errorf("unknown event type received: %T", event)
	}

	if err := o.orderStore.orderDelivered(od); err != nil {
		return err
	}

	return nil
}
