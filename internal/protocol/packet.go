package protocol

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

// A packet is a sequence of bytes with the following format:
//
// <length><packet_type><payload>
//
// The length is a 64-bit unsigned integer in big-endian format.
// The packet type is a 8-bit unsigned integer, which determine the type of the packet.
// The payload is a sequence of bytes with the length specified by the length field - 1.

// Length is the length of the packet.
type Length uint64

// ReadLength reads the length of the packet from the reader.
func ReadLength(reader io.Reader) (Length, error) {
	var length Length
	if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
		return length, err
	}
	return length, nil
}

// WriteLength writes the length of the packet to the writer.
func WriteLength(writer io.Writer, length Length) error {
	return binary.Write(writer, binary.BigEndian, length)
}

// Type is the type of the packet.
type Type uint8

// ReadType reads the type of the packet from the reader.
func ReadType(reader io.Reader) (Type, error) {
	var t Type
	if err := binary.Read(reader, binary.BigEndian, &t); err != nil {
		return t, err
	}
	return t, nil
}

// WriteType writes the type of the packet to the writer.
func WriteType(writer io.Writer, t Type) error {
	return binary.Write(writer, binary.BigEndian, t)
}

// Retrieve the type of the packet.
func (t *Type) PacketType() (reflect.Type, error) {
	packetType, ok := CodeTypeMap[uint8(*t)]
	if !ok {
		return nil, fmt.Errorf("unknown packet type: %v", *t)
	}
	return packetType, nil
}

// Payload is the payload of the packet.
type Payload []byte

// ReadPayload reads the payload of the packet from the reader.
// The length of the payload is specified by the length parameter.
// Do not subtract 1 from the length parameter; pass the length of
// the packet as is.
func ReadPayload(reader io.Reader, length Length) (Payload, error) {
	payload := make([]byte, length-1)
	if _, err := io.ReadFull(reader, payload); err != nil {
		return payload, err
	}
	return payload, nil
}

// ReadPacket reads the packet from the reader.
func ReadPacket(reader io.Reader) (*interface{}, error) {
	length, err := ReadLength(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read packet length: %v", err)
	}

	t, err := ReadType(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read packet type: %v", err)
	}

	payload, err := ReadPayload(reader, length)
	if err != nil {
		return nil, fmt.Errorf("failed to read packet payload: %v", err)
	}

	packetType, err := t.PacketType()
	if err != nil {
		return nil, fmt.Errorf("failed to get packet type: %v", err)
	}

	packet := reflect.New(packetType).Interface()
	if err := json.Unmarshal(payload, &packet); err != nil {
		return nil, fmt.Errorf("failed to unmarshal packet: %v", err)
	}

	return &packet, nil
}

// WritePacket writes the packet to the writer.
func WritePacket(writer io.Writer, packet interface{}) error {
	payload, err := json.Marshal(packet)
	if err != nil {
		return fmt.Errorf("failed to marshal packet: %v", err)
	}

	length := Length(len(payload) + 1)
	if err := WriteLength(writer, length); err != nil {
		return fmt.Errorf("failed to write packet length: %v", err)
	}

	t := Type(TypeCodeMap[reflect.TypeOf(packet).Elem()])
	if err := WriteType(writer, t); err != nil {
		return fmt.Errorf("failed to write packet type: %v", err)
	}

	if _, err := writer.Write(payload); err != nil {
		return fmt.Errorf("failed to write packet payload: %v", err)
	}

	return nil
}

// ReadServer reads the server packet from the reader.
func ReadServer(src FromServer, reader io.Reader) error {
	packet, err := ReadPacket(reader)
	if err != nil {
		return err
	}

	switch p := (*packet).(type) {
	case *KeepAlive:
		return src.KeepAlive(p)
	case *Disconnect:
		return src.Disconnect(p)
	case *TunnelAccept:
		return src.TunnelAccept(p)
	case *NewStream:
		return src.NewStream(p)
	case *StreamAccept:
		return src.StreamAccept(p)
	}

	return fmt.Errorf("unknown packet type: %v", packet)
}

// ReadClient reads the client packet from the reader.
func ReadClient(src FromClient, reader io.Reader) error {
	packet, err := ReadPacket(reader)
	if err != nil {
		return err
	}

	switch p := (*packet).(type) {
	case *Login:
		return src.Login(p)
	case *KeepAlive:
		return src.KeepAlive(p)
	case *TunnelRequest:
		return src.TunnelRequest(p)
	case *StreamRequest:
		return src.StreamRequest(p)
	}

	return fmt.Errorf("unknown packet type: %v", packet)
}
