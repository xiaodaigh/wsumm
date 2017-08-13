package wsumm

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type Conn struct {
	*websocket.Conn
	readSemaphore  chan bool
	writeSemaphore chan bool
}

type Upgrader struct {
	*websocket.Upgrader
}

func (c *Conn) createWriteSemaphore() {
	tmp := make(chan bool, 1)
	c.writeSemaphore = tmp
}

func (c *Conn) createReadSemaphore() {
	tmp := make(chan bool, 1)
	c.readSemaphore = tmp
}

func (c *Conn) WriteJSON(v interface{}) error {
	if c.readSemaphore == nil {
		c.createWriteSemaphore()
	}
	c.writeSemaphore <- true
	res := c.Conn.WriteJSON(v)
	<-c.writeSemaphore
	return res
}

func (c *Conn) WriteMessage(messageType int, data []byte) error {
	if c.readSemaphore == nil {
		c.createWriteSemaphore()
	}
	c.writeSemaphore <- true
	res := c.Conn.WriteMessage(messageType, data)
	<-c.writeSemaphore
	return res
}

func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Conn, error) {
	c, e := u.Upgrader.Upgrade(w, r, responseHeader)
	cc := &Conn{Conn: c}
	return cc, e
}
