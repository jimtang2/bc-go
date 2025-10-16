// webhook$$bc-go;dashboard/server/http.go;grok$$
package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type SocketHandler struct {
	in       chan interface{}
	subs     map[*http.Request]chan interface{}
	mu       sync.Mutex
	streamer *Streamer
}

func NewSocketHandler() *SocketHandler {
	h := &SocketHandler{
		in:       make(chan interface{}),
		subs:     map[*http.Request]chan interface{}{},
		mu:       sync.Mutex{},
		streamer: &Streamer{},
	}
	go h.proxy()
	go h.streamer.stream(h.in)
	return h
}

func (h *SocketHandler) proxy() {
	for {
		m := <-h.in
		h.mu.Lock()
		for _, out := range h.subs {
			out <- m
		}
		h.mu.Unlock()
	}
}

func (h *SocketHandler) sub(r *http.Request) chan interface{} {
	c := make(chan interface{})
	h.mu.Lock()
	h.subs[r] = c
	h.mu.Unlock()
	go func(out chan interface{}) {
		ctx := r.Context()
		select {
		case <-ctx.Done():
			h.unsub(r)
			close(out)
		}
	}(c)
	return c
}

func (h *SocketHandler) unsub(r *http.Request) {
	h.mu.Lock()
	delete(h.subs, r)
	h.mu.Unlock()
}

func (h *SocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	defer conn.Close()
	_, cancel := context.WithCancel(context.Background())
	ch := h.sub(r)
	defer h.unsub(r)
	go func() {
		for u := range ch {
			if err := conn.WriteJSON(u); err != nil {
				log.Printf("Error sending to WebSocket: %v", err)
				return
			}
		}
	}()
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			cancel()
			break
		}
	}
}

var db *sql.DB

func serveListsFunc(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		var err error
		db, err = sql.Open("postgres", viper.GetString("db.url"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	var b []byte
	if err := db.QueryRow(`select jsonb_build_object(
  'exchanges', (select array_to_json(array_agg(id)) from exchanges),
  'pairs', (select array_to_json(array_agg(name)) from pairs))`).Scan(&b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
