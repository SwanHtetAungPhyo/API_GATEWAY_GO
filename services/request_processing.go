package services

import (
	"github.com/SwanHtetAungPhyo/api-gateway/models"
	"io"
	"log"
	"net/http"
	"sync"
)

var instanceIndex = make(map[string]int)
var indexMutex = &sync.Mutex{}

func GetNextInstance(service models.Services) models.Instance {
	indexMutex.Lock()
	defer indexMutex.Unlock()

	var selectedInstance *models.Instance
	minConnections := int(^uint(0) >> 1)
	for i := range service.Instances {
		instance := &service.Instances[i]
		instance.Mu.Lock()
		if instance.Connections < minConnections {
			minConnections = instance.Connections
			selectedInstance = instance
		}
		instance.Mu.Unlock()
	}
	return *selectedInstance
}

func ForwardRequest(writer http.ResponseWriter, request *http.Request, service models.Services) {
	instance := GetNextInstance(service)

	targetURL := instance.Url + request.URL.Path[len(service.BasePath):]
	req, err := http.NewRequest(request.Method, targetURL, request.Body)
	if err != nil {
		http.Error(writer, "Error creating request", http.StatusInternalServerError)
		return
	}
	req.Header = request.Header
	log.Printf("Forwarding request to the %s", req.URL)
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		http.Error(writer, "Error forwarding request", http.StatusInternalServerError)
		return
	}

	defer func() {
		_ = response.Body.Close()
	}()

	for key, values := range response.Header {
		for _, value := range values {
			writer.Header().Add(key, value)
		}
	}

	writer.WriteHeader(response.StatusCode)

	_, err = io.Copy(writer, response.Body)
	if err != nil {
		return
	}

	instance.Mu.Lock()
	instance.Connections--
	instance.Mu.Unlock()
}
