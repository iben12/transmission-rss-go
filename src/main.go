package main

import (
	"fmt"
	// "github.com/hekmon/transmissionrpc"
	// "os"
)

func main() {
	fmt.Println("Starting")

	torrents := getFinished()

	for _, torrentId := range torrents {
		fmt.Printf("ID: %d\n", torrentId)
	}

}
