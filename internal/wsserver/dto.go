package wsserver

type wsMessage struct{
	IPAdress string `json:"address"`
	Message string `json:"message"`
	Time string `json:"time"`
}