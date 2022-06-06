package main

import "github.com/sajjad3k/contactsmanager/routes"

func main() {

	server := routes.ServerRoutes()
	server.Run()
}
