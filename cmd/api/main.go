package main

import (
	"atomic/internal/cache"
	"atomic/internal/handler"
	"atomic/internal/middleware"
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

	result, _, err := view.RebuildView(viewName, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	cache.Init()
	mongo.InitFromEnv()
	queue.Init() // ✅ Khởi động queue xử lý ghi DB

	r := mux.NewRouter()
	//Router DB Dynamic
	api := r.PathPrefix("/db").Subrouter()
	api.Use(middleware.ModelDyanmic)

	api.HandleFunc("/{model}/create", handler.CreateDynamic).Methods("POST")
	api.HandleFunc("/{model}/{id}", handler.GetDynamic).Methods("GET")
	api.HandleFunc("/{model}/{id}", handler.UpdateDynamic).Methods("PUT")
	api.HandleFunc("/{model}/{id}", handler.DeleteDynamic).Methods("DELETE")

	//Router View
	r.HandleFunc("/view/{view}", handleView).Methods("GET")

	fmt.Println("🚀 Atomic Service is running at :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
