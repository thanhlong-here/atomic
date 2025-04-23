package command

import (
	"atomic/internal/db"
	"atomic/internal/ws"
)

func DeleteCollection(msg ws.WSMessage) map[string]interface{} {
	model := msg.Payload["model"].(string)
	filter := toBson(msg.Payload["filter"])

	res, err := db.Delete(model, filter)
	if err != nil {
		return errorResp(err)
	}
	return successResp("deleted", res.DeletedCount)
}

func init() {
	ws.AutoRegister(DeleteCollection)
}
