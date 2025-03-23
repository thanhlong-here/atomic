package ws

import (
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID     string
	Conn   *websocket.Conn
	Send   chan Message
	Topics map[string]bool
}

func (c *Client) IsSubscribed(topic string) bool {
	return c.Topics[topic]
}

func (c *Client) Subscribe(topic string) {
	c.Topics[topic] = true
}

func (c *Client) Unsubscribe(topic string) {
	delete(c.Topics, topic)
}

func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("❌ Read error [%s]: %v", c.ID, err)
			break
		}

		switch msg.Type {
		case "ping":
			c.Send <- Message{Type: "pong", To: c.ID}
		case "disconnect":
			return
		case "subscribe":
			topic := strings.TrimPrefix(msg.To, "topic:")
			c.Subscribe(topic)
			hub.AddToTopic(topic, c)
		case "unsubscribe":
			topic := strings.TrimPrefix(msg.To, "topic:")
			c.Unsubscribe(topic)
			hub.RemoveFromTopic(topic, c)
		default:
			hub.Broadcast <- msg
		}
	}
}

func (c *Client) WritePump() {
	for msg := range c.Send {
		_ = c.Conn.WriteJSON(msg)
	}
}
