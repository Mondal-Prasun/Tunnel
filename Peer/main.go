package main

import (
	"log"
)

func main() {
	log.Println("this is peer1")

	// testFunction()
	makeRequiredDirs()

	go getTrackerFile()

	go announceFile("storage/blossom.mp4")

	go listenToTheOtherPeers()

	for {
	}

}
