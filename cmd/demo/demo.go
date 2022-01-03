package main

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/tehbilly/ges"
	"net/http"
)

func main() {
	// Our types
	orderStore := &OrderStore{
		orders: map[string]*ges.Order{},
	}

	logger := watermill.NewStdLogger(false, false)

	// Both a publisher and a subscriber _(handy)_
	commandBus := gochannel.NewGoChannel(gochannel.Config{
		Persistent:                     true,
		BlockPublishUntilSubscriberAck: true,
	}, logger)

	// Both a publisher and a subscriber _(still handy)_
	eventBus := gochannel.NewGoChannel(gochannel.Config{}, logger)

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	// Handle panics in cmd/event handlers
	router.AddMiddleware(middleware.Recoverer)

	facade, err := cqrs.NewFacade(cqrs.FacadeConfig{
		GenerateCommandsTopic: func(commandName string) string {
			return commandName
		},
		CommandHandlers: func(commandBus *cqrs.CommandBus, eventBus *cqrs.EventBus) []cqrs.CommandHandler {
			return []cqrs.CommandHandler{
				&CreateOrderCommandHandler{eventBus: eventBus},
				&AssignOrderCommandHandler{eventBus: eventBus},
			}
		},
		CommandsPublisher: commandBus,
		CommandsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			return commandBus, nil
		},
		GenerateEventsTopic: func(eventName string) string {
			return eventName
		},
		EventHandlers: func(commandBus *cqrs.CommandBus, eventBus *cqrs.EventBus) []cqrs.EventHandler {
			return []cqrs.EventHandler{
				&OrderCreatedEventHandler{orderStore: orderStore},
				&OrderAssignedEventHandler{orderStore: orderStore},
				&OrderLeftEventHandler{orderStore: orderStore},
				&OrderArrivedEventHandler{orderStore: orderStore},
				&OrderDeliveredEventHandler{orderStore: orderStore},
			}
		},
		EventsPublisher: eventBus,
		EventsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			return eventBus, nil
		},
		Router:                router,
		CommandEventMarshaler: cqrs.JSONMarshaler{},
		Logger:                logger,
	})
	if err != nil {
		panic(err)
	}

	go func() {
		if err := router.Run(context.Background()); err != nil {
			panic(err)
		}
	}()

	// Wait for router to start
	<-router.Running()
	fmt.Println("Router is now running. Starting the HTTP handler.")

	handler := newHTTPHandler(orderStore, facade.CommandBus(), facade.EventBus())
	orderStore.handler = handler.eventHandler

	if err := http.ListenAndServe(":8888", handler); err != nil {
		panic(err)
	}
}
