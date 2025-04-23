package command

import (
	"atomic/internal/db"
	"atomic/internal/ws"
	"strings"
)

func init() {
	ws.AutoRegister(ListCollections)
}
func ListCollections(msg ws.WSMessage) map[string]interface{} {
	filterType := strings.ToLower(msg.Payload["type"].(string)) // "view", "collection", "all"

	colls, err := db.ListCollections()
	if err != nil {
		return errorResp(err)
	}

	result := []map[string]string{}
	for _, coll := range colls {
		collType := coll["type"].(string)
		collName := coll["name"].(string)

		// Lọc theo type nếu cần
		if filterType != "all" && filterType != collType {
			continue
		}

		result = append(result, map[string]string{
			"name": collName,
			"type": collType,
		})
	}

	return successResp("collections", result)
}
