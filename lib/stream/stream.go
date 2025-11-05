package stream

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/jimtang2/bc-go/lib/kafka"
)

type Driver interface {
	Open(args []string) (*websocket.Conn, error)                          // open ws connection, handshake, etc
	Close() error                                                         // clean up ws connection, etc
	OnWebsocketMessage(messageType int, b []byte) (kafka.KMessage, error) // parse websocket message and return KMessage; returned error triggers DriverConfig.OnError so drivers implementation should only return non-nil error when absolutely necessary
}

type DriverConfig struct {
	Args       []string                          // passed to Driver.Open(args []string)
	OnKMessage func(kafka.KMessage)              // do stuff with stream payload here
	OnError    func(Driver, DriverConfig, error) // handle stream error here
}

var drivers map[string]Driver

// drivers must call this to register themselves
func Register(name string, driver Driver) {
	if drivers == nil {
		drivers = map[string]Driver{}
	}
	if _, ok := drivers[name]; ok {
		log.Fatalf("driver %v cannot be re-registered", name)
	}
	drivers[name] = driver
}

// starts the stream driver with specified name and config
func Open(driverName string, driverConfig DriverConfig) {
	driver, ok := drivers[driverName]
	if !ok {
		log.Fatalf("driver %v is not registered")
	}
	ws, err := driver.Open(driverConfig.Args)
	if err != nil {
		driverConfig.OnError(driver, driverConfig, err)
		return
	}
	go func() {
		for {
			messageType, b, err := ws.ReadMessage()
			if err != nil {
				log.Printf("[%v] parse error: %v", driverName, err)
				driverConfig.OnError(driver, driverConfig, err)
				return
			}
			kmessage, err := driver.OnWebsocketMessage(messageType, b)
			if err != nil {
				log.Printf("[%v] parse error: %v", driverName, err)
				driverConfig.OnError(driver, driverConfig, err)
			} else if kmessage != nil && kmessage.IsValid() {
				driverConfig.OnKMessage(kmessage)
			}
		}
	}()
}
