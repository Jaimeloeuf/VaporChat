package server

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var websocketUpgrader = websocket.Upgrader{
	// @todo
	// Buffers can be tuned based on average payload size and expected
	// concurrent connections, to find the sweet spot between smaller buffers
	// for lower per connection memory use and larger buffers for less syscalls
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// @todo
	// Potentially enable pool to reuse memory across connections, and to allow
	// for higher buffer sizes, potentially the default 4096
	// WriteBufferPool: &sync.Pool{},

	// @todo Allow connections from any origin for testing
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
