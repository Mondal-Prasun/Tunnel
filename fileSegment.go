package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const (
	SEGMENT_SIZE             int8   = 6
	SEGEMENT_MIN_FILE_SIZE   int64  = 1024
	SEGMENT_MAGIC_BYTES      string = "BLACKBOX"
	SEGEMENT_STORE_DIRECTORY string = "segments"
	JOINT_STORE_DIRECTORY    string = "made"
	SEGEMENT_EXT             string = ".bl"
)

type Segment struct {
	contentSize     int64
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
	allSegments         []Segment
}

func segmentFile(filePath string) (*SegmentFileMetadata, error) {

	fileInfo, err := os.Stat(filePath)

	log.Println("total file size: ", fileInfo.Size())
	if err != nil {
		return nil, err
	}

	parentFileBytes, err := os.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	segmentFileSize := fileInfo.Size() / int64(SEGMENT_SIZE)

	log.Println("Segement file size: ", segmentFileSize)

	saveFileMul := segmentFileSize / SEGEMENT_MIN_FILE_SIZE

	saveFileBuffer := make([]byte, (SEGEMENT_MIN_FILE_SIZE * (saveFileMul + 1)))

	log.Println("Save buffer size:", len(saveFileBuffer))

	f, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	allSegmentFiles := make([]Segment, SEGMENT_SIZE)

	for i := range SEGMENT_SIZE {
		n, err := f.ReadAt(saveFileBuffer, int64(len(saveFileBuffer))*int64(i))

		if err != nil {
			if err == io.EOF {
				log.Println("Content finished... proceding with 0's")
			} else {
				log.Println("Something went wrong with segmenting file:", err.Error(), "index:", i)
			}
		}

		// log.Println("Segment file num: ", i, "and content: ", saveFileBuffer[0], "readBute size:", n)
		seg, err := transfromSegmentBl(saveFileBuffer, int64(n), parentFileBytes, i)

		if err != nil {
			return nil, err
		}

		allSegmentFiles[i] = *seg
	}

	var parentFileExt string

	for i, ch := range fileInfo.Name() {
		if ch == '.' {
			parentFileExt = fileInfo.Name()[i:]
			break
		}
	}

	log.Println(parentFileExt)

	parentFileHash := sha256.Sum256(parentFileBytes)

	parentFileString := hex.EncodeToString(parentFileHash[:])

	segParent := SegmentFileMetadata{
		parentFileName:      fileInfo.Name(),
		parentFileExtention: parentFileExt,
		parentFilehash:      parentFileString,
		parentFileSize:      fileInfo.Size(),
		segmentCount:        SEGMENT_SIZE,
		allSegments:         allSegmentFiles,
	}
	return &segParent, nil
}

func transfromSegmentBl(
	segBytes []byte,
	segmentSize int64,
	parentFileByte []byte,
	segmentNum int8,
) (*Segment, error) {

	//check if the segment store folder is avaliable or not
	if _, err := os.Stat(SEGEMENT_STORE_DIRECTORY); err != nil {
		err = os.Mkdir(SEGEMENT_STORE_DIRECTORY, os.ModeDir)
		if err != nil {
			log.Println("Cannot make directory:", err.Error())
			return nil, err
		}
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

	_, err = segf.WriteString(SEGMENT_MAGIC_BYTES)

	if err != nil {
		os.Remove(segStorePath)
		return nil, err
	}

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

	//offset for the segmented byte that where it should start storing
	contentOffset := len([]byte(segMetadata)) + len("!BYTE:")

	_, err = segf.WriteAt(segBytes, int64(contentOffset))

	if err != nil {
		os.Remove(segStorePath)
		return nil, err
	}

	//return all info about the corrosponding segment
	segment := Segment{
		contentSize:     segmentSize,
		createdAt:       createdAt,
		filehash:        seghashedName,
		fileDestination: segStorePath,
		segmentNumber:   segmentNum,
		comparitionHash: parentFilehashString,
	}

	return &segment, nil

}

func jointBLFiles(segFileData SegmentFileMetadata) {

	parentFileName := segFileData.parentFileName

	if _, err := os.Stat(JOINT_STORE_DIRECTORY); err != nil {
		err = os.Mkdir(JOINT_STORE_DIRECTORY, os.ModeDir)
		if err != nil {
			log.Println("Cannot make directory:", err.Error())
			return
		}
	}

	saveFilePath := JOINT_STORE_DIRECTORY + "/" + parentFileName

	parentFile, err := os.Create(saveFilePath)

	if err != nil {
		log.Println("Cannot make directory:", err.Error())
		return
	}

	defer parentFile.Close()

	var contentOffset int64 = 0

	for i := range len(segFileData.allSegments) {
		content, err := getContent(segFileData.allSegments[i].fileDestination,
			segFileData.allSegments[i].contentSize)

		if err != nil {
			log.Println("Something went wrong while getting segment content:", err.Error())
			break
		}
		_, err = parentFile.WriteAt(content, contentOffset)

		contentOffset = contentOffset + segFileData.allSegments[i].contentSize

		if err != nil {
			log.Println("Something went wrong while writing the content bytes:", err.Error())
			return
		}
	}

	parentInfo, err := os.Stat(saveFilePath)

	if err != nil {
		log.Println("Cannot find file:", err.Error())
		return
	}

	if segFileData.parentFileSize == parentInfo.Size() {
		log.Println("ParentFile restored")
	} else {
		log.Println("😭")
		log.Println(segFileData.parentFileSize)
		log.Println(parentInfo.Size())
	}

}

func getContent(filePath string, contentSize int64) ([]byte, error) {

	if filePath[(len(filePath)-3):] != ".bl" {
		return nil, errors.New("this is not a valid segFile")
	}

	segFile, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer segFile.Close()

	segInfo, err := os.Stat(filePath)

	if err != nil {
		return nil, err
	}

	segFileCotent := make([]byte, segInfo.Size())

	_, err = segFile.Read(segFileCotent)

	if err != nil {
		log.Println(err.Error())
	}

	index := bytes.Index(segFileCotent, []byte("BYTE:"))

	content := make([]byte, contentSize)

	offSet := (int64(index) + int64(len([]byte("BYTE:"))))

	_, err = segFile.ReadAt(content, offSet)

	if err != nil {
		return nil, err
	}

	return content, nil
}
