package main

import (
	"net/http"

	"github.com/hookercookerman/lampkicking"
	"github.com/hookercookerman/lampkicking/users"

	"github.com/go-zoo/bone"
)

func main() {
	mux := bone.New()
	persistUrl := lampkicking.Getenv("PERSIST_SERVICE")

	userService := users.NewUserService(persistUrl, nil)
	mux.Post("/users", http.HandlerFunc(userService.Store))
	mux.Get("/users/:id/connections", http.HandlerFunc(userService.GetConnections))
	mux.Put("/users/:id/connections/:connection_id", http.HandlerFunc(userService.AddConnection))

	http.ListenAndServe(":"+lampkicking.Getenv("USER_SERVICE_PORT"), mux)
}
