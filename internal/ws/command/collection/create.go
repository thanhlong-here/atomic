package command

import (
	"atomic/internal/db"
	"atomic/internal/ws"
)

func CreateCollection(msg ws.WSMessage) map[string]interface{} {
	model := msg.Payload["model"].(string)
	data := msg.Payload["data"].(map[string]interface{})

	res, err := db.Create(model, data)
	if err != nil {
		return errorResp(err)
	}
	return successResp("inserted_id", res.InsertedID)
}

func init() {
	ws.AutoRegister(CreateCollection)
}
