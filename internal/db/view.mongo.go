package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Tạo View mới
func CreateView(viewName string, source string, pipeline mongo.Pipeline) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := bson.D{
		{Key: "create", Value: viewName},
		{Key: "viewOn", Value: source},
		{Key: "pipeline", Value: pipeline},
	}

	err := database.RunCommand(ctx, cmd).Err()
	if err != nil {
		return fmt.Errorf("CreateView error: %v", err)
	}
	return nil
}

// Cập nhật View: xoá và tạo lại
func UpdateView(viewName string, source string, pipeline mongo.Pipeline) error {
	if err := DeleteView(viewName); err != nil {
		return err
	}
	return CreateView(viewName, source, pipeline)
}

// Xoá View
func DeleteView(viewName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := bson.D{{Key: "drop", Value: viewName}}
	err := database.RunCommand(ctx, cmd).Err()
	if err != nil {
		return fmt.Errorf("DeleteView error: %v", err)
	}
	return nil
}

// GetViewData: Lấy dữ liệu của một view bất kỳ (hỗ trợ filter, limit, skip)
func GetViewData(viewName string, filter bson.M, limit, skip int64) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := GetCollection(viewName)
	opts := options.Find().SetLimit(limit).SetSkip(skip)

	cursor, err := col.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("GetViewData error: %v", err)
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	for cursor.Next(ctx) {
		var doc map[string]interface{}
		if err := cursor.Decode(&doc); err == nil {
			results = append(results, doc)
		}
	}
	return results, nil
}
