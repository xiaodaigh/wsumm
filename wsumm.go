//package wsumm is a collection of simple helper functions to make gorilla/websocket's WriteMessage not panicable.
//Gorilla websocket's WriteMessage function can not be run at the same by two functions, otherwise
//it will panic, so I've create a semaphore on the write functions so that only one can run
//and it block the execution of the other
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
	if c.writeSemaphore == nil {
		c.createWriteSemaphore()
	}
	c.writeSemaphore <- true
	defer func() {
		<-c.writeSemaphore
	}()

	res := c.Conn.WriteJSON(v)

	return res
}

func (c *Conn) WriteMessage(messageType int, data []byte) error {
	if c.writeSemaphore == nil {
		c.createWriteSemaphore()
	}
	c.writeSemaphore <- true
	defer func() {
		<-c.writeSemaphore
	}()

	res := c.Conn.WriteMessage(messageType, data)

	return res
}

func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Conn, error) {
	c, e := u.Upgrader.Upgrade(w, r, responseHeader)
	cc := &Conn{Conn: c}
	return cc, e
}
