package handlers

import (
	"encoding/json"
	"github.com/SwanHtetAungPhyo/api-gateway/models"
	"github.com/SwanHtetAungPhyo/api-gateway/services"
	"net/http"
	"sync"
)

var serviceRegistry = make(map[string]models.Services)
var mu = &sync.Mutex{}

func Registry(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var newService models.Services
	if err := json.NewDecoder(request.Body).Decode(&newService); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if newService.BasePath == "" {
		http.Error(writer, "Base path cannot be empty", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	existingServices, exist := serviceRegistry[newService.BasePath]
	if exist {

		existingServices.Instances = append(existingServices.Instances, newService.Instances...)
		serviceRegistry[newService.BasePath] = existingServices
	} else {

		serviceRegistry[newService.BasePath] = newService
	}

	writer.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(writer).Encode(&newService)
	if err != nil {
		http.Error(writer, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func Forwardding(writer http.ResponseWriter, request *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	for basePath, service := range serviceRegistry {
		if len(request.URL.Path) >= len(basePath) && request.URL.Path[:len(basePath)] == basePath {
			services.ForwardRequest(writer, request, service)
			return
		}
		http.Error(writer, "Cannot find the service", http.StatusInternalServerError)
	}
}

func RoutesListing(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Method Not allowed", http.StatusMethodNotAllowed)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(serviceRegistry)
	if err != nil {
		http.Error(writer, "Internal Error", http.StatusInternalServerError)
		return
	}
}
