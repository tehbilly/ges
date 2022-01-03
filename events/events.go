package events

import (
	"github.com/tehbilly/ges"
	"time"
)

// OrderCreated is the event for a new order being created
type OrderCreated struct {
	OrderID    string
	CustomerID string
	Location   string

	OrderStatus ges.OrderStatus

	TimeSubmitted time.Time
	TimePromised  time.Time
	TimeEarliest  time.Time
	TimeLatest    time.Time
}

// OrderAssigned is the event for an order being assigned
type OrderAssigned struct {
	OrderID    string
	AssignedTo string
}

// OrderLeft is the event for an order leaving a restaurant
type OrderLeft struct {
	OrderID string
}

// OrderArrived is the event for an order arriving at the destination
type OrderArrived struct {
	OrderID string
}

// OrderDelivered is the event for an order being delivered to the customer
type OrderDelivered struct {
	OrderID string
}

// OrderCancelled is the event for an order that should no longer be delivered
type OrderCancelled struct {
	OrderID string
}
