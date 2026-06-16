package main

import (
	"github.com/SigmarWater/WebSocketTutorial/internal/wsserver"
	log "github.com/sirupsen/logrus"
)

const(
	addr = "0.0.0.0:9090"
)

func main(){
	wsSrv := wsserver.NewWsServer(addr)
	log.Infof("Started ws server on: %v", addr)
	if err := wsSrv.Start(); err != nil{
		log.Errorf("Error with ws server: %v", err)
	}
}