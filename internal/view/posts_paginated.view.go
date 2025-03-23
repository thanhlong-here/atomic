package view

import (
	"atomic/internal/mongo"
	"context"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type postsPaginatedView struct{}

func (v postsPaginatedView) Name() string {
	return "posts_paginated"
}

func (v postsPaginatedView) SyncModels() []string {
	return []string{"posts"}
}

func (v postsPaginatedView) Rebuild(r *http.Request) (map[string]interface{}, bool, error) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	skip := (page - 1) * limit

	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.M{"created_at": -1})

	cursor, err := mongo.Collection("posts").Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return nil, false, err
	}
	defer cursor.Close(context.TODO())

	var results []map[string]interface{}
	for cursor.Next(context.TODO()) {
		var doc map[string]interface{}
		if err := cursor.Decode(&doc); err == nil {
			results = append(results, doc)
		}
	}

	// Chỉ cache trang 1–3
	shouldCache := page <= 3

	return map[string]interface{}{"data": results}, shouldCache, nil
}

var dview DynamicView = postsPaginatedView{}

func init() {
	Register(dview)
}
