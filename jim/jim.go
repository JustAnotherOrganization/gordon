package jim

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// ConnectSocket creates a websocket connection using the given address.
func ConnectSocket(address string) (*websocket.Conn, error) {
	var wsDialer websocket.Dialer
	conn, _, err := wsDialer.Dial(address, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "wsDialer.Dial : %s", address)
	}

	return conn, err
}

// ReadSocketMessage reads a websocket and returns the next message.
func ReadSocketMessage(conn *websocket.Conn) ([]byte, error) {
	_, byt, err := conn.ReadMessage()
	if err != nil {
		return nil, errors.Wrap(err, "conn.ReadMessage")
	}

	return byt, err
}
