package handler

import (
	"chat-demo/internal/chatServer/pkg/ws"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func ServerNotice(c *gin.Context) {
	id, err := strconv.Atoi(c.Request.URL.Query().Get("id"))
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = ws.Manager.WriteToSb(id)
	if err!=nil {
		log.Println(err.Error())
	}
}
