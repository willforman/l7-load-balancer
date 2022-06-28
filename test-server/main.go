package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type TestResponse struct {
	Name string `json:"name"`
	Cnt  int    `json:"cnt"`
}

func handle(tr *TestResponse, sleepSeconds *int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("hit!")
		time.Sleep(time.Second * time.Duration(*sleepSeconds))
		switch r.Method {
		case http.MethodGet:
			json.NewEncoder(w).Encode(tr)
		case http.MethodPost:
			var inc int
			json.NewDecoder(r.Body).Decode(&inc)
			tr.Cnt += inc
		}
	}
}

func main() {
	if len(os.Args) != 3 {
		panic("must provide exactly 1 argument")
	}
	portStr := os.Args[1]
	sleepSeconds, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	tr := TestResponse{portStr, 0}

	http.HandleFunc("/api", handle(&tr, &sleepSeconds))
	log.Fatal(http.ListenAndServe(portStr, nil))
}
