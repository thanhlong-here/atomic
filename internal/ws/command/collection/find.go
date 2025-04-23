package command

import (
	"atomic/internal/db"
	"atomic/internal/ws"
)

func FindCollection(msg ws.WSMessage) map[string]interface{} {
	model := msg.Payload["model"].(string)
	filter := toBson(msg.Payload["filter"])

	res, err := db.FindOne(model, filter)
	if err != nil {
		return errorResp(err)
	}
	return successResp("data", res)
}

func init() {
	ws.AutoRegister(FindCollection)
}
