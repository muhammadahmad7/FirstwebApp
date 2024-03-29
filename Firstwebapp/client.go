package main

import (
	"github.com/gorilla/websocket"
	"time"
)

// client represents a single chatting user.
type client struct {

	// socket is the web socket for this client.
	socket *websocket.Conn

	// send is a channel on which messages are sent.
	send chan *message

	// room is the room this client is chatting in.
	room *room
	userData map[string]interface{}
}

func (c *client) read() {
	for {
var msg *message
err :=c.socket.ReadJSON(&msg)
if err!=nil{
	return
}
msg.When=time.Now()
msg.Name = c.userData["name"].(string)
c.room.forward<-msg
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if  err != nil {
			break
		}
	}
	c.socket.Close()
}