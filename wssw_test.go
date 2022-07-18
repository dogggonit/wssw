package wssw

import (
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWSSW(t *testing.T) {
	w := createTestWebsocket[RawString](t)
	conn := createTestServer(t, w)

	testMessage := "hello world"
	if err := conn.WriteMessage(websocket.TextMessage, []byte(testMessage)); err != nil {
		t.Fatalf("%+v", err)
	}

	switch mt, m, err := conn.ReadMessage(); {
	case err != nil:
		t.Fatalf("%+v", err)
	case mt != websocket.TextMessage:
		t.Fatal("incorrect message type received")
	case string(m) != testMessage:
		t.Fatalf("received message does not match what was sent")
	}
}

func createTestServer(t *testing.T, h http.Handler) *websocket.Conn {
	t.Helper()

	s := httptest.NewServer(h)
	url := strings.ReplaceAll(strings.ReplaceAll(s.URL, "https", "wss"), "http", "ws")

	w, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(s.Close)
	t.Cleanup(func() {
		w.Close()
	})
	return w
}

func createTestWebsocket[T any](t *testing.T) Websocket[T, T] {
	t.Helper()

	w := New[T, T](&Config{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(*http.Request) bool {
				return true
			},
		},
	})

	w.OnError().Subscribe(func(err error) {
		t.Fatalf("%+v", err)
	})

	w.OnConnect().Subscribe(func(s Connection[T, T]) {
		s.OnError().Subscribe(func(err error) {
			t.Fatalf("%+v", err)
		})

		s.Subscribe(func(t T) {
			s.Publish(t)
		})
	})

	return w
}
