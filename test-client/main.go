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

func get(hostAddr string) error {
	res, err := http.Get(hostAddr)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var tr TestResponse
	err = json.NewDecoder(res.Body).Decode(&tr)
	if err != nil {
		return err
	}
	log.Printf("%s: %d\n", tr.Name, tr.Cnt)
	return nil
}

func post(hostAddr string, inc int) error {
	reqBody, err := json.Marshal(&inc)
	if err != nil {
		return err
	}

	_, err = http.Post(hostAddr, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		panic("must provide host address")
	}

	hostAddr := os.Args[1]

	var choice int
	for {
		println("0: GET\n1: POST")
		fmt.Scanln(&choice)
		var err error
		switch choice {
		case 0:
			err = get(hostAddr)
		case 1:
			err = post(hostAddr, 1)
		}
		if err != nil {
			println(err.Error())
		}
	}
}
