package wsserver

import (
	"net/http"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

const(
	templateDir = "../../web/templates/html"
	staticDir = "../../web/static/"
)

type WSServer interface {
	Start() error
}

type wsSrv struct {
	mux *http.ServeMux
	srv *http.Server
	// upgrader для http-запросов
	wsUpg *websocket.Upgrader
}

func NewWsServer(addr string) WSServer {
	mux := http.NewServeMux()

	return  &wsSrv{
		srv: &http.Server{
			Addr: addr,
			Handler: mux, 
		},
		mux: mux,
		wsUpg: &websocket.Upgrader{},
	}
}

func (ws *wsSrv) Start() error{
	// Подключаем статический файл по endpoint='/'
	ws.mux.Handle("/", http.FileServer(http.Dir(templateDir)))
	
	ws.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
	// Путь для подключения по websocket
	ws.mux.HandleFunc("/ws", ws.wsHandler)
	ws.mux.HandleFunc("/test", ws.testHandler)
	return ws.srv.ListenAndServe()
}

func (ws *wsSrv) testHandler(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Test is successful"))
}

func (ws *wsSrv) wsHandler(w http.ResponseWriter, r *http.Request){
	conn, err := ws.wsUpg.Upgrade(w, r, nil)
	if err != nil{
		log.Errorf("Error with websocket connection: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return 
	}

	// Вывод IP адреса клиента, который подключается к серверу
	log.Infof(conn.RemoteAddr().String())
}	