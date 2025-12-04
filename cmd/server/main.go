package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/user/todo_auth/internal/database"
	"github.com/user/todo_auth/internal/handlers"
)

func main() {
	db := database.InitDB()
	defer db.Close()

	router := mux.NewRouter()

	handlers.RegisterRoutes(router, db)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Task Manager API is running!")
	}).Methods("GET")

	port := ":8181"
	fmt.Printf("Starting server at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
