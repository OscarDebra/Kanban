package main

import (
	"log"
	"net/http"

	"kanban/backend/internal/db"
	"kanban/backend/internal/handlers"
	"kanban/backend/internal/middleware"
	"kanban/backend/internal/scheduler"
	"kanban/backend/internal/ws"
)

func main() {
	database := db.Connect()
	defer database.Close()

	hub := ws.NewHub()
	go scheduler.Start(database)

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("POST /auth/register", handlers.Register(database))
	mux.HandleFunc("POST /auth/login", handlers.Login(database))

	// Protected routes
	protected := http.NewServeMux()
	protected.HandleFunc("GET /boards", handlers.GetBoards(database))
	protected.HandleFunc("POST /boards", handlers.CreateBoard(database))
	protected.HandleFunc("GET /boards/{id}", handlers.GetBoard(database))
	protected.HandleFunc("DELETE /boards/{id}", handlers.DeleteBoard(database))

	protected.HandleFunc("POST /boards/{id}/columns", handlers.CreateColumn(database, hub))
	protected.HandleFunc("PATCH /columns/{id}", handlers.UpdateColumn(database, hub))
	protected.HandleFunc("DELETE /columns/{id}", handlers.DeleteColumn(database, hub))

	protected.HandleFunc("POST /columns/{id}/tasks", handlers.CreateTask(database, hub))
	protected.HandleFunc("GET /tasks/{id}", handlers.GetTask(database))
	protected.HandleFunc("PATCH /tasks/{id}", handlers.UpdateTask(database, hub))
	protected.HandleFunc("DELETE /tasks/{id}", handlers.DeleteTask(database, hub))

	protected.HandleFunc("POST /tasks/{id}/comments", handlers.CreateComment(database, hub))
	protected.HandleFunc("POST /tasks/{id}/attachments", handlers.UploadAttachment(database))
	protected.HandleFunc("GET /attachments/{id}", handlers.DownloadAttachment(database))

	protected.HandleFunc("GET /ws", handlers.HandleWebSocket(hub, database))

	mux.Handle("/", middleware.RequireAuth(protected))

	log.Println("backend listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}