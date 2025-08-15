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
	config.RunMigrations()

	r := router.HandleRoute()

	// CORS middleware
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	handler := corsMiddleware(r)

	if err := config.LoadEnv(); err != nil {
		log.Println("error reading port")
	}

	port := config.GetEnv("PORT", "8080")      // Get port from environment variable or default to 8080
	host := config.GetEnv("HOST", "localhost") // Get host from environment variable or default to localhost
	fmt.Printf("Starting server at port http://%s:%s\n", host, port)
	log.Fatal(http.ListenAndServe(host+":"+port, handler)) //starting server on PORT 8080
}
