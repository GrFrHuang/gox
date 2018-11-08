// WebSocket protocol ws://xxx.xxx:xx/xxx
// WebSocket is a two-way communication protocols, after the connection is established,
// the WebSocket server and Browser can be active to send or receive data from each other,
// like a Socket, the difference is the WebSocket is a kind of a simple simulation based on the Web Socket protocol.

// The characteristics of WebSocket:
// -Based on TCP protocol.
// -Good compatibility with HTTP protocol. The default ports are also 80 and 443, and the handshake phase USES HTTP protocol.
// -Data format is light, performance cost is low, communication efficiency is high.
// -Can send text or binary data.
// -There is no same-origin restriction, and the client can communicate with any server.
// -The protocol identifier is ws (or WSS if encrypted) and the server URL is the URL.

// 1.Browser与WebSocket服务器通过TCP三次握手建立连接，如果这个建立连接失败，那么后面的过程就不会执行，Web应用程序将收到错误消息通知。
// 2.在TCP建立连接成功后，Browser/UA通过http协议传送WebSocket支持的版本号，协议的字版本号，原始地址，主机地址等等一些列字段给服务器端。
// 3.WebSocket服务器收到Browser/UA发送来的握手请求后，如果数据包数据和格式正确，客户端和服务器端的协议版本号匹配等等，
// 就接受本次握手连接，并给出相应的数据回复status code为101，同样回复的数据包也是采用http协议传输。
// 4.Browser收到服务器回复的数据包后，如果数据包内容、格式都没有问题的话，就表示本次连接成功，触发onopen消息，
// 此时Web开发者就可以在此时通过send接口想服务器发送数据。否则，握手连接失败，Web应用程序会收到onerror消息，并且能知道连接失败的原因。

package main

import (
	"github.com/gorilla/websocket"
	"github.com/GrFrHuang/gox/log"
	"net/http"
	"encoding/json"
	"encoding/xml"
	"time"
)

type WebSocketServer struct {
	conn          *websocket.Conn
	request       *http.Request
	response      http.ResponseWriter
	readDeadLine  time.Time // Read timeout time.
	writeDeadLine time.Time // Write timeout time.
	//authKey       string    // Auth header field.
}

// todo safe many goroutines
// Initialize a web socket server.
func NewWebSocketServer(w http.ResponseWriter, r *http.Request) (*WebSocketServer) {
	// Use default options.
	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("[ws]: ", err)
	}
	return &WebSocketServer{
		conn:     c,
		request:  r,
		response: w,
	}
}

func (ws *WebSocketServer) VerifyAuth(msgType int, message interface{}) (err error) {
	var data []byte
	data, err = json.Marshal(message)
	if err != nil {
		log.Error("[ws]: ", err)
		return
	}
	err = ws.conn.WriteControl(msgType, data, ws.writeDeadLine)
	if err != nil {
		log.Error("[ws]: ", err)
	}
	return
}

// Start async listen message from client.
func (ws *WebSocketServer) StartRead() (err error) {
	go func() {
		for {
			// ReadMessage function is blocking until receive message from client.
			msgType, message, err := ws.conn.ReadMessage()
			if err != nil {
				log.Error("[ws]: ", err)
				break
			}
			log.Info("recv: ", message)
			msg := string(message) + "GrFrHuang"
			err = ws.conn.WriteMessage(msgType, []byte(msg))
			if err != nil {
				log.Error("[ws]: ", err)
				break
			}
		}
		defer func() {
			err2 := ws.conn.Close()
			if err2 != nil {
				log.Error(err2)
			}
		}()
	}()
	return
}

// Send response message to client by json format.
func (ws *WebSocketServer) WriteJson(msgType int, message interface{}) (err error) {
	var data []byte
	data, err = json.Marshal(message)
	if err != nil {
		log.Error("[ws]: ", err)
		return
	}
	err = ws.conn.WriteControl(msgType, data, ws.writeDeadLine)
	if err != nil {
		log.Error("[ws]: ", err)
	}
	return
}

// Send response message to client by xml format.
func (ws *WebSocketServer) WriteXml(msgType int, message interface{}) (err error) {
	var data []byte
	data, err = xml.Marshal(message)
	if err != nil {
		log.Error("[ws]: ", err)
		return
	}
	err = ws.conn.WriteControl(msgType, data, ws.writeDeadLine)
	if err != nil {
		log.Error("[ws]: ", err)
	}
	return
}

func (ws *WebSocketServer) SetWriteDeadLine(line time.Time) () {
	if line.Before(time.Now()) || line.Equal(time.Now()) {
		log.Panic("[ws]: time error !")
	}
	ws.writeDeadLine = line
}

func (ws *WebSocketServer) SetReadDeadLine(line time.Time) () {
	if line.Before(time.Now()) || line.Equal(time.Now()) {
		log.Panic("[ws]: time error !")
	}
	ws.readDeadLine = line
}

// todo connect pool manage
// todo 通过设置单IP可建立连接的最大连接数的方式防范
// todo 通过设置auth_token
