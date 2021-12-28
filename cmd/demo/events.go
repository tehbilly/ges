package main

const (
	EventOrderCreated        = "order:created"
	EventOrderItemsAdded     = "order:line-items:added"
	EventOrderAddressChanged = "order:line-items:added"
)

// New order created
type OrderCreated struct {
	ID              string
	Location        string
	DeliveryAddress string
	LineItems       []*LineItem
}

func (o OrderCreated) Name() string {
	return EventOrderCreated
}

// Add line items
type OrderItemsAdded struct {
	Items []*LineItem
}

func (o OrderItemsAdded) Name() string {
	return EventOrderItemsAdded
}

// Update address
type OrderAddressChanged struct {
	NewAddress string
}

func (o OrderAddressChanged) Name() string {
	return EventOrderAddressChanged
}
