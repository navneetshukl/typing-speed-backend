package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"typing-speed/internals/adapter/external/sendmail"
	db "typing-speed/internals/adapter/persistence"
	routes "typing-speed/internals/interface/rest/api"
	"typing-speed/internals/interface/rest/api/handler"
	typeSvc "typing-speed/internals/usecase/typing"
	userSvc "typing-speed/internals/usecase/user"
	"typing-speed/pkg/logs"
)

func main() {
	logChan := make(chan logs.LogEntry, 1000)
	go func() {
		for v := range logChan {
			v.CreateLog()
		}
	}()

	dbConn, err := db.ConnectToDB()
	if err != nil {
		log.Println("Error connecting to DB:", err)
		return
	}
	mailSvc := sendmail.NewGoMail("localhost", 1025)
	userDBService := db.NewUserRepository(dbConn)
	userUseCase := userSvc.NewUserService(&userDBService, mailSvc)

	testDBService := db.NewTestRepository(dbConn)
	typingUseCase := typeSvc.NewTypingService(&userDBService, mailSvc, &testDBService)

	handler := handler.NewHandler(typingUseCase, userUseCase, logChan)
	router := routes.SetUpRoutes(handler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("🚀 Server is running on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}
