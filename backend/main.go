package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	nats "github.com/nats-io/go-nats"
)

/*
 person ...
*/
type person struct {
	Name       string
	Age        int
	Occupation string
}

var wg sync.WaitGroup

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	fmt.Println("Hey Bosch")

	go sendNats()
	go recvNats()

	http.HandleFunc("/person", index)
	http.HandleFunc("/ws", socket)

	log.Fatal(http.ListenAndServe(":7654", nil))
}

func recvNats() {
	log.Println("recv nats start")
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		panic(err)
	}

	nc.Subscribe("hmi", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})

	wg.Wait()
}

func sendNats() {
	log.Println("send nats start")
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		panic(err)
	}

	for {
		nc.Publish("hmi", []byte("Hello World hmi"))
		time.Sleep(time.Millisecond * 1000)
	}
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
