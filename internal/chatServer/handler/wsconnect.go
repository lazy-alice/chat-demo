package handler

import (
	"chat-demo/internal/chatServer/model"
	"chat-demo/internal/chatServer/pkg/ws"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

func WsPage(c *gin.Context) {
	conn,err:=(&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {return true}}).Upgrade(c.Writer,c.Request,nil)
	if err!=nil {
		http.NotFound(c.Writer,c.Request)
		return
	}
	id,err:=strconv.Atoi(c.Request.URL.Query().Get("id"))
	if err!=nil {
		log.Fatal(err)
	}
	client:=&ws.Client{
		ID:id,
		Socket: conn,
		Send:make(chan []byte),
	}
	ws.Manager.Register<-client
	messageArray,err:=model.ReadMessage(id)
	fmt.Println(messageArray)
	if err!=nil {
		log.Println(err.Error())
	}

	go client.Read()
	go client.Write()

	for _,msg:=range messageArray {
		jsonMessage,_:=json.Marshal(&model.Message{
			Sender: msg.Sender,
			Content: msg.Content,
		})
		client.Send<-jsonMessage
	}
}
