package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	message := &Response{Message: "Hello from services One"}
	err := json.NewEncoder(w).Encode(message)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	fmt.Println("Handled request for /hello")
}

func main() {
	mux1 := http.NewServeMux()
	mux1.HandleFunc("/hello", HelloHandler)

	mux2 := http.NewServeMux()
	mux2.HandleFunc("/hello", HelloHandler)

	go func() {
		err := http.ListenAndServe(":3001", mux1)
		if err != nil {

		}
	}()
	go func() {
		err := http.ListenAndServe(":3002", mux2)
		if err != nil {

		}
	}()

	select {}
}
