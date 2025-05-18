package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Broadcast struct {
	room string
	from int
	msg  string
}

type WsConn struct {
	id   int
	room string
	conn *websocket.Conn
}

var (
	nextID   uint64
	rooms    = make(map[string]map[int]*WsConn)
	roomsMu  sync.Mutex
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func handleWs(c *gin.Context) {
	room := c.Param("room")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Websocket upgrade error: ", err)
		return
	}

	id := int(atomic.AddUint64(&nextID, 1))
	ws := WsConn{id, room, conn}

	// register client
	roomsMu.Lock()
	if _, ok := rooms[room]; !ok {
		rooms[room] = make(map[int]*WsConn)
	}
	rooms[room][id] = &ws
	roomsMu.Unlock()

	log.Printf("client %d joined room %s\n", id, room)

	go handleWsConnection(&ws)
}

func handleWsConnection(ws *WsConn) {
	for {
		_, msg, err := ws.conn.ReadMessage()
		if err != nil {
			log.Printf("error reading message from client %d: %v", ws.id, err)
			break
		}

		broadcast(ws.room, ws.id, string(msg))
	}

	roomsMu.Lock()
	if roomMap, ok := rooms[ws.room]; ok {
		delete(roomMap, ws.id)
		if len(roomMap) == 0 {
			delete(rooms, ws.room)
		}
	}
	roomsMu.Unlock()

	ws.conn.Close()
	log.Printf("Client %d left room %s\n", ws.id, ws.room)
}

func broadcast(room string, fromID int, msg string) {
	roomsMu.Lock()
	defer roomsMu.Unlock()

	conns, ok := rooms[room]

	if !ok {
		return
	}

	for id, conn := range conns {
		if id == fromID {
			continue
		}
		err := conn.conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Printf("write error to %d: %v", id, err)
		}
	}
}

func loadTlsConfig() tls.Config {
	certPath := filepath.Join("certs", "localhost+2.pem")
	keyPath := filepath.Join("certs", "localhost+2-key.pem")

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatalf("failed to load TLS certs: %v", err)
	}

	return tls.Config{Certificates: []tls.Certificate{cert}}
}

func main() {
	tlsCfg := loadTlsConfig()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.File("index.html")
	})

	r.GET("/ws/:room", handleWs)

	s := &http.Server{
		Addr:      ":8443",
		Handler:   r,
		TLSConfig: &tlsCfg,
	}

	fmt.Println("starting server at https://localhost:8443")
	fmt.Println("please note:")
	fmt.Println("  * Note the HTTPS in the URL; there is no HTTP -> HTTPS redirect.")
	fmt.Println("  * You'll need to accept the invalid TLS certificate as it is self-signed.")
	fmt.Println("  * Some browsers or OSs may not allow the webcam to be used by multiple pages at once.")

	log.Fatal(s.ListenAndServeTLS("", ""))
}
