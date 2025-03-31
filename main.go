package main

import (
	"crypto/sha256"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	thisPort := *(flag.Int("myport", 8080, "Current instance port"))
	otherPort := *(flag.Int("otherport", 8081, "Other instance port"))

	flag.Parse()

	http.HandleFunc("/final/", func(w http.ResponseWriter, req *http.Request) {
		key := strings.TrimPrefix(req.URL.Path, "/final/")
		if key == "" {
			http.Error(w, "can't operate on empty key", http.StatusBadRequest)
			return
		}

		handlerMux(w, req, key)
	})

	http.HandleFunc("/keys/", func(w http.ResponseWriter, req *http.Request) {
		key := strings.TrimPrefix(req.URL.Path, "/keys/")

		if key == "" {
			http.Error(w, "can't operate on empty key", http.StatusBadRequest)
			return
		}

		hash := sha256.Sum224([]byte(key))
		hashInt := binary.BigEndian.Uint16(hash[:2])

		destination := hashInt % 2

		if destination == 0 {
			handlerMux(w, req, key)
		} else {
			fmt.Printf("shit goes to port: %v\n", otherPort)
		}

		fmt.Printf("val: %v\n", hashInt)
	})

	ports := fmt.Sprintf(":%v", thisPort)

	log.Fatal(http.ListenAndServe(ports, nil))
}

func handlerMux(w http.ResponseWriter, req *http.Request, key string) {
	switch req.Method {
	case http.MethodPut:
		handlePut(w, req, key)
	case http.MethodGet:
		handleGet(w, req, key)
	case http.MethodDelete:
		handleDelete(w, req, key)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handlePut(w http.ResponseWriter, req *http.Request, key string) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()
	val := string(body)

	Put(key, val)

	w.WriteHeader(http.StatusOK)
}

func handleGet(w http.ResponseWriter, _ *http.Request, key string) {
	value, err := Get(key)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(value))
	if err != nil {
		fmt.Printf("Error writing value for key: %s  err:%v\n", key, err)
	}

}

func handleDelete(w http.ResponseWriter, _ *http.Request, key string) {
	err := Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
