package badger

import (
	"log"
	"sync"
	"time"
)

type Badge struct {
	sync.Mutex
	data        []string
	batchSize   int
	currentSize int
	flush       chan bool
	flushTime   time.Duration
}

func NewBadge(batchSize int, flushTime time.Duration) *Badge {
	data := make([]string, 0, batchSize)
	flush := make(chan bool)

	badge := &Badge{
		data:        data,
		batchSize:   batchSize,
		currentSize: 0,
		flush:       flush,
		flushTime:   flushTime,
	}

	go badge.waitTillPublish()
	go timeout(badge.flush, badge.flushTime)

	return badge
}

func (b *Badge) Push(element string) {
	b.Lock()

	b.data = append(b.data, element)
	b.currentSize++
	if b.currentSize == b.batchSize {
		b.flush <- false
	}

	b.Unlock()
}

func (b *Badge) publish() {
	b.Lock()

	log.Println("batch-size", b.currentSize, b.data)

	b.currentSize = 0
	b.data = nil
	b.flush = make(chan bool)

	b.Unlock()

	go b.waitTillPublish()
	go timeout(b.flush, b.flushTime)

}

func (b *Badge) waitTillPublish() {
	timeout := <-b.flush

	log.Println("---------------------------------------------------------")
	defer log.Println("---------------------------------------------------------")

	if timeout {
		log.Println("timed-out")
	} else {
		log.Println("batch full")
	}

	b.publish()
}

func timeout(flush chan bool, flushTime time.Duration) {
	time.Sleep(flushTime)
	flush <- true
}
