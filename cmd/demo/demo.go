package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/tehbilly/ges"
	"github.com/tehbilly/ges/aggregate"
	"github.com/tehbilly/ges/store"
)

func main() {
	ctx := context.Background()

	eventStore := store.NewInMemoryStore()
	repo := aggregate.NewRepository(OrderAggregate, eventStore)

	streamID := store.StreamID{
		Type: OrderAggregate.Name(),
		ID:   uuid.NewString(),
	}

	// Create the order
	fmt.Println("writing: OrderCreated")
	if err := eventStore.Write(ctx, streamID, OrderCreated{
		ID:       streamID.ID,
		Location: "00070",
	}); err != nil {
		panic(err)
	}

	// Update the address
	fmt.Println("writing: OrderAddressChanged")
	if err := eventStore.Write(ctx, streamID, OrderAddressChanged{
		NewAddress: "142 Made Up Pl. Townville, ST 12345",
	}); err != nil {
		panic(err)
	}

	// Add some items
	fmt.Println("writing: OrderItemsAdded")
	if err := eventStore.Write(ctx, streamID, OrderItemsAdded{
		Items: []*LineItem{
			{
				ID:          "food:sammich",
				Name:        "Chicken Sammich",
				RetailPrice: 3.75,
				Quantity:    2,
			},
			{
				ID:          "food:fries",
				Name:        "Waffle Fries",
				RetailPrice: 1.99,
				Quantity:    1,
			},
		},
	}); err != nil {
		panic(err)
	}

	// Add some more items
	fmt.Println("writing: OrderItemsAdded")
	if err := eventStore.Write(ctx, streamID, OrderItemsAdded{
		Items: []*LineItem{
			{
				ID:          "drink:sweet-tea",
				Name:        "Sweet Tea",
				RetailPrice: 0.99,
				Quantity:    1,
			},
		},
	}); err != nil {
		panic(err)
	}

	// Get the order from the repo _(created by creating empty instance then applying events
	fmt.Println("requesting Order from repo")
	o, err := repo.Get(ctx, aggregate.StringID(streamID.ID))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Got order from repo: %s\n", jsonPretty(o))
}

// OrderAggregate is used for aggregate.Root instances to create new instances of Order and the like
var OrderAggregate = aggregate.NewType("order", func() aggregate.Root {
	return new(Order)
})

type LineItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	RetailPrice float32 `json:"retail_price"`
	Quantity    int     `json:"quantity"`
}

type Order struct {
	aggregate.BaseRoot

	// ID is the CFA Order ID
	ID string `json:"id"`

	Location        string      `json:"location"`
	DeliveryAddress string      `json:"delivery_address"`
	LineItems       []*LineItem `json:"line_items"`
}

func (o *Order) Apply(event ges.Event) error {
	switch e := event.(type) {
	case OrderCreated:
		fmt.Printf("Order.Apply(OrderCreated) with event: %s\n", jsonify(e))
		o.ID = e.ID
		o.Location = e.Location
		o.DeliveryAddress = e.DeliveryAddress
	case OrderItemsAdded:
		fmt.Printf("Order.Apply(OrderItemsAdded) with event: %s\n", jsonify(e))
		o.LineItems = append(o.LineItems, e.Items...)
	case OrderAddressChanged:
		fmt.Printf("Order.Apply(OrderAddressChanged) with event: %s\n", jsonify(e))
		o.DeliveryAddress = e.NewAddress
	default:
		return fmt.Errorf("unknown event type %T: %#v", e, e)
	}
	return nil
}

func (o *Order) AggregateID() aggregate.ID {
	return aggregate.StringID(o.ID)
}

func jsonify(o interface{}) string {
	bytes, _ := json.Marshal(o)
	return string(bytes)
}

func jsonPretty(o interface{}) string {
	bytes, _ := json.MarshalIndent(o, "", "  ")
	return string(bytes)
}
