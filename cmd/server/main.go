package main

import (
	"log"
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

	// connnect to db
	dbConn, err := db.ConnectToDB()
	if err != nil {
		log.Println("Error conneting to DB")
		return
	}

	typingDBService := db.NewUserRepository(dbConn)
	typingUseCase := typing.NewTypingService(&typingDBService)

	authDBService := db.NewAuthRepository(dbConn)
	authUseCase := auth.NewAuthService(&authDBService)
	handler := handler.NewHandler(typingUseCase, authUseCase, logChan)
	app := routes.SetUpRoutes(handler)

	err = app.Run(":8080")
	if err != nil {
		log.Println("Error in starting the server")
		return
	}
}
