package main

import (
	"log"

	"github.com/SigmarWater/WebSocketTutorial/internal/wsserver"
)

const(
	addr = "localhost:9090"
)

func main(){
	wsSrv := wsserver.NewWsServer(addr)
	log.Printf("Started ws server on: %v\n", addr)
	if err := wsSrv.Start(); err != nil{
		log.Fatalf("Error with ws server: %v", err)
	}
}