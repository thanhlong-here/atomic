package queue

import (
	"atomic/internal/mongo"
	"fmt"
)

type writeTask struct {
	Model string
	Data  map[string]interface{}
}

var writeQueue chan writeTask

func Init() {
	writeQueue = make(chan writeTask, 100)
	go worker()
}

func Enqueue(model string, data map[string]interface{}) {
	writeQueue <- writeTask{Model: model, Data: data}
}

func EnqueueAndReturnID(model string, data map[string]interface{}) (interface{}, error) {
	return mongo.Insert(model, data)
}

func worker() {
	for task := range writeQueue {
		_, err := mongo.Insert(task.Model, task.Data)
		if err != nil {
			fmt.Println("❌ Mongo Insert Error:", err)
		} else {
			fmt.Println("✅ Mongo Inserted:", task.Model)
		}
	}
}

