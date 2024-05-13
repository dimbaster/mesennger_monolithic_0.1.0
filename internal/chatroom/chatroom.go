package chatroom

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	id   int
	conn *websocket.Conn
}

type message struct {
	Senderid int    `json:"-"`
	Text     string `json:"text"`
}

type ChatRoom struct {
	id                    int
	registeredConnections map[int]*websocket.Conn
	toUnregister          chan int
	toRegister            chan client
	messages              chan message
	shouldDelete          chan struct{}
}

func New() *ChatRoom {
	cr := &ChatRoom{
		registeredConnections: make(map[int]*websocket.Conn),
		toUnregister:          make(chan int),
		toRegister:            make(chan client),
		messages:              make(chan message, 100),
		shouldDelete:          make(chan struct{}),
	}

	go cr.startPolling()
	return cr
}

func (c *ChatRoom) startPolling() {
	defer func() {
		c.shouldDelete <- struct{}{}
	}()

	for {
		select {
		case id := <-c.toUnregister:
			delete(c.registeredConnections, id)
			if len(c.registeredConnections) == 0 {
				close(c.shouldDelete)
				return
			}
		case newClient := <-c.toRegister:
			c.registeredConnections[newClient.id] = newClient.conn
			go c.connReader(newClient)
		default:
		}

		select {
		case msg := <-c.messages:
			for clid, conn := range c.registeredConnections {
				err := conn.WriteJSON(message{
					Senderid: msg.Senderid,
					Text:     msg.Text,
				})
				if err != nil {
					log.Println("failed to write json to a conn of a user: ", clid, err)
					c.toUnregister <- clid
				}
			}
		case <-time.After(time.Minute * 3):
			return
		}
	}
}

func (c *ChatRoom) AddClient(userid int, conn *websocket.Conn) {
	c.toRegister <- client{
		id:   userid,
		conn: conn,
	}
}

func (c *ChatRoom) RemoveClient(id int) {
	c.toUnregister <- id
}

func (c *ChatRoom) GetDoneChan() chan struct{} {
	return c.shouldDelete
}

func (c *ChatRoom) connReader(cl client) {
	for {
		var msg message
		err := cl.conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			c.toUnregister <- cl.id
			break
		}

		msg.Senderid = cl.id
		c.messages <- msg
	}
}
