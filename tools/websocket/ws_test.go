package main

import (
	"flag"
	"net/http"
	"github.com/GrFrHuang/gox/log"
	"github.com/gorilla/websocket"
	"html/template"
	"testing"
	"fmt"
)

var serverAddr = flag.String("serverAddr", "localhost:7080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func TestWebSocket(t *testing.T) {
	flag.Parse()
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*serverAddr, nil))
}

func echo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=====", r.Header)
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("[ws]: ", err)
		return
	}
	defer c.Close()
	for {
		// ReadMessage function is blocking until receive message from client.
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Error("[ws]: ", err)
			break
		}
		log.Info("recv: ", message)
		msg := string(message) + "GrFrHuang"
		err = c.WriteMessage(mt, []byte(msg))
		if err != nil {
			log.Error("[ws]: ", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server,
"Send" to send a message to the server and "Close" to close the connection.
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello ">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
