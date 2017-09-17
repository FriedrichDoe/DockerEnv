package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	pb "workspace/DockerEnv/backend/tutorial"

	"github.com/golang/protobuf/proto"
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
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		panic(err)
	}

	book2 := &pb.AddressBook{}

	nc.Subscribe("hmi", func(m *nats.Msg) {
		if err := proto.Unmarshal(m.Data, book2); err != nil {
			log.Fatalln("lol ->", err)
		}

		bs, err1 := json.Marshal(book2)
		if err1 != nil {
			panic(err1)
		}
	})

	wg.Wait()
}

func sendNats() {
	log.Println("send nats start")
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		panic(err)
	}

	p1 := &pb.Person{
		Id:    1235123515,
		Name:  "Fedja Doe",
		Email: "fdoe@spassst.com",
		Phones: []*pb.Person_PhoneNumber{
			{Number: "555-4321", Type: pb.Person_WORK},
		},
	}

	p2 := &pb.Person{
		Id:    87967707,
		Name:  "deine Mutter",
		Email: "mudda@deine.org",
		Phones: []*pb.Person_PhoneNumber{
			{Number: "6666-12345", Type: pb.Person_HOME},
		},
	}

	book := &pb.AddressBook{
		People: []*pb.Person{
			p1,
			p2,
		},
	}

	out, err := proto.Marshal(book)
	if err != nil {
		log.Fatalln("failed to encode book", err)
	}

	for {
		nc.Publish("hmi", out)
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
