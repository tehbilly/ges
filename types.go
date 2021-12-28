package ges

type Command Message

type Event Message

type Payload interface {
	Name() string
}

type Message interface {
	Payload
}
