package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jimtang2/bc-go/lib/db"
	_ "github.com/lib/pq"
)

type SocketHandler struct {
	channel  chan []byte
	upgrader websocket.Upgrader
	mu       sync.Mutex
	subs     map[*http.Request]chan []byte
}

func NewSocketHandler() *SocketHandler {
	socket := &SocketHandler{
		channel: make(chan []byte),
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
		b := <-socket.channel
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

type AlphaHandler struct {
	mu        sync.Mutex
	cache     []byte
	cacheTime time.Time
}

func (h *AlphaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if time.Now().Sub(h.cacheTime) > 2*time.Second {
		h.mu.Lock()
		var err error
		h.cache, err = db.AlphaItems()
		h.mu.Unlock()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.cacheTime = time.Now()
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(h.cache)
}
