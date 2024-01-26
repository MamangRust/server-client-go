package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

type Item struct {
	ItemID      int    `json:"item_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

func main() {
	var rootCmd = &cobra.Command{Use: "client"}

	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get items",
		Run: func(cmd *cobra.Command, args []string) {
			getItems()
		},
	}

	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add item",
		Run: func(cmd *cobra.Command, args []string) {
			addItem()
		},
	}

	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update item",
		Run: func(cmd *cobra.Command, args []string) {
			updateItem()
		},
	}

	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete item",
		Run: func(cmd *cobra.Command, args []string) {
			deleteItem()
		},
	}

	var helpCmd = &cobra.Command{
		Use:   "help",
		Short: "Show available commands",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Available commands:")
			fmt.Println("  get    : Get items")
			fmt.Println("  add    : Add item")
			fmt.Println("  update : Update item")
			fmt.Println("  delete : Delete item")
		},
	}

	rootCmd.AddCommand(getCmd, addCmd, updateCmd, deleteCmd, helpCmd)
	rootCmd.Execute()
}

func getItems() {
	resp, err := http.Get("http://localhost:8080/items")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	var items []Item
	err = json.NewDecoder(resp.Body).Decode(&items)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	fmt.Println("Items:")
	for _, item := range items {
		fmt.Printf("ID: %d, Name: %s, Price: %d\n", item.ItemID, item.Name, item.Price)
	}
}

func addItem() {
	var newItem Item

	fmt.Print("Enter item name: ")
	fmt.Scan(&newItem.Name)
	fmt.Print("Enter new item description: ")
	fmt.Scan(&newItem.Description)
	fmt.Print("Enter item price: ")
	fmt.Scan(&newItem.Price)

	// Membuat payload JSON dari newItem
	jsonData, err := json.Marshal(newItem)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Membungkus payload JSON ke dalam bytes.Buffer
	buffer := bytes.NewBuffer(jsonData)

	resp, err := http.Post("http://localhost:8080/items", "application/json", buffer)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// Memeriksa status code dari response
	if resp.StatusCode == http.StatusOK {
		var addedItem Item
		err = json.NewDecoder(resp.Body).Decode(&addedItem)
		if err != nil {
			fmt.Println("Error decoding response:", err)
			return
		}

		fmt.Printf("Item added: ID: %d, Name: %s, Price: %d\n", addedItem.ItemID, addedItem.Name, addedItem.Price)
	} else {
		fmt.Println("Failed to add item. Status Code:", resp.StatusCode)
	}
}

func updateItem() {
	var updatedItem Item

	fmt.Print("Enter item ID to update: ")
	fmt.Scan(&updatedItem.ItemID)
	fmt.Print("Enter new item name: ")
	fmt.Scan(&updatedItem.Name)
	fmt.Print("Enter new item description: ")
	fmt.Scan(&updatedItem.Description)
	fmt.Print("Enter new item price: ")
	fmt.Scan(&updatedItem.Price)

	// Create JSON payload from updatedItem
	jsonData, err := json.Marshal(updatedItem)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Wrap JSON payload into bytes.Buffer
	buffer := bytes.NewBuffer(jsonData)

	client := &http.Client{}
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://localhost:8080/items/%d", updatedItem.ItemID), buffer)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// Check the status code of the response
	if resp.StatusCode == http.StatusOK {
		var responseItem Item
		err = json.NewDecoder(resp.Body).Decode(&responseItem)
		if err != nil {
			fmt.Println("Error decoding response:", err)
			return
		}

		fmt.Printf("Item updated: ID: %d, Name: %s, Price: %d\n", responseItem.ItemID, responseItem.Name, responseItem.Price)
	} else if resp.StatusCode == http.StatusNotFound {
		fmt.Println("Item not found")
	} else {
		fmt.Println("Failed to update item. Status Code:", resp.StatusCode)
	}
}

func deleteItem() {
	var itemID int

	fmt.Print("Enter item ID to delete: ")
	fmt.Scan(&itemID)

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:8080/items/%d", itemID), nil)
	if err != nil {
		log.Fatal("Error creating DELETE request:", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending DELETE request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Item deleted successfully")
	} else if resp.StatusCode == http.StatusNotFound {
		fmt.Println("Item not found")
	} else {
		fmt.Println("Failed to delete item. Status Code:", resp.StatusCode)
	}

}
