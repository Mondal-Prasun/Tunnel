package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Start")

	filePath := "storage/cat.mp4"

	if err := segmentFile(filePath); err != nil {
		log.Println("Something went wrong: ", err.Error())
		return
	}

}
