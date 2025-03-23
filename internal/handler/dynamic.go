package handler

import (
	"atomic/internal/cache"
	"atomic/internal/mongo"
	"atomic/internal/queue"
	"atomic/internal/view"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// POST /db/{model}/create
func CreateDynamic(w http.ResponseWriter, r *http.Request) {
	model := mux.Vars(r)["model"]
	var data map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Gửi vào queue xử lý ghi Mongo
	insertedID, err := queue.EnqueueAndReturnID(model, data)
	if err != nil {
		http.Error(w, "Failed to insert", http.StatusInternalServerError)
		return
	}
	view.TriggerSync(model)
	// Gắn _id vào phản hồi cho client
	data["_id"] = insertedID
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

// GET /db/{model}/{id}
func GetDynamic(w http.ResponseWriter, r *http.Request) {
	model := mux.Vars(r)["model"]
	id := mux.Vars(r)["id"]

	// Kiểm tra cache trước
	if doc, found := cache.Get(model, id); found {
		json.NewEncoder(w).Encode(doc)
		return
	}

	// Nếu không có → lấy từ Mongo và cache lại
	result, err := mongo.Find(model, id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	cache.Set(model, id, result)
	json.NewEncoder(w).Encode(result)
}

// PUT /db/{model}/{id}
func UpdateDynamic(w http.ResponseWriter, r *http.Request) {
	model := mux.Vars(r)["model"]
	id := mux.Vars(r)["id"]

	var update map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := mongo.Update(model, id, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	view.TriggerSync(model)
	// Trả lại object mới đã cập nhật
	update["id"] = id
	json.NewEncoder(w).Encode(update)
}

// DELETE /db/{model}/{id}
func DeleteDynamic(w http.ResponseWriter, r *http.Request) {
	model := mux.Vars(r)["model"]
	id := mux.Vars(r)["id"]

	err := mongo.Delete(model, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	view.TriggerSync(model)
	w.WriteHeader(http.StatusNoContent)
}

func GetDynamicView(w http.ResponseWriter, r *http.Request) {
	viewName := mux.Vars(r)["view"]
	cacheKey := "view:" + viewName + r.URL.RawQuery

	if doc, found := cache.GetRaw(cacheKey); found {
		json.NewEncoder(w).Encode(doc)
		return
	}

	data, _, err := view.RebuildView(viewName, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
}
