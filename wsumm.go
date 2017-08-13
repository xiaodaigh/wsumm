//package wsumm is a collection of simple helper functions to make gorilla/websocket's WriteMessage not panicable.
//Gorilla websocket's WriteMessage function can not be run at the same by two functions, otherwise
//it will panic, so I've create a semaphore on the write functions so that only one can run
//and it block the execution of the other
package wsumm

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
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
	log.SetPrefix("wsumm writing to channel: ")
	if c.writeSemaphore == nil {
		log.Println("yes it's null so initilising")
		c.createWriteSemaphore()
		log.Println(c.writeSemaphore == nil)
	}
	then := time.Now()
	c.writeSemaphore <- true
	defer func() {
		<-c.writeSemaphore
		log.Println(time.Since(then))
	}()

	res := c.Conn.WriteJSON(v)

	return res
}

func (c *Conn) WriteMessage(messageType int, data []byte) error {
	log.SetPrefix("wsumm writing to channel: ")
	if c.writeSemaphore == nil {
		log.Println("yes it's null so initilising")
		c.createWriteSemaphore()
		log.Println(c.writeSemaphore == nil)
	}
	then := time.Now()
	c.writeSemaphore <- true
	defer func() {
		<-c.writeSemaphore
		log.Println(time.Since(then))
	}()

	res := c.Conn.WriteMessage(messageType, data)

	return res
}

func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Conn, error) {
	c, e := u.Upgrader.Upgrade(w, r, responseHeader)
	cc := &Conn{Conn: c}
	return cc, e
}
