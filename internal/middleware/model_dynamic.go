package middleware

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func ModelDyanmic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if model, ok := vars["model"]; ok {
			if !strings.HasPrefix(model, "dynamic_") {
				vars["model"] = "dynamic_" + model
			}
		}
		next.ServeHTTP(w, r)
	})
}
