package main

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/tehbilly/ges"
	"github.com/tehbilly/ges/commands"
	"github.com/tehbilly/ges/events"
	"time"
)

// CreateOrderCommandHandler is the command handler for commands.CreateOrder
type CreateOrderCommandHandler struct {
	eventBus *cqrs.EventBus
}

func (o *CreateOrderCommandHandler) HandlerName() string {
	return "CreateOrderCommandHandler"
}

func (o *CreateOrderCommandHandler) NewCommand() interface{} {
	return &commands.CreateOrder{}
}

func (o *CreateOrderCommandHandler) Handle(ctx context.Context, cmd interface{}) error {
	oc, ok := cmd.(*commands.CreateOrder)
	if !ok {
		return fmt.Errorf("unknown command received: %T", cmd)
	}

	// ... pretend we're doing order creation stuff here, calculating promise times, etc.

	if err := o.eventBus.Publish(ctx, &events.OrderCreated{
		OrderID:       oc.OrderID,
		CustomerID:    oc.CustomerID,
		Location:      oc.Location,
		OrderStatus:   ges.OrderStatusNew,
		TimeSubmitted: oc.TimeSubmitted,
		TimePromised:  oc.TimeSubmitted.Add(15 * time.Minute),
		TimeEarliest:  oc.TimeSubmitted.Add(10 * time.Minute),
		TimeLatest:    oc.TimeSubmitted.Add(20 * time.Minute),
	}); err != nil {
		return err
	}

	return nil
}

// OrderCreatedEventHandler is the handler for OrderCreated events
type OrderCreatedEventHandler struct {
	orderStore *OrderStore
}

func (o *OrderCreatedEventHandler) HandlerName() string {
	return "OrderCreatedEventHandler"
}

func (o *OrderCreatedEventHandler) NewEvent() interface{} {
	return &events.OrderCreated{}
}

func (o *OrderCreatedEventHandler) Handle(ctx context.Context, event interface{}) error {
	oc, ok := event.(*events.OrderCreated)
	if !ok {
		return fmt.Errorf("unknown event type received: %T", event)
	}

	if err := o.orderStore.orderCreated(oc); err != nil {
		return err
	}

	return nil
}
