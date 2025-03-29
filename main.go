package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

var store = make(map[string]string)
var storeMutex = sync.RWMutex{}

func main() {

	storeMutex.Lock()

	store["hello"] = "world"
	store["general"] = "kenobi"

	storeMutex.Unlock()

	http.HandleFunc("/keys/", func(w http.ResponseWriter, req *http.Request) {
		key := strings.TrimPrefix(req.URL.Path, "/keys/")
		if key == "" {
			http.Error(w, "can't operate on empty key", http.StatusBadRequest)
			return
		}

		switch req.Method {
		case http.MethodPut:
			put(w, req, key)
		case http.MethodGet:
			get(w, req, key)
		case http.MethodDelete:
			delete(w, req, key)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func put(w http.ResponseWriter, req *http.Request, key string) {

	body, err := io.ReadAll(req.Body)

	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	val := string(body)

	storeMutex.Lock()
	defer storeMutex.Unlock()

	store[key] = val

	w.WriteHeader(http.StatusOK)

}
func get(w http.ResponseWriter, req *http.Request, key string) {
	storeMutex.RLock()
	defer storeMutex.RUnlock()

	value, found := store[key]

	if !found {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(value))
	if err != nil {
		fmt.Printf("Error writing value for key: %s  err:%v\n", key, err)
	}

}
func delete(w http.ResponseWriter, req *http.Request, key string) {
	if req.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed) // find exact error
	}

	fmt.Fprintf(w, "delete\n")
}
