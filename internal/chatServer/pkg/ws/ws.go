package ws

import (
	"chat-demo/internal/chatServer/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
)

type ClientManager struct {
	clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

type Client struct {
	ID     int
	Socket *websocket.Conn
	Send   chan []byte
}

var Manager = ClientManager{
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func (c *ClientManager) Start() {
	for {
		select {
		case conn := <-c.Register:
			c.clients[conn] = true
			jsonMessage, _ := json.Marshal(&model.Message{Content: "socket " + strconv.Itoa(conn.ID) + " has connected."})
			c.Send(jsonMessage, conn)
		case conn := <-c.Unregister:
			if _, ok := c.clients[conn]; ok {
				close(conn.Send)
				delete(c.clients, conn)
				jsonMessage, _ := json.Marshal(&model.Message{Content: "socket " + strconv.Itoa(conn.ID) + " has disconnected."})
				c.Send(jsonMessage, conn)
			}
		case message := <-c.Broadcast:
			for conn := range c.clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(c.clients, conn)
				}
			}
		}
	}
}

func (c *ClientManager) Send(message []byte, ignore *Client) {
	for conn := range c.clients {
		if conn != ignore {
			conn.Send <- message
		}
	}
}

func (c *ClientManager) SendToOne(message []byte, conn *Client) {
	conn.Send <- message
}

//Read 读取前端发来的消息
func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		fmt.Println(string(message))
		if err != nil {
			Manager.Unregister <- c
			c.Socket.Close()
			break
		}
		var msg = model.Message{}
		err = json.Unmarshal(message, &msg)
		if err != nil {
			msg.Sender = c.ID
			msg.Recipient = 0
			msg.Content = string(message)
			err = model.WriteMessage(msg)
			if err != nil {
				log.Println(err.Error())
				return
			}
			jsonMessage, _ := json.Marshal(&model.Message{Sender: c.ID, Content: string(message)})
			Manager.Broadcast <- jsonMessage
		} else {
			msg.Sender = c.ID
			err = model.WriteMessage(msg)
			if err != nil {
				log.Println(err.Error())
				return
			}
			for client, isConnect := range Manager.clients {
				if client.ID == msg.Recipient {
					if isConnect == true {
						jsonMessage, _ := json.Marshal(&model.Message{Sender: c.ID, Content: msg.Content})
						client.Send <- jsonMessage
					}
				}
			}
		}
	}
}

//Write 写回消息给前端
func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.TextMessage, []byte{})
				return
			}

			_ = c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// WriteToSb 服务端发送消息给某个客户端
func (c *ClientManager) WriteToSb(id int) error {
	var message = model.Message{
		Sender: 0,//0会自动过滤，前端json不显示
		Recipient: id,
		Content:   "借阅即将到期",
	}
	err := model.WriteMessage(message)
	if err != nil {
		return err
	}
	for k, isConnect := range c.clients {
		if k.ID == id {
			if isConnect == true {
				jsonMessage, err := json.Marshal(&model.Message{Sender: message.Id, Content: message.Content})
				if err != nil {
					log.Println(err.Error())
				}
				k.Send <- jsonMessage
			}
		}
	}
	return nil
}
