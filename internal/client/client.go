package client

import (
	"net"
	"sync"

	"github.com/corp4/gateway4/internal/protocol"
	log "github.com/sirupsen/logrus"
)

// Client is a client connection.
type Client struct {
	sync.Locker
	*log.Entry

	net.Conn
	close chan struct{}
}

// NewClient creates a new client from a connection.
func NewClient(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Locker: &sync.Mutex{},
		Entry:  log.WithField("client", conn.RemoteAddr().String()),
		Conn:   conn,
	}
	return client, nil
}

// KeepAlive is called when a keep alive packet is received.
func (c *Client) KeepAlive(p *protocol.KeepAlive) error {
	c.Debugf("Received keep alive: %d", p.ID)

	// Send a keep alive packet.
	keepAlive := protocol.KeepAlive{ID: p.ID}
	if err := protocol.WritePacket(c.Conn, &keepAlive); err != nil {
		c.Errorf("failed to write packet: %v", err)
		return err
	}
	return nil
}

// Disconnect is called when a disconnect packet is received.
func (c *Client) Disconnect(p *protocol.Disconnect) error {
	c.Debugf("Received disconnect: %s", p.Reason)
	c.Infof("Disconnecting client: %s", p.Reason)
	c.Close()
	return nil
}

// TunnelAccept is called when a tunnel accept packet is received.
func (c *Client) TunnelAccept(p *protocol.TunnelAccept) error {
	c.Debugf("Received tunnel accept: %s", p.Host)
	return nil
}

// NewStream is called when a new stream packet is received.
func (c *Client) NewStream(p *protocol.NewStream) error {
	c.Debugf("Received new stream: %s", p.UUID)
	return nil
}

// SteamAccept is called when a stream accept packet is received.
func (c *Client) StreamAccept(p *protocol.StreamAccept) error {
	c.Debugf("Received stream accept.")
	return nil
}

// Close closes the client connection.
func (c *Client) Close() error {
	close(c.close)
	return c.Conn.Close()
}
