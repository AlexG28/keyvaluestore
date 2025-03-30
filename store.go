package main

import (
	"fmt"
	"sync"
)

var store = make(map[string]string)
var storeMutex = sync.RWMutex{}

func Put(key string, value string) {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	store[key] = value
}

func Get(key string) (string, error) {
	storeMutex.RLock()
	defer storeMutex.RUnlock()

	value, found := store[key]
	if !found {
		return "", fmt.Errorf("failed to get")
	}
	return value, nil
}

func Delete(key string) error {
	storeMutex.RLock()
	defer storeMutex.RUnlock()

	_, found := store[key]

	if !found {
		return fmt.Errorf("failed to delete")
	}

	delete(store, key)
	return nil
}
