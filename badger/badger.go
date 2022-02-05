package badger

import (
	"log"
	"sync"
	"time"
)

func Initiate(batchSize int, timeout time.Duration) chan string {
	pushToBatch := make(chan string, 1000)
	publish := make(chan PublisherMessage, 100)

	badgersSize := 3
	publishersSize := 1

	for i := 0; i < badgersSize; i++ {
		NewBadge(pushToBatch, publish, batchSize, timeout)
		time.Sleep(time.Millisecond * 500)
	}
	log.Println("started", badgersSize, "badges")

	for i := 0; i < publishersSize; i++ {
		go publishWrite(publish)
	}
	log.Println("started", badgersSize, "publishers", publishersSize)

	return pushToBatch
}

type Badge struct {
	sync.Mutex
	data            []string
	batchSize       int
	currentSize     int
	skipNextTimeout bool
	timeout         time.Duration
}

type PublisherMessage struct {
	timeout bool
	data    []string
}

func NewBadge(pull chan string, publish chan PublisherMessage, batchSize int, timeout time.Duration) {
	data := make([]string, 0, batchSize)

	badge := Badge{
		data:            data,
		batchSize:       batchSize,
		currentSize:     0,
		skipNextTimeout: false,
		timeout:         timeout,
	}

	go badge.listen(pull, publish)
	go badge.publishAfterTimeout(publish)
}

func (b *Badge) listen(pull chan string, publish chan PublisherMessage) {
	for {
		newElement := <-pull
		b.Lock()
		b.data = append(b.data, newElement)
		b.currentSize++
		if b.currentSize == b.batchSize {
			publish <- PublisherMessage{
				timeout: false,
				data:    b.data,
			}
			b.data = nil
			b.currentSize = 0
			b.skipNextTimeout = true
		}
		b.Unlock()
	}
}

func (b *Badge) publishAfterTimeout(publish chan PublisherMessage) {
	for {
		time.Sleep(b.timeout)
		b.Lock()
		if b.skipNextTimeout {
			b.skipNextTimeout = false
			b.Unlock()
			log.Println("skipping as last badge was full")
			continue
		}
		if b.currentSize == 0 {
			b.Unlock()
			log.Println("badge empty: skipping")
			continue
		}

		publish <- PublisherMessage{
			timeout: true,
			data:    b.data,
		}
		b.data = nil
		b.currentSize = 0
		b.skipNextTimeout = false
		b.Unlock()
	}

}

func publishWrite(pull chan PublisherMessage) {
	for {
		msg := <-pull
		log.Println("----------------------------------------------------------------------------------------")
		if msg.timeout {
			log.Println("timeout")
		} else {
			log.Println("badge full")

		}
		log.Println("badge size", len(msg.data))
		log.Println(msg.data)
		log.Println("----------------------------------------------------------------------------------------")
	}
}
