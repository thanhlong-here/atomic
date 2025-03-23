package ws

type Message struct {
	Type    string                 `json:"type"`    // subscribe, unsubscribe, emit, direct, ping, disconnect
	From    string                 `json:"from"`    // service-id hoặc user-id
	To      string                 `json:"to"`      // service-id, "broadcast", hoặc topic:xxx
	Payload map[string]interface{} `json:"payload"` // nội dung
}
