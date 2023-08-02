// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	chatspace *Chatspace

	messages []*ChatMsg

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *ChatMsg

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub(chatspace *Chatspace) *Hub {
	return &Hub{
		chatspace:  chatspace,
		messages:   make([]*ChatMsg, 0),
		broadcast:  make(chan *ChatMsg),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			client.send <- &ChatspaceInfo{
				Name:     h.chatspace.Name,
				Users:    h.chatspace.Users,
				Channels: h.chatspace.Channels,
				// TODO: filter messages to exclude DMs to others and msgs to channels to which this client is not subscribed.
				Messages: append(h.messages[:0:0], h.messages...), // TODO: use linked list instead, to avoid copy?
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			h.messages = append(h.messages, message)
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
