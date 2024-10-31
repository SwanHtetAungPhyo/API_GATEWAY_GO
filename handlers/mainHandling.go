package handlers

import (
	"encoding/json"
	"github.com/SwanHtetAungPhyo/api-gateway/models"
	"github.com/SwanHtetAungPhyo/api-gateway/services"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"gopkg.in/yaml.v2"
)

var serviceRegistry = make(map[string]models.Services)
var mu = &sync.Mutex{}

func LoadServicesFromYAML(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var services struct {
		Services []models.Services `yaml:"services"`
	}

	err = yaml.Unmarshal(data, &services)
	if err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()

	for _, service := range services.Services {
		for i := range service.Instances {
			if service.Instances[i].Connections < 0 {
				service.Instances[i].Connections = 0
			}
		}

		serviceRegistry[service.BasePath] = service
		log.Printf("Registered service: %s with base path: %s", service.Name, service.BasePath)
	}

	return nil
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
