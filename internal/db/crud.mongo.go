package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create: Thêm 1 document vào collection động
func Create(model string, data map[string]interface{}) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	col := GetCollection(model)
	res, err := col.InsertOne(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("Create error: %v", err)
	}
	return res, nil
}

// Update: Cập nhật 1 document trong collection động
func Update(model string, filter bson.M, update map[string]interface{}) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	col := GetCollection(model)
	res, err := col.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return nil, fmt.Errorf("Update error: %v", err)
	}
	return res, nil
}

// Delete: Xoá 1 document trong collection động
func Delete(model string, filter bson.M) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	col := GetCollection(model)
	res, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("Delete error: %v", err)
	}
	return res, nil
}

// FindOne: Tìm 1 document trong collection động
func FindOne(model string, filter bson.M) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	col := GetCollection(model)
	var result map[string]interface{}
	err := col.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("FindOne error: %v", err)
	}
	return result, nil
}

// FindMany: Tìm nhiều document trong collection động
func FindMany(model string, filter bson.M, limit int64, skip int64) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := GetCollection(model)

	opts := options.Find().SetLimit(limit).SetSkip(skip)
	cursor, err := col.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("FindMany error: %v", err)
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	for cursor.Next(ctx) {
		var doc map[string]interface{}
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		results = append(results, doc)
	}
	return results, nil
}

// SoftDelete: Đánh dấu document là đã xoá (không xoá khỏi DB)
func SoftDelete(model string, filter bson.M) (*mongo.UpdateResult, error) {
	update := bson.M{"deleted": true, "deletedAt": time.Now()}
	return Update(model, filter, update)
}
