package server

import (
	"io"
	"math/rand"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/corp4/gateway4/internal/protocol"
)

// Connection is a connection coming from a client.
type Connection struct {
	sync.Locker
	*log.Entry

	conn  net.Conn
	close chan struct{}

	keepAliveId chan uint64
}

// NewConnection creates a new connection from a net.Conn.
func NewConnection(conn net.Conn) *Connection {
	connection := &Connection{
		Locker: &sync.Mutex{},
		Entry: log.WithFields(log.Fields{
			"remote_addr": conn.RemoteAddr().String(),
		}),
		conn:        conn,
		keepAliveId: make(chan uint64, 1),
		close:       make(chan struct{}),
	}

	go connection.HealthCheck()
	return connection
}

// HealthCheck checks if the connection is still alive.
func (c *Connection) HealthCheck() {
	for {
		// Wait 15 seconds before sending a keep alive packet.
		select {
		case <-time.NewTimer(15 * time.Second).C:
		case <-c.close:
			return
		}

		// Send a keep alive packet.
		randId := rand.Uint64()
		keepAlive := protocol.KeepAlive{ID: randId}
		now := time.Now()
		if err := protocol.WritePacket(c.conn, &keepAlive); err != nil {
			c.Errorf("failed to write packet: %v", err)
			return
		}
		c.Debugf("Sent keep alive: %d", randId)

		// Wait for the keep alive packet to be received.
		select {
		case receivedId := <-c.keepAliveId:
			if receivedId != randId {
				c.Errorf("received keep alive id is not the same as the sent one")
				return
			}
			c.Debugf("Valid keep alive: %d - Ping %s", receivedId, time.Since(now))
			continue
		case <-time.NewTimer(15 * time.Second).C:
			c.Errorf("keep alive timeout")
			c.Close()
			return
		case <-c.close:
			return
		}
	}
}

// KeepAlive handles the keep alive packet.
func (c *Connection) KeepAlive(p *protocol.KeepAlive) error {
	c.Debugf("Received keep alive: %d", p.ID)
	c.keepAliveId <- p.ID
	return nil
}

// Login handles the login packet.
func (c *Connection) Login(p *protocol.Login) error {
	c.Debugf("Received login: %s", p.Token)
	return nil
}

// TunnelRequest handles the tunnel request packet.
func (c *Connection) TunnelRequest(p *protocol.TunnelRequest) error {
	c.Debugf("Received tunnel request: %s - %s", p.UUID, p.Protocol)
	return nil
}

// StreamRequest handles the stream request packet.
func (c *Connection) StreamRequest(p *protocol.StreamRequest) error {
	c.Debugf("Received stream request: %s - %s", p.Token, p.UUID)
	io.Copy(c.conn, c.conn)
	return nil
}

// Close closes the connection.
func (c *Connection) Close() error {
	c.Debugf("closing connection")
	close(c.close)
	return c.conn.Close()
}
