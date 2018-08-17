package main

import (
	"theAmazingCodeExample/app/common"
	"theAmazingCodeExample/app/router"
)

func main() {
	common.ConnectToDatabase()
	//common.CreateAWSSession()
	common.ConnectToRabbitMQ()
	router.CreateRouter()
	router.RunRouter()
}
