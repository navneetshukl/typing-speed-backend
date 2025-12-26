package main

import (
	"log"
	"net/http"
	"os"
)

// func main() {
// 	logChan := make(chan logs.LogEntry, 1000)
// 	go func() {
// 		for v := range logChan {
// 			v.CreateLog()
// 		}
// 	}()

// 	_ = godotenv.Load()
// 	// if err != nil {
// 	// 	log.Println("Error Loading the env ", err)
// 	// 	return
// 	// }

// 	//var dbConn *sql.DB
// 	//dbConn, err := db.ConnectToDB()
// 	// if err != nil {
// 	// 	log.Println("Error connecting to DB:", err)
// 	// 	return
// 	// }
// 	// var mailSvc sendmail.MailSender
// 	// //mailSvc = sendmail.NewGoMail("localhost", 1025)
// 	// userDBService := db.NewUserRepository(dbConn)
// 	// userUseCase := userSvc.NewUserService(userDBService, mailSvc)

// 	// typingDBService := db.NewTestRepository(dbConn)
// 	// typingUseCase := typeSvc.NewTypingService(userDBService, mailSvc, typingDBService)

// 	// handler := handler.NewHandler(typingUseCase, userUseCase, logChan)
// 	// router := routes.SetUpRoutes(handler)

// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080" // Default to 8080 for local testing
// 	}
// 	srv := &http.Server{
// 		Addr:    fmt.Sprintf(":%s", port),
// 		//Handler: router,
// 	}

// 	go func() {
// 		log.Println("🚀 Server is running on port # ", port)
// 		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("listen: %s\n", err)
// 		}
// 	}()

// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit
// 	log.Println("🛑 Shutting down server gracefully...")

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	if err := srv.Shutdown(ctx); err != nil {
// 		log.Fatalf("Server forced to shutdown: %v", err)
// 	}

// 	log.Println("✅ Server exited gracefully")
// }

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Listening on port", port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
