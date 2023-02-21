package protocol

import "reflect"

var (
	// CodeTypeMap maps the packet type code to the packet type.
	CodeTypeMap map[uint8]reflect.Type

	// TypeCodeMap maps the packet type to the packet type code.
	TypeCodeMap map[reflect.Type]uint8
)

func init() {
	// Build the map with the packet types.
	CodeTypeMap = map[uint8]reflect.Type{
		0x00: reflect.TypeOf(Login{}),
		0x01: reflect.TypeOf(KeepAlive{}),
		0x02: reflect.TypeOf(Disconnect{}),
		0x03: reflect.TypeOf(TunnelRequest{}),
		0x04: reflect.TypeOf(TunnelAccept{}),
		0x05: reflect.TypeOf(NewStream{}),
		0x06: reflect.TypeOf(StreamRequest{}),
		0x07: reflect.TypeOf(StreamAccept{}),
	}

	// Build the reverse map.
	TypeCodeMap = map[reflect.Type]uint8{}
	for code, typ := range CodeTypeMap {
		TypeCodeMap[typ] = code
	}
}

// Login is sent by the client to the server to authenticate.
type Login struct {
	Token string // The authentication token.
}

// KeepAlive is sent by to server to check client connection regularly.
type KeepAlive struct {
	ID uint64 // Randomly generated ID.
}

// Disconnect is sent by the server to the client to disconnect with a reason.
type Disconnect struct {
	Reason string // The reason of the disconnection.
}

// TunnelRequest is sent by the client to the server to request a tunnel.
type TunnelRequest struct {
	UUID     string // The UUID of the tunnel. (randomly generated)
	Protocol string // The protocol of the tunnel (e.g. tcp, http, etc.)
}

// TunnelAccept is sent by the server to the client to accept a tunnel request.
type TunnelAccept struct {
	UUID string // The UUID of the tunnel.
	Host string // The host of the tunnel.
}

// NewStream is sent by the server to notify the client that a new stream is available.
type NewStream struct {
	UUID string // The UUID of the new stream. (randomly generated)
}

// StreamRequest is sent by the client to the server to open a new stream.
type StreamRequest struct {
	Token string // The authentication token.
	UUID  string // The UUID of the stream to open.
}

// StreamAccept is sent by the server to the client to accept a stream request.
type StreamAccept struct {
}

// FromClient interface is implemented by the server.
type FromClient interface {
	Login(*Login) error
	KeepAlive(*KeepAlive) error
	TunnelRequest(*TunnelRequest) error
	StreamRequest(*StreamRequest) error
}

// FromServer interface is implemented by the client.
type FromServer interface {
	KeepAlive(*KeepAlive) error
	Disconnect(*Disconnect) error
	TunnelAccept(*TunnelAccept) error
	NewStream(*NewStream) error
	StreamAccept(*StreamAccept) error
}
