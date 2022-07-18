package main

import (
	"flag"
	"github.com/doggonit/wssw"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	host = flag.String("host", "localhost:8080", "host to serve on")
)

func init() {
	flag.Parse()
}

func main() {
	http.Handle("/", getEchoWS[wssw.RawString]())
	log.Printf("listening on %s", *host)
	log.Fatal(http.ListenAndServe(*host, nil))
}

func getEchoWS[T any]() wssw.Websocket[T, T] {
	w := wssw.New[T, T](&wssw.Config{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(*http.Request) bool {
				return true
			},
		},
	})

	w.OnError().Subscribe(func(err error) {
		log.Printf("err: %+v", err)
	})

	w.OnConnect().Subscribe(func(c wssw.Connection[T, T]) {
		log.Printf("[%s][%s] connected", c.GetConn().RemoteAddr().String(), c.GetID().String())

		c.OnError().Subscribe(func(e error) {
			log.Printf("[%s][%s] err: %+v", c.GetConn().RemoteAddr().String(), c.GetID().String(), e)
		})
		c.Subscribe(func(m T) {
			log.Printf(
				"[%s][%s] received: '%+v'",
				c.GetConn().RemoteAddr().String(),
				c.GetID(),
				m,
			)
			c.Publish(m)
		})
		c.OnClose().Subscribe(func(wssw.Connection[T, T]) {
			log.Printf("[%s][%s] closed", c.GetConn().RemoteAddr().String(), c.GetID().String())
		})
	})

	w.OnDisconnect().Subscribe(func(c wssw.Connection[T, T]) {
		log.Printf("[%s][%s] disconnected", c.GetConn().RemoteAddr().String(), c.GetID().String())
	})

	return w
}
