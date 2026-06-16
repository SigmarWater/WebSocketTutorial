package wsserver

import "net/http"

type WSServer interface {
	Start() error
}

type wsSrv struct {
	mux *http.ServeMux
	srv *http.Server
}

func NewWsServer(addr string) WSServer {
	mux := http.NewServeMux()

	return  &wsSrv{
		srv: &http.Server{
			Addr: addr,
			Handler: mux, 
		},
		mux: mux,
	}
}

func (ws *wsSrv) Start() error{
	ws.mux.HandleFunc("/test", ws.testHandler)
	return ws.srv.ListenAndServe()
}

func (ws *wsSrv) testHandler(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Test is successful"))
}