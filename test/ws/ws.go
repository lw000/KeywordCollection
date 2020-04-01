package ws

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// TLSDialer ...
	TLSDialer = &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
	}
)

// FastWsClient ...
type FastWsClient struct {
	tag            int
	conn           *websocket.Conn
	onMessage      func([]byte)
	onConnected    func()
	onDisConnected func()
}

// Open ...
func (fc *FastWsClient) Open(scheme string, host, path string) error {
	u := url.URL{Scheme: scheme, Host: host, Path: path}
	var er error
	if scheme == "wss" {
		fc.conn, _, er = TLSDialer.Dial(u.String(), nil)
	} else if scheme == "ws" {
		fc.conn, _, er = websocket.DefaultDialer.Dial(u.String(), nil)
	} else {
		return errors.New("未知Scheme")
	}

	if er != nil {
		return er
	}

	fc.onConnected()

	return nil
}

// HandleConnected ...
func (fc *FastWsClient) HandleConnected(f func()) {
	fc.onConnected = f
}

// HandleDisConnected ...
func (fc *FastWsClient) HandleDisConnected(f func()) {
	fc.onDisConnected = f
}

// HandleMessage ...
func (fc *FastWsClient) HandleMessage(f func(data []byte)) {
	fc.onMessage = f
}

// SendMessage ...
func (fc *FastWsClient) SendMessage(data []byte) error {
	er := fc.conn.WriteMessage(websocket.TextMessage, data)
	if er != nil {
		log.Println(er)
		return er
	}
	return nil
}

// Ping ...
func (fc *FastWsClient) Ping() error {
	er := fc.conn.WriteMessage(websocket.PingMessage, []byte(""))
	if er != nil {
		log.Println(er)
		return er
	}
	return nil
}

// Run ...
func (fc *FastWsClient) Run() {
	for {
		_, message, err := fc.conn.ReadMessage()
		if err != nil {
			log.Println(err)

			fc.onDisConnected()

			return
		}
		fc.onMessage(message)
	}
}
