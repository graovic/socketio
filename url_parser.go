package socketio

import (
	"fmt"
	"net/url"
	"strings"
)

// Parse raw url string and make valid handshake or websockets socket.io url.
type UrlParser struct {
	raw    string
	parsed *url.URL
}

func NewUrlParser(raw string) (*UrlParser, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	if parsed.Scheme == "" {
		parsed.Scheme = "http"
	}
	return &UrlParser{raw: raw, parsed: parsed}, nil
}

func (u *UrlParser) Handshake() string {
	return fmt.Sprintf("%s/socket.io/1", u.parsed.String())
}

func (u *UrlParser) Websocket(sessionId string) string {
	var host string
	if u.parsed.Scheme == "https" {
		host = strings.Replace(u.parsed.String(), "https://", "wss://", 1)
	} else {
		host = strings.Replace(u.parsed.String(), "http://", "ws://", 1)
	}
	return fmt.Sprintf("%s/socket.io/1/websocket/%s", host, sessionId)
}