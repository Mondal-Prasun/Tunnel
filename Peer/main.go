package main

import (
	"log"
)

func main() {
	log.Println("this is peer")

	makeRequiredDirs()

	// go getTrackerFile()

	// go announceFile("storage/blossom.mp4")

	go listenToTheOtherPeers()

	go requestSegments("a4406c90d99c8c3be19fa8cf6c5d34c7bd462d61edb98a00f01324cafa21c052")

	for {
	}

}
