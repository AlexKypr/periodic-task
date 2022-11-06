package utils

import (
	"encoding/json"
	"net/http"
)

func ToJSON(rw http.ResponseWriter, data interface{}, code int) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.WriteHeader(code)
	json.NewEncoder(rw).Encode(data)
}
