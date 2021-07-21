package model

import (
	"chat-demo/internal/chatServer/pkg/db"
	"log"
)

type Message struct {
	Id        int    `json:"id,omitempty"`
	Sender    int    `json:"sender,omitempty"`
	Recipient int    `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
	Status    int    `json:"status,omitempty"`
}

// ReadMessage reader上线就将发给他的消息读到前端，将未读消息状态改成已读
func ReadMessage(recipient int) ([]Message, error) {
	var msg = Message{}
	var message []Message
	//sql := "select * from message where recipient = " + strconv.Itoa(recipient)
	rows, err := db.DB.Query("select * from message where recipient = ?",recipient)
	if err != nil {
		return message, err
	}
	for rows.Next() {
		err:=rows.Scan(&msg.Id, &msg.Content, &msg.Sender, &msg.Recipient, &msg.Status)
		if err!= nil {
			log.Println(err.Error())
		}
		message = append(message, msg)
	}
	for _, v := range message {
		if v.Status == 0 {
			db.DB.Exec("update message set status = 1 where id = ?",v.Id)
		}
	}
	//fmt.Println(message)
	return message, nil
}

// WriteMessage 将各种消息写入message表
func WriteMessage(message Message) error {
	_, err := db.DB.Exec("insert into message(content,recipient,sender) values(?,?,?)", message.Content,
		message.Recipient, message.Sender)
	if err != nil {
		return err
	}
	return nil
}
