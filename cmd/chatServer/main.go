package main

import (
	"chat-demo/internal/chatServer/handler"
	"chat-demo/internal/chatServer/pkg/db"
	"chat-demo/internal/chatServer/pkg/ws"
	"github.com/gin-gonic/gin"
	"log"
)

//var upGrader = websocket.Upgrader{
//	CheckOrigin:func (r *http.Request) bool {
//		return true
//	},
//}
//
//var MWS map[*websocket.Conn]int



//func ping(c *gin.Context) {
//	ws,err:=upGrader.Upgrade(c.Writer,c.Request,nil)
//	if err!= nil {
//		return
//	}
//	id,err:=strconv.Atoi(c.Request.URL.Query().Get("id"))
//	if err!= nil {
//		return
//	}
//	MWS[ws]=id
//
//	go Receive1(ws)
//}
//
//func Receive1(ws *websocket.Conn) {
//	for {
//		mt,message,err:=ws.ReadMessage()
//		fmt.Println(message)
//		if err!= nil {
//			fmt.Println(err.Error())
//			break
//		}
//
//		if string(message)=="ping" {
//			message=[]byte("pong")
//		}
//		err=ws.WriteMessage(mt,message)
//
//		if err!=nil {
//			break
//		}
//
//	}
//	fmt.Println("11")
//}
//
//func SendSomeone(c *gin.Context) {
//	id,err:=strconv.Atoi(c.Request.URL.Query().Get("id"))
//	if err!=nil {
//		return
//	}
//	for k,v:=range MWS {
//		if v==id {
//			err:=k.WriteMessage(websocket.TextMessage,[]byte("nice to meet "+strconv.Itoa(id)))
//			if err!= nil {
//				return
//			}
//		}
//	}
//}

func main() {
	//MWS=make(map[*websocket.Conn]int)
	//bindAddress:="localhost:8081"
	//r:=gin.Default()
	//r.GET("/ping",ping)
	//r.GET("/sendone",SendSomeone)
	//r.Run(bindAddress)


	err:=db.InitDB()
	if err!= nil {
		log.Fatalln(err.Error())
	}
	go ws.Manager.Start()
	bindAddress:="127.0.0.1:8081"
	r:=gin.Default()
	r.GET("/ws",handler.WsPage)
	r.GET("/notice",handler.ServerNotice)
	r.Run(bindAddress)

}