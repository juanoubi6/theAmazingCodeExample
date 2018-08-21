package main

import (
	"theAmazingCodeExample/app/common"
	"theAmazingCodeExample/app/router"
)

func main() {
	common.ConnectToDatabase()
	//common.CreateAWSSession()
	common.ConnectToRabbitMQ()
	common.ConnectToNats()
	router.CreateRouter()
	router.RunRouter()
}
