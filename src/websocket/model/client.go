package websocket_model

import "github.com/gofiber/contrib/websocket"

type Client[T any] struct {
	Connection *websocket.Conn
	Data       T
}

func CreateClient[T any](data T, conn *websocket.Conn) *Client[T] {
	return &Client[T]{
		Connection: conn,
		Data:       data,
	}
}
