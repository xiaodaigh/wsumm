package wsumm

import "github.com/gorilla/websocket"

type Conn struct {
	Conn           *websocket.Conn
	readSemaphore  chan bool
	writeSemaphore chan bool
}

//
func (c *Conn) WriteJSON(v interface{}) error {
	c.writeSemaphore <- true
	res := c.Conn.WriteJSON(v)
	<-c.writeSemaphore
	return res
}

func (c Conn) WriteMessage(messageType int, data []byte) error {
	c.writeSemaphore <- true
	res := c.Conn.WriteMessage(messageType, data)
	<-c.writeSemaphore
	return res
}
