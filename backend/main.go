package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

/*
 person ...
*/
type person struct {
	Name       string
	Age        int
	Occupation string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	fmt.Println("Hey Bosch")

	http.HandleFunc("/person", index)
	http.HandleFunc("/ws", socket)

	// log.Println(http.ListenAndServeTLS(":7654", "cert.pem", "key.pem", nil))
	log.Fatal(http.ListenAndServe(":7654", nil))
}

func socket(w http.ResponseWriter, r *http.Request) {

	log.Println("websocket call")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("websocket connection open")

	for {
		messageType, recvMsg, err := conn.ReadMessage()
		if err != nil {
			log.Println("recv from socket:", recvMsg)
			break
		}
		if err := conn.WriteMessage(messageType, []byte("pong")); err != nil {
			log.Println("error on write to websocket")
			break
		}
		log.Println("answer websocket successful")
	}

}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("rest call")

	p := person{
		"Fedja",
		26,
		"Software Developer",
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(p); err != nil {
		panic(err)
	}
}
