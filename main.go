package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Start")

	filePath := "storage/cat.mp4"

	segFileData, err := segmentFile(filePath)

	if err != nil {
		log.Println("Cannot segment Files:", err.Error())
		return
	}

	jointBLFiles(*segFileData)

}
