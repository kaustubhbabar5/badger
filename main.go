package main

import (
	"fmt"
	"time"

	"github.com/kaustubhbabar5/badger/badger"
)

func main() {
	badge := badger.NewBadge(3, time.Second*5)

	for {
		var input string
		fmt.Scanln(&input)
		badge.Push(input)
	}

}
