package main

import (
	"log"

	"github.com/kiwamunet/image-optim/server"
)

func main() {
	log.Println("image-optim starting....")

	server.Serve()
}
