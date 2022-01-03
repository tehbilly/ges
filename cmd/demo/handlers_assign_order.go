package main

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/tehbilly/ges/commands"
	"github.com/tehbilly/ges/events"
)

// AssignOrderCommandHandler is the command handler for commands.AssignOrder
type AssignOrderCommandHandler struct {
	eventBus *cqrs.EventBus
}

func (a AssignOrderCommandHandler) HandlerName() string {
	return "AssignOrderCommandHandler"
}

func (a AssignOrderCommandHandler) NewCommand() interface{} {
	return &commands.AssignOrder{}
}

func (a AssignOrderCommandHandler) Handle(ctx context.Context, cmd interface{}) error {
	ao, ok := cmd.(*commands.AssignOrder)
	if !ok {
		return fmt.Errorf("unknown command received: %T", cmd)
	}

	// ... pretend we're calling Fleet Engine or other assignment work here

	if err := a.eventBus.Publish(ctx, &events.OrderAssigned{
		OrderID:    ao.OrderID,
		AssignedTo: ao.AssignTo,
	}); err != nil {
		return err
	}

	return nil
}

// OrderAssignedEventHandler is the event handler for OrderAssigned
type OrderAssignedEventHandler struct {
	orderStore *OrderStore
}

func (o *OrderAssignedEventHandler) HandlerName() string {
	return "OrderAssignedEventHandler"
}

func (o *OrderAssignedEventHandler) NewEvent() interface{} {
	return &events.OrderAssigned{}
}

func (o *OrderAssignedEventHandler) Handle(ctx context.Context, event interface{}) error {
	oa, ok := event.(*events.OrderAssigned)
	if !ok {
		return fmt.Errorf("unknown event type received: %T", event)
	}

	if err := o.orderStore.orderAssigned(oa); err != nil {
		return err
	}

	return nil
}
