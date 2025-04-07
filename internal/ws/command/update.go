package command

import (
	"atomic/internal/db"
	"atomic/internal/ws"
)

func HandleUpdate(msg ws.WSMessage) map[string]interface{} {
	model := msg.Payload["model"].(string)
	filter := toBson(msg.Payload["filter"])
	update := msg.Payload["data"].(map[string]interface{})

	res, err := db.Update(model, filter, update)
	if err != nil {
		return errorResp(err)
	}
	return successResp("updated", res.ModifiedCount)
}

func init() {
	ws.AutoRegister(HandleUpdate)
}
