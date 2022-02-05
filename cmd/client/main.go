package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	go worker("01")
	go worker("02")
	go worker("03")
	go worker("04")
	go worker("05")

	time.Sleep(time.Hour * 1)
}

func worker(id string) {
	i := 0
	for {
		rand.Seed(time.Now().UnixNano())
		randomInt := rand.Intn(1000)

		log.Println(i, "sleeping for ms:", randomInt)
		time.Sleep(time.Millisecond * time.Duration(randomInt))

		input := fmt.Sprintf("worker-no=%s,message-no=%d", id, i)

		url := "http://127.0.0.1:8000"

		res, err := http.Post(url, "string", bytes.NewBuffer([]byte(input)))
		if err != nil {
			panic(err)
		}

		log.Println(fmt.Sprintf("%d status, %s", res.StatusCode, input))

		i++
	}
}
