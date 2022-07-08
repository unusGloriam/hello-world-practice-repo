package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const message = "Hello, World!!!"
const port = ":80"

var upgrader = websocket.Upgrader{ //an upgrader from TCP to WebSocket
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ContextDesc(context *gin.Context) { //a function that describes the 'Hello, World!!!' context with OK status-code
	context.String(http.StatusOK, message)
}
func ServerImp(a *http.Server) { //a server behaviour implementation
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		connection, error_code := upgrader.Upgrade(w, r, nil) //upgrading a connection to a WebSocket one
		if error_code != nil {
			log.Println("[ServerImpUpgrade]The world won't be greeted right now due to " + error_code.Error())
		}
		message_type, new_message, error_code := connection.ReadMessage() //trying to read an incoming message
		if error_code != nil {
			log.Println("[ServerReadMSG]The world won't be greeted right now due to " + error_code.Error())
		}
		error_code = connection.WriteMessage(message_type, new_message) //sending back the message from a 'client'
		if error_code != nil {
			log.Println("[ServerWriteMSG]The world won't be greeted right now due to " + error_code.Error())
		}

	})

	error_code := a.ListenAndServe()
	if error_code != nil {
		log.Println("[ServerListen]Server closed")
	}
}
func ServerTwoImp(a *http.Server) { //a server behaviour implementation
	error_code := a.ListenAndServe()
	if error_code != nil {
		log.Println("[ServerListen]Server closed")
	}
}

func main() {
	//printing 'Hello, World!!!' without libs
	log.Println(message)
	//printing 'Hello, World!!!' with websocket a-la ping-pong
	srv := &http.Server{
		Addr: "localhost" + port,
	}
	go ServerImp(srv) //starting a 'server' as a separate GoRoutine
	//---------------------|A client behaviour implementation[start]|---------------------//
	u := url.URL{
		Scheme: "ws",
		Host:   "localhost" + port,
		Path:   "/",
	}
	connection, _, error_code := websocket.DefaultDialer.Dial(u.String(), nil)
	if error_code != nil { //if error - error message pops up
		log.Println("[WebsocketDial]The world won't be greeted right now due to " + error_code.Error())
	}
	error_code = connection.WriteMessage(websocket.TextMessage, []byte(message)) //a 'client' has written a message (made a *ping*)
	if error_code != nil {                                                       //if error - error message pops up
		log.Println("[WebsocketWrite]The world won't be greeted right now due to " + error_code.Error())
	}
	_, new_message, error_code := connection.ReadMessage() //recieved a response from a 'server', (got a *pong*)
	log.Printf("%s", new_message)
	connection.Close() //closing the connection
	srv.Shutdown(context.TODO())
	//---------------------|A client behaviour implementation[finish]|---------------------//
	//printing 'Hello, World!!!' with Gin
	gin_router := gin.Default()      //made a default Gin router
	gin_router.GET("/", ContextDesc) //trying to GET an empty resource from localhost with the 'Hello, World!!!' status message
	srv_two := &http.Server{
		Addr:    "localhost" + port,
		Handler: gin_router,
	}
	go ServerTwoImp(srv_two)
	duration, _ := time.ParseDuration("60s")
	time.Sleep(duration)
	srv_two.Shutdown(context.TODO())
}
