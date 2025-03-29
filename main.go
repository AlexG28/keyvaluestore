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
	http.HandleFunc("/keys/", func(w http.ResponseWriter, req *http.Request) {
		key := strings.TrimPrefix(req.URL.Path, "/keys/")
		if key == "" {
			http.Error(w, "can't operate on empty key", http.StatusBadRequest)
			return
		}

		switch req.Method {
		case http.MethodPut:
			handlePut(w, req, key)
		case http.MethodGet:
			handleGet(w, req, key)
		case http.MethodDelete:
			handleDelete(w, req, key)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlePut(w http.ResponseWriter, req *http.Request, key string) {

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
func handleGet(w http.ResponseWriter, _ *http.Request, key string) {
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
func handleDelete(w http.ResponseWriter, _ *http.Request, key string) {
	storeMutex.RLock()
	defer storeMutex.RUnlock()

	_, found := store[key]

	if !found {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}

	delete(store, key)

	w.WriteHeader(http.StatusOK)
}
