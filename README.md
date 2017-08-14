 # wsumm
 
Wrapper functions to make gorilla/websocket's Write*Functions* i.e. WriteMessage and WriteJSON be blocking if another Write*function* of the same channel is being run.

# Installation
In command line

```
go get github.com/xiaodaigh/wsumm
```

# Example Usage
```go
package test_wsumm

import (
	"github.com/gorilla/websocket"
	"github.com/xiaodaigh/wsumm"
	"net/http"
)

var upgrader = wsumm.Upgrader{
	Upgrader: &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}}

func someHandler(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	go func() {
		for i := 0; i < 1000; i++ {
			conn.WriteMessage(websocket.TextMessage, []byte("testing"))
		}
	}()
	go func() {
		for i := 0; i < 1000; i++ {
			conn.WriteJSON(struct{}{})
		}
	}()
}
```
