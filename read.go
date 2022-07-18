package wssw

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"reflect"
)

func read[T any](r *websocket.Conn) (v T, err error, closed bool) {
	if mt, b, e := r.ReadMessage(); e != nil {
		closed = true
	} else {
		n, isPtr := create[T]()
		if isPtr {
			err = decode(mt, b, n)
		} else {
			err = decode(mt, b, &n)
		}
		v = n
	}
	return
}

func decode(mt int, b []byte, v any) error {
	if m, ok := v.(Unmarshall); ok {
		return m.UnmarshallWS(mt, b)
	}
	return json.Unmarshal(b, v)
}

func create[T any]() (T, bool) {
	if typ := reflect.TypeOf(*new(T)); typ.Kind() == reflect.Pointer {
		return reflect.New(typ.Elem()).Interface().(T), true
	} else {
		return *new(T), false
	}
}
