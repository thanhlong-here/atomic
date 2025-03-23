package main

import (
	"atomic/internal/cache"
	"atomic/internal/handler"
	"atomic/internal/mongo"
	"atomic/internal/queue"
	"atomic/internal/view"
	"encoding/json"

	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handleView(w http.ResponseWriter, r *http.Request) {
	viewName := mux.Vars(r)["view"]
	cacheKey := "view:" + viewName + r.URL.RawQuery

	// Check cache trước
	if doc, found := view.GetCache(cacheKey); found {
		json.NewEncoder(w).Encode(doc)
		return
	}

	// Nếu không có cache → build lại
	result, shouldCache, err := view.RebuildView(viewName, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Trả kết quả
	if shouldCache {
		view.SetCache(cacheKey, result)
	}
	json.NewEncoder(w).Encode(result)
}

func main() {
	cache.Init()
	mongo.InitFromEnv()
	queue.Init() // ✅ Khởi động queue xử lý ghi DB

	r := mux.NewRouter()
	//Router DB
	r.HandleFunc("/db/{model}/create", handler.CreateDynamic).Methods("POST")
	r.HandleFunc("/db/{model}/{id}", handler.GetDynamic).Methods("GET")
	r.HandleFunc("/db/{model}/{id}", handler.UpdateDynamic).Methods("PUT")
	r.HandleFunc("/db/{model}/{id}", handler.DeleteDynamic).Methods("DELETE")
	//Router View
	r.HandleFunc("/view/{view}", handleView).Methods("GET")

	fmt.Println("🚀 Atomic Service is running at :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
