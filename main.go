package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-redis/redis/v9"
)

type setRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "",
		DB:           0,
		MaxIdleConns: 5,
	})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			res, _ := json.Marshal(response{
				Code:    http.StatusMethodNotAllowed,
				Message: "invalid method",
			})
			w.Write(res)
			return
		}

		if r.Header.Get("content-type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			res, _ := json.Marshal(response{
				Code:    http.StatusUnsupportedMediaType,
				Message: "invalid content-type",
			})
			w.Write(res)
			return
		}

		payload, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var req setRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(response{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("failed to parse payload, cause: %v", err),
			})
			w.Write(res)
			return
		}

		if req.Key == "" || req.Value == "" {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(response{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("key and value is required."),
			})
			w.Write(res)
			return
		}

		expired := time.Now().Add(time.Minute)
		if err := rdb.Set(r.Context(), req.Key, req.Value, time.Until(expired)).Err(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(response{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("failed to set key-value, cause: %v", err),
			})
			w.Write(res)
			return
		}

		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(response{
			Code:    http.StatusOK,
			Message: fmt.Sprintf("key [%s] and value [%s] have been stored and will be expired at %s", req.Key, req.Value, expired.Local().Format(time.RFC3339)),
		})
		w.Write(res)
	})

	http.HandleFunc("/get/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			res, _ := json.Marshal(response{
				Code:    http.StatusMethodNotAllowed,
				Message: "invalid method",
			})
			w.Write(res)
			return
		}

		queries := r.URL.Query()
		if _, exist := queries["key"]; !exist {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(response{
				Code:    http.StatusBadRequest,
				Message: "invalid query",
			})
			w.Write(res)
			return
		}

		key := queries["key"][0]
		value, err := rdb.Get(r.Context(), key).Result()
		if err != nil {
			if err == redis.Nil {
				w.WriteHeader(http.StatusBadRequest)
				res, _ := json.Marshal(response{
					Code:    http.StatusBadRequest,
					Message: fmt.Sprintf("there is no stored value of key [%s], cause: key is incorrect or value has expired", key),
				})
				w.Write(res)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(response{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("failed to get stored value of key [%s], cause: %v", key, err),
			})
			w.Write(res)
			return
		}

		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(response{
			Code:    http.StatusOK,
			Message: fmt.Sprintf("stored value of key [%s] is [%s]", key, value),
		})
		w.Write(res)
	})

	if err := http.ListenAndServe(":8888", nil); err != nil {
		panic(err)
	}
}
