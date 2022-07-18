package wssw

import (
	"errors"
	"github.com/dogggonit/rxxr"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
)

type (
	Websocket[S, R any] interface {
		http.Handler
		// OnError allows subscribing to errors in the websocket handler.
		OnError() rxxr.Subscribable[error]
		// OnConnect allows subscribing to new connections.
		OnConnect() rxxr.Subscribable[Connection[S, R]]
		// OnDisconnect allows subscribing to connection disconnections.
		OnDisconnect() rxxr.Subscribable[Connection[S, R]]
	}

	Connection[S, R any] interface {
		rxxr.Subscribable[R]
		Publish(v ...S)

		// GetID gets an id that uniquely identifies this connection.
		GetID() uuid.UUID

		// GetConn gets the underlying websocket connection object.
		GetConn() *websocket.Conn
		// Close closes the connection
		Close()

		// OnError allows subscribing to errors that the connection encounters.
		OnError() rxxr.Subscribable[error]
		// OnClose allows subscribing to the websocket closing.
		// After the websocket has been closed, all subscribers will immediately get a value when subscribing.
		OnClose() rxxr.Subscribable[Connection[S, R]]
	}
)

type (
	// Marshaller allows a type to define how it marshals itself into a websocket message.
	// Possible message types are websocket.TextMessage and websocket.BinaryMessage.
	Marshaller interface {
		MarshalWS() (b []byte, mt int, err error)
	}

	// Unmarshall allow a type to define how it unmarshalls itself from a websocket message.
	// Possible message types are websocket.TextMessage and websocket.BinaryMessage.
	Unmarshall interface {
		UnmarshallWS(mt int, b []byte) error
	}
)

var (
	// ErrWrongMessageType should be returned from Unmarshall when the wrong message type is encountered.
	// Possible message types are websocket.TextMessage and websocket.BinaryMessage.
	ErrWrongMessageType = errors.New("incorrect message type")
)
