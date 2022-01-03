package ges

import "time"

type Order struct {
	OrderID    string
	CustomerID string
	Location   string

	OrderStatus OrderStatus

	TimeSubmitted time.Time
	TimePromised  time.Time
	TimeEarliest  time.Time
	TimeLatest    time.Time

	IsAssigned bool
	AssignedTo string
}

type OrderStatus string

const (
	OrderStatusNew       OrderStatus = "NEW"
	OrderStatusAssigned  OrderStatus = "ASSIGNED"
	OrderStatusEnRoute   OrderStatus = "EN_ROUTE"
	OrderStatusArrived   OrderStatus = "ARRIVED"
	OrderStatusDelivered OrderStatus = "DELIVERED"
)
