package badger

import (
	"log"
	"time"
)

type Badge struct {
	data        []string
	batchSize   int
	currentSize int
	flushTime   time.Duration
}

func NewBadge(batchSize int, flushTime time.Duration) *Badge {
	data := make([]string, 0, batchSize)

	badge := &Badge{
		data:        data,
		batchSize:   batchSize,
		currentSize: 0,
		flushTime:   flushTime,
	}

	go badge.timeout()
	return badge
}

func (b *Badge) Push(element string) {
	b.data = append(b.data, element)
	b.currentSize++
	if b.currentSize == b.batchSize {
		b.Flush()
	}

}

func (b *Badge) Flush() {
	log.Println(b.data)

	b.currentSize = 0
	b.data = nil
	b.timeout()
}

func (b *Badge) timeout() {
	time.Sleep(b.flushTime)
	b.Flush()
}
