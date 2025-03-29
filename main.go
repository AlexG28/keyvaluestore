package main

import (
	"fmt"
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

	// http.HandleFunc("/keys/", put)
	http.HandleFunc("/keys/", get)
	// http.HandleFunc("/keys/", delete)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func put(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed) // find exact error
	}

	fmt.Fprintf(w, "delete\n")

}
func get(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	key := strings.TrimPrefix(req.URL.Path, "/keys/")
	if key == "" {
		http.Error(w, "can't operate on empty key", http.StatusBadRequest)
		return
	}

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
func delete(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed) // find exact error
	}

	fmt.Fprintf(w, "delete\n")
}
