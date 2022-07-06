package main

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const message = "Hello, World!!!"

var upgrader = websocket.Upgrader{ //an upgrader from TCP to WebSocket
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ContextDesc(context *gin.Context) { //a function that describes the 'Hello, World!!!' context with OK status-code
	context.String(http.StatusOK, message)
}

func main() {
	//printing 'Hello, World!!!' without libs
	println(message)
	//printing 'Hello, World!!!' with websocket a-la ping-pong
	println("A ping-pong 'client' is being set up...")
	u := url.URL{
		Scheme: "ws:/",
		Host:   "localhost:8080",
		Path:   "/",
	}
	connection, _, error_code := websocket.DefaultDialer.Dial(u.String(), nil)
	if error_code != nil { //if error - error message pops up
		println("[WebsocketDial]The world won't be greeted right now due to " + error_code.Error())
	}
	error_code = connection.WriteMessage(websocket.TextMessage, []byte(message)) //a 'client' has written a message, *ping*
	if error_code != nil {                                                       //if error - error message pops up
		println("[WebsocketWrite]The world won't be greeted right now due to " + error_code.Error())
	}

	_, new_message, error_code := connection.ReadMessage() //a 'server' has read a message, *pong*
	println(new_message)
	connection.Close()
	//printing 'Hello, World!!!' with Gin
	gin_router := gin.Default()      //made a default Gin router
	gin_router.GET("/", ContextDesc) //trying to GET an empty resource from localhost with the 'Hello, World!!!' status message
	error_code = gin_router.Run()    //starting the Gin server
	if error_code != nil {           //if error - error message pops up
		panic("[GinRun]The world won't be greeted right now due to " + error_code.Error())
	}
}
