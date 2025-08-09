// @title Think Battleground API
// @version 1.0
// @description This is my sample API with Gorilla Mux
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"fmt"
	"log"
	"net/http"
	"thinkbattleground-apis/config"
	"thinkbattleground-apis/router"
)

func main() {
	config.DbConnection()

	r := router.HandleRoute()

	if err := config.LoadEnv(); err != nil {
		log.Println("error reading port")
	}

	port := config.GetEnv("PORT", "8080") // Get port from environment variable or default to 8080
	host := config.GetEnv("HOST", "localhost") // Get host from environment variable or default to localhost
	fmt.Printf("Starting server at port http://%s:%s\n", host, port)
	log.Fatal(http.ListenAndServe(host+":"+port, r)) //starting server on PORT 8080
}
