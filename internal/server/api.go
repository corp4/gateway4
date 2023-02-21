package server

import (
	"net"
	"sync"

	"github.com/corp4/gateway4/internal/protocol"
	log "github.com/sirupsen/logrus"
)

// API is the server API.
type API struct {
	sync.Locker
}

// NewAPI creates a new API.
func NewAPI() *API {
	return &API{}
}

// ListenAndServe starts the server.
func (a *API) ListenAndServe(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		// Accept a connection.
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		// Handle the connection.
		log.Debugf("New connection from %s", conn.RemoteAddr())
		go func() {
			defer conn.Close()

			connection := NewConnection(conn)
			for {
				if err := protocol.ReadClient(connection, conn); err != nil {
					log.Errorf("%v", err)
					return
				}
			}
		}()
	}
}
