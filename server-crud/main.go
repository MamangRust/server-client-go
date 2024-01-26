package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Item struct {
	ItemID      int    `json:"item_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type DataStore struct {
	mu     sync.RWMutex
	items  map[int]Item
	nextID int
}

func NewDataStore() *DataStore {
	return &DataStore{
		items:  make(map[int]Item),
		nextID: 1,
	}
}

func (ds *DataStore) Create(item Item) int {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	item.ItemID = ds.nextID
	ds.items[item.ItemID] = item
	ds.nextID++

	return item.ItemID
}

func (ds *DataStore) Read(itemID int) (Item, bool) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	item, ok := ds.items[itemID]
	return item, ok
}

func (ds *DataStore) Update(itemID int, newItem Item) bool {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if _, ok := ds.items[itemID]; ok {
		newItem.ItemID = itemID
		ds.items[itemID] = newItem
		return true
	}

	return false
}

func (ds *DataStore) Delete(itemID int) bool {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if _, ok := ds.items[itemID]; ok {
		delete(ds.items, itemID)
		return true
	}

	return false
}

func main() {
	dataStore := NewDataStore()

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getItems(w, dataStore)
		case http.MethodPost:
			createItem(w, r, dataStore)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/items/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getItem(w, r, dataStore)
		case http.MethodPut:
			updateItem(w, r, dataStore)
		case http.MethodDelete:
			deleteItem(w, r, dataStore)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}

func getItems(w http.ResponseWriter, dataStore *DataStore) {
	dataStore.mu.RLock()
	defer dataStore.mu.RUnlock()

	items := make([]Item, 0, len(dataStore.items))
	for _, item := range dataStore.items {
		items = append(items, item)
	}

	sendJSONResponse(w, items)
}

func createItem(w http.ResponseWriter, r *http.Request, dataStore *DataStore) {
	var newItem Item
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	itemID := dataStore.Create(newItem)

	newItem.ItemID = itemID

	sendJSONResponse(w, newItem)
}

func getItem(w http.ResponseWriter, r *http.Request, dataStore *DataStore) {
	itemID := extractIDFromURL(r)
	item, ok := dataStore.Read(itemID)
	if !ok {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	sendJSONResponse(w, item)
}

func updateItem(w http.ResponseWriter, r *http.Request, dataStore *DataStore) {
	itemID := extractIDFromURL(r)

	var newItem Item
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if !dataStore.Update(itemID, newItem) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	sendJSONResponse(w, map[string]string{"message": "Item updated successfully"})
}

func deleteItem(w http.ResponseWriter, r *http.Request, dataStore *DataStore) {
	itemID := extractIDFromURL(r)

	if !dataStore.Delete(itemID) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	sendJSONResponse(w, map[string]string{"message": "Item deleted successfully"})
}

func extractIDFromURL(r *http.Request) int {
	var id int
	fmt.Sscanf(r.URL.Path, "/items/%d", &id)
	return id
}

func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
