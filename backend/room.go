package main

import "github.com/garyburd/redigo/redis"

// ClientHandler Instance of client
type ClientHandler interface {
	Send(interface{})
}

// Room is the instance of running room
type Room struct {
	name    string
	clients map[ClientHandler]bool

	register   chan ClientHandler
	unregister chan ClientHandler
}

// Make initializes room instance
func (r *Room) Make(name string) {
	r.name = name
	r.clients = make(map[ClientHandler]bool)
	r.register = make(chan ClientHandler)
	r.unregister = make(chan ClientHandler)

	go r.run()
}

func (r *Room) run() {
	for {
		select {
		case c := <-r.register:
			r.clients[c] = true
			r.announceClients()
		case c := <-r.unregister:
			if _, ok := r.clients[c]; ok {
				delete(r.clients, c)
			}
			r.announceClients()
		}
	}
}

func (r *Room) announceClients() {
	count := len(r.clients)
	r.SendToEveryone(nil, map[string]interface{}{
		"event":   "room:clients",
		"clients": count,
	})
}

// GetSecret returns room secret
func (r *Room) GetSecret() (string, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.String(conn.Do("HGET", "room:"+r.name, "secret"))
}

// SendToEveryone sends data to all clients in the room, optionally excluding one client
func (r *Room) SendToEveryone(except *Client, data interface{}) {
	for c := range r.clients {
		if c != except {
			c.Send(data)
		}
	}
}
