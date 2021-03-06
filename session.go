package socketio

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"log"
	"fmt"
)

// Session holds the configuration variables received from the socket.io
// server.
type Session struct {
	ID                 string
	HeartbeatTimeout   time.Duration
	ConnectionTimeout  time.Duration
	SupportedProtocols []string
}

// NewSession receives the configuraiton variables from the socket.io
// server.
func NewSession(url string) (*Session, error) {
	urlParser, err := newURLParser(url)
	if err != nil {
		return nil, err
	}
	log.Println("url Parser ", urlParser)
	response, err := http.Get(urlParser.handshake())
	if err != nil {
		return nil, err
	}
	log.Println("response , ", response)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	response.Body.Close()
	
	log.Println(string(body))
	sessionVars := strings.Split(string(body), ":")
	if len(sessionVars) != 4 {
		log.Println("Session vars ", sessionVars)
		return nil, errors.New(fmt.Sprintf("%s %d ", "Session variables is not 4, we received session with id ", len(sessionVars)))
	}

	id := sessionVars[0]

	heartbeatTimeoutSec, _ := strconv.Atoi(sessionVars[1])
	connectionTimeoutSec, _ := strconv.Atoi(sessionVars[2])

	heartbeatTimeout := time.Duration(heartbeatTimeoutSec) * time.Second
	connectionTimeout := time.Duration(connectionTimeoutSec) * time.Second

	supportedProtocols := strings.Split(string(sessionVars[3]), ",")

	return &Session{id, heartbeatTimeout, connectionTimeout, supportedProtocols}, nil
}

// SupportProtocol checks if the given protocol is supported by the
// socket.io server.
func (session *Session) SupportProtocol(protocol string) bool {
	for _, supportedProtocol := range session.SupportedProtocols {
		if protocol == supportedProtocol {
			return true
		}
	}
	return false
}
