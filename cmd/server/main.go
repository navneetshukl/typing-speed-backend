package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	db "typing-speed/internals/adapter/persistence"
	routes "typing-speed/internals/interface/rest/api"
	"typing-speed/internals/interface/rest/api/handler"
	"typing-speed/internals/usecase/auth"
	"typing-speed/internals/usecase/typing"
	"typing-speed/pkg/logs"
)

func main() {
	// writing log to log file
	logChan := make(chan logs.LogEntry, 1000)
	go func() {
		for v := range logChan {
			v.CreateLog()
		}
	}()

	// connect to db
	dbConn, err := db.ConnectToDB()
	if err != nil {
		log.Println("Error connecting to DB:", err)
		return
	}

	typingDBService := db.NewUserRepository(dbConn)
	typingUseCase := typing.NewTypingService(&typingDBService)

	authDBService := db.NewAuthRepository(dbConn)
	authUseCase := auth.NewAuthService(&authDBService)

	handler := handler.NewHandler(typingUseCase, authUseCase, logChan)
	router := routes.SetUpRoutes(handler)

	// Create the HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Run server in a goroutine so it doesn’t block
	go func() {
		log.Println("🚀 Server is running on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for interrupt or terminate signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Block until a signal is received
	log.Println("🛑 Shutting down server gracefully...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}
