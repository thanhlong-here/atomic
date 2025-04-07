package command

import (
	"atomic/internal/db"
	"atomic/internal/ws"
)

func HandleFindMany(msg ws.WSMessage) map[string]interface{} {
	model, err := GetModel(msg.Payload)
	if err != nil {
		return errorResp(err)
	}
	filter := GetFilter(msg.Payload)
	limit := GetInt64(msg.Payload, "limit")
	skip := GetInt64(msg.Payload, "skip")

	res, err := db.GetPaginated(model, filter, skip, limit)
	if err != nil {
		return errorResp(err)
	}
	return successResp("data", res)
}

func init() {
	ws.AutoRegister(HandleFindMany)
}
