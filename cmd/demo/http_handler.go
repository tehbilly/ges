package main

import (
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/tehbilly/ges"
	"github.com/tehbilly/ges/commands"
	"github.com/tehbilly/ges/events"
	"io"
	"net/http"
	"time"
)

type httpHandler struct {
	orderStore *OrderStore
	cb         *cqrs.CommandBus
	eb         *cqrs.EventBus

	mux       *http.ServeMux
	conns     []*websocket.Conn
	listeners []chan *event
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *httpHandler) createOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "/create only accepts POST", http.StatusBadRequest)
		return
	}
	cors(&w)

	orderID := uuid.NewString()

	// Create an order by sending a command to create one. This will, after successfully creating the order, create an
	// "OrderCreated" event that we can respond to.
	if err := h.cb.Send(r.Context(), &commands.CreateOrder{
		OrderID:       orderID,
		CustomerID:    "cfa-cust-wmcgann",
		Location:      "00070",
		TimeSubmitted: time.Now(),
	}); err != nil {
		http.Error(w, "Unable to send CreateOrder command", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusAccepted)

	rb, _ := json.Marshal(map[string]interface{}{
		"OrderID": orderID,
	})
	_, _ = w.Write(rb)
}

func (h *httpHandler) assignOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "/assign only accepts POST", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	cors(&w)

	bb, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read request body: %s", err), http.StatusInternalServerError)
		return
	}

	// NB: This could technically be an instance of the command type
	type assignRequest struct {
		OrderID  string
		AssignTo string
	}
	var ar assignRequest

	if err := json.Unmarshal(bb, &ar); err != nil {
		http.Error(w, fmt.Sprintf("Unable to unmarshal request body: %s", err), http.StatusInternalServerError)
		return
	}

	if err := h.cb.Send(r.Context(), &commands.AssignOrder{
		OrderID:  ar.OrderID,
		AssignTo: ar.AssignTo,
	}); err != nil {
		http.Error(w, fmt.Sprintf("Unable to publish AssignOrder command: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusAccepted)
}

func (h *httpHandler) listOrders(w http.ResponseWriter, r *http.Request) {
	cors(&w)

	var orders []*ges.Order

	for _, order := range h.orderStore.orders {
		orders = append(orders, order)
	}

	bytes, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to marshal orders: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusAccepted)

	_, _ = w.Write(bytes)
}

func (h *httpHandler) emit(w http.ResponseWriter, r *http.Request) {
	cors(&w)
	defer r.Body.Close()

	bb, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read emit request body: %s", err), http.StatusInternalServerError)
		return
	}

	type emitReq struct {
		OrderID   string
		EventName string
	}
	var er emitReq
	if err := json.Unmarshal(bb, &er); err != nil {
		http.Error(w, fmt.Sprintf("Unable to unmarshal emit request body: %s", err), http.StatusInternalServerError)
		return
	}

	var e interface{}

	switch er.EventName {
	case "OrderLeft":
		e = &events.OrderLeft{OrderID: er.OrderID}
	case "OrderArrived":
		e = &events.OrderArrived{OrderID: er.OrderID}
	case "OrderDelivered":
		e = &events.OrderDelivered{OrderID: er.OrderID}
	default:
		http.Error(w, fmt.Sprintf("Unknown event type: %s", er.EventName), http.StatusBadRequest)
		return
	}

	if err := h.eb.Publish(r.Context(), e); err != nil {
		http.Error(w, fmt.Sprintf("Unable to emit event: %s", err), http.StatusInternalServerError)
		return
	}
}

type event struct {
	Event string     `json:"event"`
	Order *ges.Order `json:"order"`
}

func (h *httpHandler) eventHandler(eventName string, order *ges.Order) {
	for _, l := range h.listeners {
		l <- &event{
			Event: eventName,
			Order: order,
		}
	}
}

func (h *httpHandler) _eventHandler(eventName string, order *ges.Order) {
	var toRemove []int

	for i, conn := range h.conns {
		writer, err := conn.NextWriter(websocket.TextMessage)
		if err != nil {
			fmt.Printf("Unable to obtain writer: %s\n", err)
			_ = conn.Close()
			toRemove = append(toRemove, i)
			continue
		}

		e := &event{
			Event: eventName,
			Order: order,
		}

		ob, err := json.Marshal(e)
		if err != nil {
			fmt.Printf("Unable to marshal event: %s\n", err)
			_ = conn.Close()
			toRemove = append(toRemove, i)
			continue
		}

		if _, err := writer.Write(ob); err != nil {
			fmt.Printf("Unable to write event: %s\n", err)
			_ = conn.Close()
			toRemove = append(toRemove, i)
			continue
		}
	}

	for _, i := range toRemove {
		h.conns = append(h.conns[:i], h.conns[i+1:]...)
	}
}

func newHTTPHandler(
	orderStore *OrderStore,
	cb *cqrs.CommandBus,
	eb *cqrs.EventBus,
) *httpHandler {
	handler := &httpHandler{
		mux:        http.NewServeMux(),
		orderStore: orderStore,
		cb:         cb,
		eb:         eb,
	}

	handler.mux.HandleFunc("/create", handler.createOrder)
	handler.mux.HandleFunc("/assign", handler.assignOrder)
	handler.mux.HandleFunc("/list", handler.listOrders)
	handler.mux.HandleFunc("/emit", handler.emit)

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	handler.mux.HandleFunc("/listen", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Got a request we're going to try to upgrade to a websocket.\n")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("Unable to upgrade WebSocket connection: %s", err)
			return
		}

		//handler.conns = append(handler.conns, conn)
		ch := make(chan *event)

		go func() {
			for e := range ch {
				if err := conn.WriteJSON(e); err != nil {
					fmt.Printf("Unable to send event as JSON over websocket: %s\n", err)
					continue
				}
			}
		}()

		handler.listeners = append(handler.listeners, ch)

		// The reader
		defer conn.Close()
		conn.SetReadLimit(1024)
		conn.SetReadDeadline(time.Now().Add(1 * time.Minute))
		conn.SetPongHandler(func(appData string) error {
			conn.SetReadDeadline(time.Now().Add(1 * time.Minute))
			return nil
		})

		for {
			messageType, _, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("Error reading from websocket: %s\n", err)
				break
			}
			fmt.Printf("Received message of type: %#v\n", messageType)
		}
	})

	return handler
}

func cors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
