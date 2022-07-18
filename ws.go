package wssw

import (
	"github.com/dogggonit/rxxr"
	"github.com/gorilla/websocket"
	"net/http"
)

type (
	ws[S, R any] struct {
		config       Config
		errors       rxxr.Pipe[error]
		onConnect    rxxr.Pipe[Connection[S, R]]
		onDisconnect rxxr.Pipe[Connection[S, R]]
	}

	// Config allows configuration of the
	Config struct {
		// If the type being sent/received implements Marshaller or Unmarshall this will be ignored,
		// by default the websocket will convert your type to JSON. Enabling this option will
		// indent the JSON being sent.
		IndentJSON bool
		// This allows you to change the settings used when upgrading a connection to a websocket.
		Upgrader websocket.Upgrader
	}
)

// New creates a new websocket handler.
// S is the type that will be sent from the server to the websocket client.
// R is the type that will be received from the websocket client and used by your program.
// If conf is nil the default configuration will be used.
func New[S, R any](conf *Config) Websocket[S, R] {
	if conf == nil {
		conf = &Config{}
	}

	w := &ws[S, R]{
		config:       *conf,
		errors:       rxxr.New[error](nil),
		onConnect:    rxxr.New[Connection[S, R]](nil),
		onDisconnect: rxxr.New[Connection[S, R]](nil),
	}
	return w
}

func (ws *ws[S, R]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.config.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		ws.errors.Publish(err)
		return
	}

	c := newConnection(ws, conn)
	defer c.Close()

	ws.onConnect.Publish(c)
	defer ws.onDisconnect.Publish(c)
	c.holdConnectionOpen()
}

func (ws *ws[S, R]) OnError() rxxr.Subscribable[error] {
	return ws.errors
}

func (ws *ws[S, R]) OnConnect() rxxr.Subscribable[Connection[S, R]] {
	return ws.onConnect
}

func (ws *ws[S, R]) OnDisconnect() rxxr.Subscribable[Connection[S, R]] {
	return ws.onDisconnect
}
