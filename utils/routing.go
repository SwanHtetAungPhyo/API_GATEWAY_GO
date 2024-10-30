package utils

import (
	"github.com/SwanHtetAungPhyo/api-gateway/handlers"
	"net/http"
)

func InitRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/register", handlers.Registry)
	mux.HandleFunc("/", handlers.Forwardding)
	mux.HandleFunc("/routes", handlers.RoutesListing)
	return mux
}
