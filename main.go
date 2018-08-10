package main

import (
	"theAmazingCodeExample/app/common"
	"theAmazingCodeExample/app/router"
)

func main() {
	common.ConnectToDatabase()
	//common.CreateAWSSession()
	router.CreateRouter()
	router.RunRouter()
}
