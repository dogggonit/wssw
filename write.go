package wssw

import (
	"encoding/json"
	"github.com/dogggonit/rxxr"
	"github.com/gorilla/websocket"
)

func write[T any](conn *websocket.Conn, indentJSON bool, errors rxxr.Pipe[error]) func(T) {
	return func(t T) {
		if pm, ok := any(t).(*websocket.PreparedMessage); ok {
			if err := conn.WritePreparedMessage(pm); err != nil {
				errors.Publish(err)
			}
		} else if b, m, err := prepareMessage(t, indentJSON); err != nil {
			errors.Publish(err)
		} else if err = conn.WriteMessage(m, b); err != nil {
			errors.Publish(err)
		}
	}
}

func prepareMessage[T any](v T, indentJSON bool) ([]byte, int, error) {
	var ptr any
	if _, isPtr := create[T](); isPtr {
		ptr = v
	} else {
		ptr = &v
	}

	if m, ok := ptr.(Marshaller); ok {
		return encodeMarshaller(m)
	}
	return encodeAny(v, indentJSON)
}

func encodeAny(v any, indentJSON bool) ([]byte, int, error) {
	var marshaller func(v any) ([]byte, error)
	if indentJSON {
		marshaller = func(v any) ([]byte, error) {
			return json.MarshalIndent(v, "", "    ")
		}
	} else {
		marshaller = json.Marshal
	}

	b, err := marshaller(v)
	if err != nil {
		return nil, 0, err
	}
	return b, websocket.TextMessage, nil
}

func encodeMarshaller(m Marshaller) ([]byte, int, error) {
	b, mt, err := m.MarshalWS()
	if err != nil {
		return nil, 0, err
	}
	return b, mt, nil
}
