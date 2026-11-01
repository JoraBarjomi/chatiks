// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "time"

type Message struct {
	Type string
	User string
	Text string
	Time time.Time
	Online int
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func (h *Hub) broadcastMsg(msg Message) {
	for client := range h.clients {
		select {
		case client.send <- msg:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) broadcastOnline() {
	h.broadcastMsg(Message{Type: "Online", Online: len(h.clients)})
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.broadcastMsg(Message{Type: "System", User: "Hub", Text: client.username + " connected", Time: time.Now()})
			h.broadcastOnline()
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.broadcastMsg(Message{Type: "System", User: "Hub", Text: client.username + " disconnected", Time: time.Now()})
				h.broadcastOnline()
			}
		case message := <-h.broadcast:
			h.broadcastMsg(message)
		}
	}
}
