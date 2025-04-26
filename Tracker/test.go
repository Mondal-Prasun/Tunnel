package main

import (
	"log"
	"os"
)

func testFunction() {

	parentFile := "storage/blossom.mp4"
	madeFile := "made/blossom.mp4"

	tfm, err := segmentFile(parentFile)

	if err != nil {
		log.Panicln("Error:", err.Error())
	}

	JointBLFiles(tfm)

	pb, err := os.ReadFile(parentFile)

	if err != nil {
		log.Panicln("Error:", err.Error())
	}

	mb, err := os.ReadFile(madeFile)

	if err != nil {
		log.Panicln("Error:", err.Error())
	}

	for i, b := range pb {

		if mb[i] != b {
			log.Println(mb[i], "!=", b)
		}
	}

}
