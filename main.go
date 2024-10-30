package main

import (
	"github.com/SwanHtetAungPhyo/api-gateway/utils"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	router := utils.InitRouter()
	err := godotenv.Load()
	if err != nil {
		return
	}

	port := os.Getenv("PORT")
	if port == " " {
		port = "5000"
	}
	log.Println("ApiGateway is listening on the port :8080")
	if err = http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Cannot server at the Port: %s", err.Error())
		return
	}

}
