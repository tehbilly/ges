package commands

import "time"

type CreateOrder struct {
	OrderID       string
	CustomerID    string
	Location      string
	TimeSubmitted time.Time
}
