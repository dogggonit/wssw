package wssw

import (
	"github.com/dogggonit/rxxr"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sync"
)

type connection[S, R any] struct {
	id uuid.UUID

	onClose     rxxr.Pipe[Connection[S, R]]
	open        bool
	openLock    sync.RWMutex
	sendPipe    rxxr.Pipe[S]
	receivePipe rxxr.Pipe[R]

	conn   *websocket.Conn
	errors rxxr.Pipe[error]
}

func newConnection[S, R any](w *ws[S, R], conn *websocket.Conn) *connection[S, R] {
	c := &connection[S, R]{
		id:          uuid.New(),
		sendPipe:    rxxr.New[S](nil),
		receivePipe: rxxr.New[R](nil),
		conn:        conn,
		open:        true,
		errors:      rxxr.New[error](nil),
		onClose:     new(rxxr.Config[Connection[S, R]]).SendOnSubscribe(true).Build(),
	}
	c.sendPipe.Subscribe(write[S](conn, w.config.IndentJSON, c.errors))
	return c
}

func (c *connection[S, R]) holdConnectionOpen() {
	for {
		if v, err, closed := read[R](c.conn); err != nil {
			c.errors.Publish(err)
			return
		} else if closed {
			c.Close()
			return
		} else {
			c.receivePipe.Publish(v)
		}
	}
}

func (c *connection[S, R]) Subscribe(fn func(R)) rxxr.Subscription {
	c.openLock.RLock()
	defer c.openLock.RUnlock()
	if c.open {
		return c.receivePipe.Subscribe(fn)
	}
	return nil
}

func (c *connection[S, R]) Unsubscribe(subscription rxxr.Subscription) {
	if subscription == nil {
		return
	}
	c.openLock.RLock()
	defer c.openLock.RUnlock()
	if c.open {
		c.receivePipe.Unsubscribe(subscription)
	}
}

func (c *connection[S, R]) Publish(v ...S) {
	c.openLock.RLock()
	defer c.openLock.RUnlock()
	if c.open {
		c.sendPipe.Publish(v...)
	}
}

func (c *connection[S, R]) GetID() uuid.UUID {
	return c.id
}

func (c *connection[S, R]) GetConn() *websocket.Conn {
	return c.conn
}

func (c *connection[S, R]) Close() {
	c.openLock.Lock()
	defer c.openLock.Unlock()
	if !c.open {
		return
	}

	c.onClose.Publish(c)
	c.onClose.Close()
	c.sendPipe.Close()
	c.receivePipe.Close()
	c.errors.Close()

	_ = c.conn.Close()
	c.open = false
}

func (c *connection[S, R]) OnError() rxxr.Subscribable[error] {
	return c.errors
}

func (c *connection[S, R]) OnClose() rxxr.Subscribable[Connection[S, R]] {
	return (*closeChecker[S, R])(c)
}

type closeChecker[S, R any] connection[S, R]

func (c *closeChecker[S, R]) GetID() uuid.UUID {
	return uuid.Nil
}

func (c *closeChecker[S, R]) Subscribed() bool {
	return false
}

func (c *closeChecker[S, R]) Subscribe(f func(Connection[S, R])) rxxr.Subscription {
	c.openLock.RLock()
	defer c.openLock.RUnlock()
	if c.open {
		return c.onClose.Subscribe(f)
	} else if f != nil {
		f((*connection[S, R])(c))
	}
	return c
}

func (c *closeChecker[S, R]) Unsubscribe(s rxxr.Subscription) {
	c.openLock.RLock()
	defer c.openLock.RUnlock()
	if c.open {
		c.onClose.Unsubscribe(s)
	}
}
