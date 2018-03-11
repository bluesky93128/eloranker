package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"net"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool { return true },
}

// Client is the instance of a connection
type Client struct {
	ws     *websocket.Conn
	sent   chan []byte
	closed bool
	ip     string
	// Deprecated: unused now.
	// TODO: Remove if sessions won't be used anymore
	session string

	// Room special
	room   *Room
	secret string
	voting [2]string
}

// Close connection and unregister client from the room
func (c *Client) Close() {
	if c.closed {
		return
	}

	c.closed = true
	c.ws.Close()
	close(c.sent)
	if c.room != nil {
		c.room.unregister <- c
	}
}

func (c *Client) readPump() {
	defer c.Close()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		var request requestMessage
		err = json.Unmarshal(message, &request)
		if err != nil {
			// Invalid input
			c.Error(err.Error(), "")
		} else {
			c.handleMessage(&request)
		}
	}
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer c.Close()
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-c.sent:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.ws.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.sent)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.sent)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWS is the http request handler which creates Client
func ServeWS(w http.ResponseWriter, r *http.Request) {
	header := make(http.Header)
	session, err := r.Cookie("session")
	if err != nil {
		session = &http.Cookie{
			Name:   "session",
			Value:  generateSession(),
			Path:   "/",
			MaxAge: int((time.Hour * 24 * 30 * 12 * 3).Seconds()),
		}
		header.Add("Set-Cookie", session.String())
	}

	ws, err := upgrader.Upgrade(w, r, header)
	if err != nil {
		return
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	client := &Client{
		ws:      ws,
		sent:    make(chan []byte, 256),
		ip:      ip,
		session: session.Value,
	}

	go client.readPump()
	go client.writePump()
}
