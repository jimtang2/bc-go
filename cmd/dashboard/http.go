package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jimtang2/bc-go/lib/db"
	"github.com/jimtang2/bc-go/pkg/pb/v1"
	_ "github.com/lib/pq"
	"google.golang.org/protobuf/proto"
)

type MatchesHandler struct {
	channel  chan []byte
	upgrader websocket.Upgrader
	mu       sync.Mutex
	subs     map[*http.Request]chan []byte
	lastRcv  time.Time
}

func NewMatchesHandler() *MatchesHandler {
	handler := &MatchesHandler{
		channel: make(chan []byte),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		mu:   sync.Mutex{},
		subs: map[*http.Request]chan []byte{},
	}
	go handler.proxyStream()
	return handler
}

// proxy global channel messages to all subs
func (handler *MatchesHandler) proxyStream() {
	for {
		b := <-handler.channel
		handler.lastRcv = time.Now()
		handler.mu.Lock()
		for _, subChan := range handler.subs {
			subChan <- b
		}
		handler.mu.Unlock()
	}
}

func (handler *MatchesHandler) subscribe(r *http.Request) chan []byte {
	handler.mu.Lock()
	handler.subs[r] = make(chan []byte)
	handler.mu.Unlock()
	go func(subChan chan []byte) {
		ctx := r.Context()
		select {
		case <-ctx.Done():
			handler.unsubscribe(r)
		}
	}(handler.subs[r])
	return handler.subs[r]
}

func (handler *MatchesHandler) unsubscribe(r *http.Request) {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	close(handler.subs[r])
	delete(handler.subs, r)
}

func (handler *MatchesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := db.LogEvent("connect matches", r); err != nil {
		log.Printf("Failed to log visit: %v", err)
	}
	defer db.LogEvent("disconnect matches", r)
	conn, err := handler.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	defer conn.Close()
	_, cancel := context.WithCancel(context.Background())
	// listen to sub chan
	go func() {
		subChan := handler.subscribe(r)
		for b := range subChan {
			if err := conn.WriteMessage(websocket.BinaryMessage, b); err != nil {
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

func (handler *MatchesHandler) Health(ctx context.Context) error {
	if time.Now().Sub(handler.lastRcv) > 20*time.Second {
		return fmt.Errorf("stream broken")
	}
	return nil
}

type AlphaHandler struct {
}

func (handler *AlphaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := db.LogEvent("get alpha", r); err != nil {
		log.Printf("Failed to log visit: %v", err)
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	resp, err := db.ProfitableMatches(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := proto.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/protobuf")
	w.Write(b)
}

func (handler *AlphaHandler) Health(ctx context.Context) error {
	return db.Ping()
}

type TickersHandler struct {
	channel  chan []byte
	upgrader websocket.Upgrader
	mu       sync.Mutex
	subs     map[*http.Request]chan []byte
	lastRcv  time.Time
}

func NewTickersHandler() *TickersHandler {
	handler := &TickersHandler{
		channel: make(chan []byte),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		mu:   sync.Mutex{},
		subs: map[*http.Request]chan []byte{},
	}
	go handler.proxyStream()
	return handler
}

// proxy global channel messages to all subs
func (handler *TickersHandler) proxyStream() {
	for {
		b := <-handler.channel
		handler.lastRcv = time.Now()
		handler.mu.Lock()
		for _, subChan := range handler.subs {
			subChan <- b
		}
		handler.mu.Unlock()
	}
}

func (handler *TickersHandler) subscribe(r *http.Request) chan []byte {
	handler.mu.Lock()
	handler.subs[r] = make(chan []byte)
	handler.mu.Unlock()
	go func(subChan chan []byte) {
		ctx := r.Context()
		select {
		case <-ctx.Done():
			handler.unsubscribe(r)
		}
	}(handler.subs[r])
	return handler.subs[r]
}

func (handler *TickersHandler) unsubscribe(r *http.Request) {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	close(handler.subs[r])
	delete(handler.subs, r)
}

func (handler *TickersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := handler.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	defer conn.Close()
	_, cancel := context.WithCancel(context.Background())
	// listen to sub chan
	go func() {
		subChan := handler.subscribe(r)
		for b := range subChan {
			if err := conn.WriteMessage(websocket.BinaryMessage, b); err != nil {
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

func (handler *TickersHandler) Health(ctx context.Context) error {
	if time.Now().Sub(handler.lastRcv) > 20*time.Second {
		return fmt.Errorf("stream broken")
	}
	return nil
}

type ExchangesHandler struct {
}

func (h *ExchangesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := &pb.ExchangesResponse{
		Exchanges: db.Exchanges(),
	}
	b, err := proto.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/protobuf")
	w.Write(b)
}

func (h *ExchangesHandler) Health(ctx context.Context) error {
	return db.Ping()
}

type HealthDependency interface {
	Health(context.Context) error
}

func HealthHandler(deps ...HealthDependency) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, dep := range deps {
			if err := dep.Health(r.Context()); err != nil {
				http.Error(w, "unhealthy", http.StatusServiceUnavailable)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
