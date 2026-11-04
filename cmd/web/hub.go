// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "time"

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
	// clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
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

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[string]map[*Client]bool),
	}
}

func (h *Hub) broadcastOnline(room string) {
	h.broadcastMsg(Message{Room: room, Type: "Online", Online: len(h.rooms[room])})
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
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
			h.broadcastMsg(message)
		}
	}
}
