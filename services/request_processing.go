package services

import (
	"github.com/SwanHtetAungPhyo/api-gateway/models"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
)

var instanceIndex = make(map[string]int)
var indexMutex = &sync.Mutex{}

func GetNextInstance(service models.Services) *models.Instance {
	var selectedInstance *models.Instance
	minConnections := int32(^int32(0) >> 1)
	for i := range service.Instances {
		instance := &service.Instances[i]
		connections := atomic.LoadInt32(&instance.Connections)
		if instance.Connections < minConnections {
			minConnections = connections
			selectedInstance = instance
		}
	}
	return selectedInstance
}

func ForwardRequest(writer http.ResponseWriter, request *http.Request, service models.Services) {
	instance := GetNextInstance(service)

	if instance == nil {

		http.Error(writer, "No available instances", http.StatusServiceUnavailable)
		return
	}

	instance.Mu.Lock()
	atomic.AddInt32(&instance.Connections, 1)
	defer func() {
		atomic.AddInt32(&instance.Connections, -1)
		instance.Mu.Unlock()
	}()

	targetURL := instance.Url + request.URL.Path[len(service.BasePath):]
	req, err := http.NewRequest(request.Method, targetURL, request.Body)
	if err != nil {
		http.Error(writer, "Error creating request", http.StatusInternalServerError)
		return
	}
	req.Header = request.Header
	log.Printf("Forwarding request to %s", req.URL)

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		http.Error(writer, "Error forwarding request", http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	for key, values := range response.Header {
		for _, value := range values {
			writer.Header().Add(key, value)
		}
	}

	writer.WriteHeader(response.StatusCode)
	_, err = io.Copy(writer, response.Body)
	if err != nil {
		log.Printf("Error copying response: %v", err)
	}
}
