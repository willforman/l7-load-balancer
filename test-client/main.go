package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type TestResponse struct {
	Name string `json:"name"`
	Cnt  int    `json:"cnt"`
}

func get(hostAddr string) {
	res, err := http.Get(hostAddr)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var tr TestResponse
	err = json.NewDecoder(res.Body).Decode(&tr)
	if err != nil {
		panic(err)
	}
	log.Printf("%s: %d\n", tr.Name, tr.Cnt)
}

func post(hostAddr string, inc int) {
	reqBody, err := json.Marshal(&inc)
	if err != nil {
		panic(err)
	}

	_, err = http.Post(hostAddr, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) != 2 {
		panic("Must provide host address")
	}

	hostAddr := os.Args[1]

	var choice int
	for {
		println("0: GET\n1: POST")
		fmt.Scanln(&choice)
		switch choice {
		case 0:
			get(hostAddr)
		case 1:
			post(hostAddr, 10)
		}
	}
}
