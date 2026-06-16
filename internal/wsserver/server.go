package wsserver

import (
	"context"
	"net"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"

	"time"

	"github.com/gorilla/websocket"
)

const(
	templateDir = "./web/templates/html"
	staticDir = "./web/static/"
)

type WSServer interface {
	Start() error
	Stop() error 
}

type wsSrv struct {
	mux *http.ServeMux
	srv *http.Server
	// upgrader для http-запросов
	wsUpg *websocket.Upgrader
	// Сессии - подключенные клиенты
	wsClients map[*websocket.Conn]struct{}
	// Мьютекс для защиты мапы
	mutex *sync.RWMutex
	// Канал для рассылки сообщений
	broadcast chan *wsMessage
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
		wsClients: make(map[*websocket.Conn]struct{}),
		mutex: &sync.RWMutex{},
		// Используем указатель для того, чтобы экономить память, используя указатель, а не копировать структуру
		broadcast: make(chan *wsMessage),
	}
}

func (ws *wsSrv) Start() error{
	// Подключаем статический файл по endpoint='/'
	ws.mux.Handle("/", http.FileServer(http.Dir(templateDir)))
	
	ws.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
	// Путь для подключения по websocket
	ws.mux.HandleFunc("/ws", ws.wsHandler)
	ws.mux.HandleFunc("/test", ws.testHandler)
	go ws.writeToClientsBroadCast()
	return ws.srv.ListenAndServe()
}

func (ws *wsSrv) Stop() error{
	log.Info("Before", ws.wsClients)
	
	close(ws.broadcast)

	ws.mutex.RLock()
	for conn := range ws.wsClients{
		// Закрываем соединение
		if err := conn.Close(); err != nil{
			log.Errorf("Error with closing: %v", err)
		}
		delete(ws.wsClients, conn)
	}
	ws.mutex.RUnlock()

	log.Info("After close", ws.wsClients)

	return ws.srv.Shutdown(context.Background())
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

	log.Infof("Client with address %s connected", conn.RemoteAddr().String())

	ws.mutex.Lock()
	ws.wsClients[conn] = struct{}{}
	ws.mutex.Unlock()

	go ws.readFromClient(conn)
}	

// Функция считывающая сообщения из каждого коннекта
func (ws *wsSrv) readFromClient(conn *websocket.Conn){
	for{
		// Выделяем память под структуру wsMessage и возвращает указатель на эту область памяти
		msg := new(wsMessage)

		// Читаем JSON в структуру msg
		if err := conn.ReadJSON(msg); err != nil{
			// Если мы не читаем, только из-за завершения работы сервера, то просто break
				if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseNormalClosure,
				) {
				log.Errorf("Unexpected WebSocket error: %v", err)
				}
				break
		}

		// Получаем только порт
		host, _, err :=  net.SplitHostPort(conn.LocalAddr().String())
		if err != nil{
			log.Errorf("Error with address split: %v", err)
		}

		// Обогощаем сообщение
		msg.IPAdress = host 
		msg.Time = time.Now().Format(time.RFC1123)

		// Отправляем в канал
		ws.broadcast <- msg 
	}

	// Если случилась ошибка при чтении, то удаляем этот connection
	ws.mutex.Lock()
	delete(ws.wsClients, conn)
	ws.mutex.Unlock()
}

func (ws *wsSrv) writeToClientsBroadCast(){
	for msg := range ws.broadcast{
		ws.mutex.RLock()
		for clients := range ws.wsClients{
			go func(){
				if err := clients.WriteJSON(msg); err != nil{
					log.Errorf("Error with writing message: %v", err)
				}
			}()
		}
		ws.mutex.RUnlock()
	}
}