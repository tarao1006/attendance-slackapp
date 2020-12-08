package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	a := "15:00"
	b := "15:00"

	at, err := time.Parse("15:04", a)
	if err != nil {
		log.Fatal(err)
	}
	bt, err := time.Parse("15:04", b)
	if err != nil {
		log.Fatal(err)
	}

	if !bt.After(at) {
		fmt.Println("hoge")
	}
}
