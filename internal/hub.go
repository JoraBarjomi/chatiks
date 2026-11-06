package internal

import (
	"context"
	"log"
	"time"
)

type Message struct {
	Room   string
	Type   string
	User   string
	Text   string
	Time   time.Time
	Online int
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	rooms map[string]map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	store *Store
	historyLen int
}

func NewHub(store *Store) *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[string]map[*Client]bool),
		store: store,
		historyLen: 50,
	}
}

func (h *Hub) unregClient(client *Client) {
	close(client.send)
	delete(h.rooms[client.room], client)
	h.broadcastMsg(Message{Room: client.room, Type: "System", User: "Hub", Text: client.username + " disconnected", Time: time.Now()})
	h.broadcastOnline(client.room)
}

func (h *Hub) broadcastMsg(msg Message) {
	for client := range h.rooms[msg.Room] {
		select {
		case client.send <- msg:
		default:
			h.unregClient(client)
		}
	}
}

func (h *Hub) broadcastOnline(room string) {
	h.broadcastMsg(Message{Room: room, Type: "Online", Online: len(h.rooms[room])})
}

func (h *Hub) sendHistory(c *Client) {
	if h.store == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	msgs, err := h.store.LastMessages(ctx, c.room, h.historyLen)
	if err != nil {
		log.Println("history for room %q: %v", c.room, err)
		return
	}

	for _, m := range msgs {
		select {
		case c.send <- m:
		default:
			return
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.sendHistory(client)

			room, ok := h.rooms[client.room]
			if !ok {
				room = make(map[*Client]bool)
				h.rooms[client.room] = room
			}
			room[client] = true
			h.broadcastMsg(Message{Room: client.room, Type: "System", User: "Hub", Text: client.username + " connected", Time: time.Now()})
			h.broadcastOnline(client.room)
		case client := <-h.unregister:
			if _, ok := h.rooms[client.room][client]; ok {
				h.unregClient(client)
			}
		case message := <-h.broadcast:
			if h.store != nil && message.Type == "" {
				ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
				if err := h.store.SaveMessage(ctx, message); err != nil {
					log.Println("save message: %v", err)
				}
				cancel()
			}
			h.broadcastMsg(message)
		}
	}
}
