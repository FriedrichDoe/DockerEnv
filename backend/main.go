package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

/*
 person ...
*/
type person struct {
	Name       string
	Age        int
	Occupation string
}

func main() {
	fmt.Println("Hey Bosch")
	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe(":7654", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
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
