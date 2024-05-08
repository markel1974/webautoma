package bidi

import (
	"context"
	"fmt"
	"net/http"

	"nhooyr.io/websocket"
)

type WebSocket struct {
	conn *websocket.Conn
}

func (ws *WebSocket) Connect(ctx context.Context, wsURL string, header http.Header) error {
	if ws.conn != nil {
		return fmt.Errorf("duplicated connection: %s", wsURL)
	}

	c, _, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{ //nolint bodyclose
		HTTPHeader: header,
	})
	if err != nil {
		return err
	}

	ws.conn = c

	return nil
}

func (ws *WebSocket) Close() error {
	return ws.conn.Close(websocket.StatusNormalClosure, "")
}

func (ws *WebSocket) Write(ctx context.Context, p []byte) error {
	return ws.conn.Write(ctx, websocket.MessageText, p)
}

func (ws *WebSocket) Read(ctx context.Context) ([]byte, error) {
	typ, p, err := ws.conn.Read(ctx)
	if err != nil {
		return nil, err
	}

	if typ != websocket.MessageText {
		return nil, fmt.Errorf("unsupported messageType: %s", typ)
	}

	return p, nil
}
