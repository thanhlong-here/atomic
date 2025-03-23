package ws

type Hub struct {
	Clients    map[string]*Client
	Topics     map[string]map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Topics:     make(map[string]map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

func (h *Hub) AddToTopic(topic string, client *Client) {
	if _, ok := h.Topics[topic]; !ok {
		h.Topics[topic] = make(map[string]*Client)
	}
	h.Topics[topic][client.ID] = client
}

func (h *Hub) RemoveFromTopic(topic string, client *Client) {
	if subs, ok := h.Topics[topic]; ok {
		delete(subs, client.ID)
		if len(subs) == 0 {
			delete(h.Topics, topic)
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.ID] = client
			go client.WritePump()

		case client := <-h.Unregister:
			delete(h.Clients, client.ID)
			for topic := range client.Topics {
				h.RemoveFromTopic(topic, client)
			}
			close(client.Send)

		case msg := <-h.Broadcast:
			if len(msg.To) > 6 && msg.To[:6] == "topic:" {
				topic := msg.To[6:]
				if subs, ok := h.Topics[topic]; ok {
					for _, c := range subs {
						c.Send <- msg
					}
				}
			} else if msg.To == "broadcast" {
				for _, c := range h.Clients {
					c.Send <- msg
				}
			} else if target, ok := h.Clients[msg.To]; ok {
				target.Send <- msg
			}
		}
	}
}
