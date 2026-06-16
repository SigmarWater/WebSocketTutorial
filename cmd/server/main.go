package main

import (
	"time"

	"github.com/SigmarWater/WebSocketTutorial/internal/wsserver"
	"errors"
	"net/http"
	log "github.com/sirupsen/logrus"
)

const(
	addr = "0.0.0.0:9090"
)

func main(){
	wsSrv := wsserver.NewWsServer(addr)
	log.Infof("Started ws server on: %v", addr)
	go func(){
		if err := wsSrv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed){
			log.Errorf("Error with ws server: %v", err)
		}
	}()

	time.Sleep(time.Second * 30)
	if err := wsSrv.Stop(); err != nil{
		log.Errorf("Error stopping server")
	}
}