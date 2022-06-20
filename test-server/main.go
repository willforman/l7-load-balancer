package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type TestResponse struct {
	Name string `json:"name"`
	Cnt  int    `json:"cnt"`
}

func handle(tr *TestResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("hit!")
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
	if len(os.Args) != 2 {
		panic("must provide exactly 1 argument")
	}
	portStr := os.Args[1]
	tr := TestResponse{portStr, 0}

	http.HandleFunc("/api", handle(&tr))
	log.Fatal(http.ListenAndServe(portStr, nil))
}
