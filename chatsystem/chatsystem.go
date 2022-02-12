package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type User struct {
	username   string
	connection *websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var savedConns = make(map[net.Addr]*User)
var savedDay = time.Now().Day()

func handleEcho(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// Check if message is initialization message
		splitMsg := strings.Split(string(msg), " ")
		if splitMsg[0] == "Init" {
			// Saves connection
			savedConns[conn.RemoteAddr()] = &User{username: strings.Join(splitMsg[1:], " "), connection: conn}
			conn.WriteMessage(msgType, getMessages())
		} else {
			currentUsername := savedConns[conn.RemoteAddr()].username
			// Prints remote address and message to console
			fmt.Printf("%s sent: %s\n", currentUsername, string(msg))
			addMessageToFile(currentUsername + ": " + string(msg))

			// Sends message to all connections
			for _, savedConn := range savedConns {
				savedConn.connection.WriteMessage(msgType, append([]byte(currentUsername+": "), msg...))
			}
		}
	}
}

func addMessageToFile(message string) {
	messagesFile, _ := os.OpenFile("messages.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	messagesFile.WriteString("|" + message)
}

func getMessages() []byte {
	data, _ := os.ReadFile("messages.txt")
	return data
}

func run() {
	for {
		currentDay := time.Now().Day()
		if currentDay != savedDay {
			savedDay = currentDay
			os.Truncate("messages.txt", 0)
		}

		time.Sleep(time.Second * 5)
	}
}

func main() {
	go run()

	http.HandleFunc("/echo", handleEcho)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "chatsystem.html")
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":8080", nil)
}
