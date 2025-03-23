package mongo

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var databaseName string
var writeQueue chan writeTask

type writeTask struct {
	model string
	data  map[string]interface{}
}

func InitFromEnv() {
	_ = godotenv.Load(".env")
	uri := os.Getenv("MONGO_URI")
	databaseName = os.Getenv("MONGO_DATABASE")

	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	writeQueue = make(chan writeTask, 100)
	go writeWorker()
}

func Enqueue(model string, data map[string]interface{}) {
	writeQueue <- writeTask{model: model, data: data}
}

func writeWorker() {
	for task := range writeQueue {
		collection := client.Database(databaseName).Collection(task.model)
		_, err := collection.InsertOne(context.TODO(), task.data)
		if err != nil {
			fmt.Println("⚠️ Ghi Mongo lỗi:", err)
		} else {
			fmt.Println("✅ Đã ghi Mongo:", task.model)
		}
	}
}

func Insert(model string, data map[string]interface{}) (interface{}, error) {
	coll := client.Database(databaseName).Collection(model)
	result, err := coll.InsertOne(context.TODO(), data)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func Find(model, id string) (map[string]interface{}, error) {
	coll := client.Database(databaseName).Collection(model)
	var result map[string]interface{}
	err := coll.FindOne(context.TODO(), bson.M{"id": id}).Decode(&result)
	return result, err
}
func FindAll(model string) ([]map[string]interface{}, error) {
	coll := client.Database(databaseName).Collection(model)

	cursor, err := coll.Find(context.TODO(), map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []map[string]interface{}
	for cursor.Next(context.TODO()) {
		var doc map[string]interface{}
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		results = append(results, doc)
	}

	return results, nil
}
func Collection(model string) *mongo.Collection {
	return client.Database(databaseName).Collection(model)
}

func Update(model, id string, update map[string]interface{}) error {
	coll := client.Database(databaseName).Collection(model)
	_, err := coll.UpdateOne(context.TODO(), bson.M{"id": id}, bson.M{"$set": update})
	return err
}

func Delete(model, id string) error {
	coll := client.Database(databaseName).Collection(model)
	_, err := coll.DeleteOne(context.TODO(), bson.M{"id": id})
	return err
}
