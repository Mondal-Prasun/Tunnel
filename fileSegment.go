package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	SEGMENT_SIZE_EVEN        int8   = 4
	SEGMENT_SIZE_ODD         int8   = 3
	SEGEMENT_STORE_DIRECTORY string = "segments"
	SEGEMENT_EXT             string = ".bl"
)

type Segment struct {
	fileSize        int64
	createdAt       string
	filehash        string
	fileDestination string
	comparitionHash string
	segmentNumber   int8
}

type SegmentFileMetadata struct {
	parentFileName      string
	parentFileExtention string
	parentFilehash      string
	parentFileSize      int64
	segmentCount        int8
}

func segmentFile(filePath string) error {

	fileInfo, err := os.Stat(filePath)

	if err != nil {
		return err
	}

	fileBytes, err := os.ReadFile(filePath)

	if err != nil {
		return err
	}

	var segmentCount int8

	if fileInfo.Size()%2 == 0 {
		segmentCount = SEGMENT_SIZE_EVEN
	} else {
		segmentCount = SEGMENT_SIZE_ODD
	}

	segmentFileSize := fileInfo.Size() / int64(segmentCount)

	log.Println("Segement file size: ", segmentFileSize)

	saveFileBuffer := make([]byte, segmentFileSize)

	f, err := os.Open(filePath)

	if err != nil {
		return err
	}

	defer f.Close()

	for i := range segmentCount {
		_, err = f.ReadAt(saveFileBuffer, segmentFileSize*int64(i))

		if err != nil {
			log.Panicln("Something went wrong while segementing file: ", err.Error())
		}

		log.Println("Segment file num: ", i, "and content: ", saveFileBuffer[0])
		transfromSegmentBl(saveFileBuffer, segmentFileSize, fileBytes, i)
	}

	return nil
}

func transfromSegmentBl(
	segBytes []byte,
	segmentSize int64,
	parentFileByte []byte,
	segmentNum int8,
) (*Segment, error) {

	//check if the segment store folder is avaliable or not
	if _, err := os.Stat(SEGEMENT_STORE_DIRECTORY); err != nil {
		os.Mkdir(SEGEMENT_STORE_DIRECTORY, os.ModeDir)
	}

	//convert segmented byte to its corrosponding hash string
	byteHash := sha256.Sum256(segBytes)
	seghashedName := hex.EncodeToString(byteHash[:])

	log.Println(string(seghashedName))

	//filepath of the segmented file [eg:segments/example.bl]
	segStorePath := SEGEMENT_STORE_DIRECTORY + "/" + seghashedName + SEGEMENT_EXT

	//creates the segmented file
	segf, err := os.Create(segStorePath)
	if err != nil {
		return nil, err
	}

	//writes all the neseccery data to the segmented file
	parentFileHash := sha256.Sum256(parentFileByte)
	parentFilehashString := hex.EncodeToString(parentFileHash[:])
	createdAt := time.Now().Format(time.UnixDate)

	segMetadata := "!METADATA:%" + "parentFilehash=" + parentFilehashString + "childContentHash=" + seghashedName + "segFileSize:" + fmt.Sprint(segmentSize) + "segmentNumber: " + fmt.Sprint(segmentNum) + "createdAt:" + createdAt + "%"

	_, err = segf.WriteString(segMetadata)

	if err != nil {
		os.Remove(segStorePath)
		return nil, err
	}
	_, err = segf.WriteString("!BYTE:")

	if err != nil {
		os.Remove(segStorePath)
		return nil, err
	}

	//oofset for the segmented byte that where it should start storing
	contentOffset := len([]byte(segMetadata)) + len("!BYTE:")

	_, err = segf.WriteAt(segBytes, int64(contentOffset))

	if err != nil {
		os.Remove(segStorePath)
		return nil, err
	}

	//return all info about the corrosponding segment
	segment := Segment{
		fileSize:        segmentSize,
		createdAt:       createdAt,
		filehash:        seghashedName,
		fileDestination: segStorePath,
		segmentNumber:   segmentNum,
		comparitionHash: parentFilehashString,
	}

	return &segment, nil

}
