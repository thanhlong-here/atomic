package command

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func toBson(v interface{}) bson.M {
	if v == nil {
		return bson.M{}
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return bson.M{}
	}
	return m
}

func int64From(v interface{}) int64 {
	switch v := v.(type) {
	case float64:
		return int64(v)
	case int:
		return int64(v)
	case int64:
		return v
	default:
		return 0
	}
}

func successResp(key string, value interface{}) map[string]interface{} {
	return map[string]interface{}{"status": "ok", key: value}
}

func errorResp(err error) map[string]interface{} {
	return map[string]interface{}{"status": "error", "error": err.Error()}
}

//help

func GetString(payload map[string]interface{}, key string) (string, error) {
	v, ok := payload[key]
	if !ok {
		return "", fmt.Errorf("missing key: %s", key)
	}
	str, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("key %s is not a string", key)
	}
	return str, nil
}

func GetMap(payload map[string]interface{}, key string) (map[string]interface{}, error) {
	v, ok := payload[key]
	if !ok {
		return map[string]interface{}{}, nil // không có thì trả rỗng
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("key %s is not a map", key)
	}
	return m, nil
}

func GetInt64(payload map[string]interface{}, key string) int64 {
	v, ok := payload[key]
	if !ok {
		return 0
	}
	switch v := v.(type) {
	case float64:
		return int64(v)
	case int:
		return int64(v)
	case int64:
		return v
	default:
		return 0
	}
}

func GetFilter(payload map[string]interface{}) bson.M {
	m, _ := GetMap(payload, "filter")
	return m
}

func GetModel(payload map[string]interface{}) (string, error) {
	return GetString(payload, "model")
}
