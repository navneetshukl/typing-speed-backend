package main

import (
	"log"
	db "typing-speed/internals/adapter/persistence"
	routes "typing-speed/internals/interface/rest/api"
	"typing-speed/internals/interface/rest/api/handler"
	"typing-speed/internals/usecase/typing"
)

func main(){
	// connnect to db

	dbConn,err:=db.ConnectToDB()
	if err!=nil{
		log.Println("Error conneting to DB")
		return
	}

	typingDBService:=db.NewUserRepository(dbConn)
	typingUseCase:=typing.NewTypingService(&typingDBService)
	typingHandler:=handler.NewTypingHandler(typingUseCase)
	app:=routes.SetUpRoutes(typingHandler)

	err=app.Run(":8080")
	if err!=nil{
		log.Println("Error in starting the server")
		return
	}
}