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
	go worker("06")
	go worker("07")
	go worker("08")
	go worker("09")

	//TODO use waitgroup
	time.Sleep(time.Hour * 1)
}

func worker(id string) {
	i := 1
	for {
		rand.Seed(time.Now().UnixNano())

		randomInt := rand.Intn(1000)

		// log.Println(i, "sleeping for ms:", randomInt)
		time.Sleep(time.Millisecond * time.Duration(randomInt))

		input := fmt.Sprintf("worker-no=%s,message-no=%d", id, i)

		url := "http://127.0.0.1:8000"

		_, err := http.Post(url, "string", bytes.NewBuffer([]byte(input)))
		if err != nil {
			panic(err)
		}
		if i%100 == 0 {
			time.Sleep(time.Second * time.Duration(rand.Intn(10)))
			log.Println("working going to sleep", id)
		}

		// log.Println(fmt.Sprintf("%d status, %s", res.StatusCode, input))

		i++
	}
}
