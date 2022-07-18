package wssw

import "github.com/gorilla/websocket"

// RawString is a type used to help with sending/receiving basic strings from a websocket connection.
type RawString string

func (rs *RawString) MarshalWS() (b []byte, mt int, err error) {
	return []byte(*rs), websocket.TextMessage, nil
}

func (rs *RawString) UnmarshallWS(mt int, b []byte) error {
	if mt != websocket.TextMessage {
		return ErrWrongMessageType
	}
	*rs = RawString(b)
	return nil
}

// Value returns this value as a string.
func (rs *RawString) Value() string {
	return string(*rs)
}

// RawBytes is a type used to help send/receive binary messages from a websocket connection.
type RawBytes []byte

func (rb *RawBytes) MarshalWS() (b []byte, mt int, err error) {
	return *rb, websocket.BinaryMessage, nil
}

func (rb *RawBytes) UnmarshallWS(mt int, b []byte) error {
	if mt != websocket.BinaryMessage {
		return ErrWrongMessageType
	}
	*rb = b
	return nil
}

// Value returns this value as a slice of bytes
func (rb *RawBytes) Value() []byte {
	return *rb
}

// Raw is used to send/receive raw messages from a websocket connection.
type Raw struct {
	MessageType int
	Data        []byte
}

func (rb *Raw) MarshalWS() (b []byte, mt int, err error) {
	return rb.Data, rb.MessageType, nil
}

func (rb *Raw) UnmarshallWS(mt int, b []byte) error {
	rb.MessageType = mt
	rb.Data = b
	return nil
}

// IsString returns true if the message type is a string.
func (rb *Raw) IsString() bool {
	return rb.MessageType == websocket.TextMessage
}

// IsBytes returns true if the message type is a binary message.
func (rb *Raw) IsBytes() bool {
	return rb.MessageType == websocket.BinaryMessage
}

// Bytes returns this value as a slice of bytes.
func (rb *Raw) Bytes() []byte {
	return rb.Data
}

// String returns this value as a string.
func (rb *Raw) String() string {
	return string(rb.Data)
}
