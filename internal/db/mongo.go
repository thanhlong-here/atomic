package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"atomic/internal/cache"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var database *mongo.Database

func Connect() {

	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	database = client.Database(dbName)
	fmt.Println("✅ Connected to MongoDB")
}

// ListCollections: Trả danh sách tất cả collections trong database (bao gồm view)
func ListCollections() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.ListCollections(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("ListCollections error: %v", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	for cursor.Next(ctx) {
		var col bson.M
		if err := cursor.Decode(&col); err == nil {
			results = append(results, col)
		}
	}
	return results, nil
}

func GetCollection(name string) *mongo.Collection {
	if database == nil {
		panic("MongoDB is not initialized. Call db.Connect(uri) first.")
	}
	return database.Collection(name)
}

// CollectionInfo: thông tin metadata của 1 collection hoặc view
type CollectionInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Options  bson.M `json:"options,omitempty"`
	ViewOn   string `json:"viewOn,omitempty"`
	Pipeline bson.A `json:"pipeline,omitempty"`
}

// GetCollectionInfo: Trả về thông tin metadata của 1 collection bất kỳ
func GetCollectionInfo(name string) (*CollectionInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"name": name}
	cursor, err := database.ListCollections(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("GetCollectionInfo error: %v", err)
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var col bson.M
		if err := cursor.Decode(&col); err != nil {
			return nil, err
		}

		info := &CollectionInfo{
			Name:    col["name"].(string),
			Type:    col["type"].(string),
			Options: bson.M{},
		}

		// Nếu là view thì bổ sung thêm thông tin
		if info.Type == "view" {
			optionsMap, _ := col["options"].(bson.M)
			info.Options = optionsMap
			if optionsMap != nil {
				if v, ok := optionsMap["viewOn"]; ok {
					info.ViewOn = fmt.Sprint(v)
				}
				if p, ok := optionsMap["pipeline"]; ok {
					info.Pipeline = p.(bson.A)
				}
			}
		}

		return info, nil
	}

	return nil, fmt.Errorf("collection '%s' not found", name)
}

// Count: Đếm số document trong collection theo filter
func Count(model string, filter bson.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	col := GetCollection(model)
	total, err := col.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("Count error: %v", err)
	}
	return total, nil
}

// GetPaginated: Truy vấn bất kỳ collection hoặc view có phân trang
func GetPaginated(name string, filter bson.M, skip, limit int64) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := GetCollection(name)
	opts := options.Find().SetLimit(limit).SetSkip(skip)

	cursor, err := col.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("GetPaginated error: %v", err)
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

func GetPaginatedCached(model string, filter bson.M, skip, limit int64) ([]map[string]interface{}, error) {
	cacheKey := fmt.Sprintf("paginated:%s:skip:%d:limit:%d:filter:%v", model, skip, limit, filter)

	var results []map[string]interface{}

	err := cache.AutoCache(cacheKey, 60, func() (interface{}, error) {
		return GetPaginated(model, filter, skip, limit)
	}, &results)

	return results, err
}
