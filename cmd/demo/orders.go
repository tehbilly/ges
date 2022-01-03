package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/tehbilly/ges"
	"github.com/tehbilly/ges/events"
)

type OrderStore struct {
	orders map[string]*ges.Order

	handler func(event string, order *ges.Order)
}

func (s *OrderStore) Get(orderID string) (*ges.Order, bool) {
	order, ok := s.orders[orderID]
	if !ok {
		return nil, false
	}

	// Copy order object to prevent callers making changes
	// NB: This normally wouldn't be needed for this struct as there are no arrays, slices, or pointers
	orderCopy := &ges.Order{
		OrderID:       order.OrderID,
		CustomerID:    order.CustomerID,
		Location:      order.Location,
		OrderStatus:   order.OrderStatus,
		TimeSubmitted: order.TimeSubmitted,
		TimePromised:  order.TimePromised,
		TimeEarliest:  order.TimeEarliest,
		TimeLatest:    order.TimeLatest,
		IsAssigned:    order.IsAssigned,
		AssignedTo:    order.AssignedTo,
	}

	return orderCopy, true
}

func (s *OrderStore) orderCreated(e *events.OrderCreated) error {
	if _, ok := s.orders[e.OrderID]; ok {
		return errors.New("order already exists: " + e.OrderID)
	}

	s.orders[e.OrderID] = &ges.Order{
		OrderID:       e.OrderID,
		CustomerID:    e.CustomerID,
		Location:      e.Location,
		OrderStatus:   e.OrderStatus,
		TimeSubmitted: e.TimeSubmitted,
		TimePromised:  e.TimePromised,
		TimeEarliest:  e.TimeEarliest,
		TimeLatest:    e.TimeLatest,
	}

	printOrderEvent("OrderCreated", s.orders[e.OrderID])
	s.invokeHandler("OrderCreated", s.orders[e.OrderID])
	return nil
}

func (s *OrderStore) orderAssigned(e *events.OrderAssigned) error {
	order, ok := s.orders[e.OrderID]
	if !ok {
		return errors.New("cannot assign unknown order: " + e.OrderID)
	}

	order.IsAssigned = true
	order.AssignedTo = e.AssignedTo
	order.OrderStatus = ges.OrderStatusAssigned

	printOrderEvent("OrderAssigned", map[string]interface{}{
		"OrderID":     order.OrderID,
		"OrderStatus": order.OrderStatus,
		"IsAssigned":  order.IsAssigned,
		"AssignedTo":  order.AssignedTo,
	})
	s.invokeHandler("OrderAssigned", s.orders[e.OrderID])
	return nil
}

func (s *OrderStore) orderLeft(e *events.OrderLeft) error {
	order, ok := s.orders[e.OrderID]
	if !ok {
		return errors.New("got OrderLeft on unknown order: " + e.OrderID)
	}

	order.OrderStatus = ges.OrderStatusEnRoute

	printOrderEvent("OrderLeft", map[string]interface{}{
		"OrderID":     order.OrderID,
		"OrderStatus": order.OrderStatus,
	})
	s.invokeHandler("OrderLeft", s.orders[e.OrderID])
	return nil
}

func (s OrderStore) orderArrived(e *events.OrderArrived) error {
	order, ok := s.orders[e.OrderID]
	if !ok {
		return errors.New("got OrderArrived on unknown order: " + e.OrderID)
	}

	order.OrderStatus = ges.OrderStatusArrived

	printOrderEvent("OrderArrived", map[string]interface{}{
		"OrderID":     order.OrderID,
		"OrderStatus": order.OrderStatus,
	})
	s.invokeHandler("OrderArrived", s.orders[e.OrderID])
	return nil
}

func (s OrderStore) orderDelivered(e *events.OrderDelivered) error {
	order, ok := s.orders[e.OrderID]
	if !ok {
		return errors.New("got OrderArrived on unknown order: " + e.OrderID)
	}

	order.OrderStatus = ges.OrderStatusDelivered

	printOrderEvent("OrderDelivered", map[string]interface{}{
		"OrderID":     order.OrderID,
		"OrderStatus": order.OrderStatus,
	})
	s.invokeHandler("OrderDelivered", s.orders[e.OrderID])
	return nil
}

func (s *OrderStore) invokeHandler(event string, order *ges.Order) {
	if s.handler != nil {
		s.handler(event, order)
	}
}

func printOrderEvent(event string, order interface{}) {
	mb, _ := json.Marshal(order)
	fmt.Printf(
		"[%s.%s] %s\n",
		color.New(color.FgBlue).Sprint("OrderStore"),
		color.New(color.FgCyan).Sprintf("%-14s", event),
		string(mb),
	)
}
