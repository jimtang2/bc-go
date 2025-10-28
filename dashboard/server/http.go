// webhook$$bc-go;dashboard/server/http.go;grok$$
package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

type SocketHandler struct {
	upgrader websocket.Upgrader
	mu       sync.Mutex
	subs     map[*http.Request]chan []byte
}

func NewSocketHandler() *SocketHandler {
	socket := &SocketHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		mu:   sync.Mutex{},
		subs: map[*http.Request]chan []byte{},
	}
	go socket.proxy()
	return socket
}

// proxy global channel messages to all subs
func (socket *SocketHandler) proxy() {
	for {
		b := <-channel
		socket.mu.Lock()
		for _, subChan := range socket.subs {
			subChan <- b
		}
		socket.mu.Unlock()
	}
}

func (socket *SocketHandler) subscribe(r *http.Request) chan []byte {
	socket.mu.Lock()
	socket.subs[r] = make(chan []byte)
	socket.mu.Unlock()
	go func(subChan chan []byte) {
		ctx := r.Context()
		select {
		case <-ctx.Done():
			socket.unsubscribe(r)
		}
	}(socket.subs[r])
	return socket.subs[r]
}

func (socket *SocketHandler) unsubscribe(r *http.Request) {
	socket.mu.Lock()
	defer socket.mu.Unlock()
	close(socket.subs[r])
	delete(socket.subs, r)
}

func (socket *SocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := socket.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	defer conn.Close()
	_, cancel := context.WithCancel(context.Background())
	// listen to sub chan
	go func() {
		subChan := socket.subscribe(r)
		for b := range subChan {
			if err := conn.WriteMessage(websocket.TextMessage, b); err != nil {
				log.Printf("Error sending to WebSocket: %v", err)
				return
			}
		}
	}()
	// cancel sub when read errors
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			cancel()
			break
		}
	}
}
