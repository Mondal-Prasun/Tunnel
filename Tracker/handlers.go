package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	// "github.com/google/uuid"
)

type SegDet struct {
	FileSegmentHash string `json:"fileSegmentHash"`
	SegmentNumber   int8   `json:"segNum"`
	ContentSize     int64  `json:"contentSize"`
	SegFileSize     int64  `json:"segFileSize"`
}

type SegmentFileAddress struct {
	FileSegmentHash string   `json:"fileSegmentHash"`
	SegContentSize  int64    `json:"segContentSize"`
	SegFileSize     int64    `json:"segFileSize"`
	SegmentNumber   int8     `json:"SegmentNumber"`
	SegAddress      []string `json:"segAddress"`
}

type TunnelTracerContent struct {
	FileHash         string               `json:"fileHash"`
	FileName         string               `json:"fileName"`
	FileImage        string               `json:"fileImage"`
	FileDescription  string               `json:"fileDescription"`
	FileSize         int64                `json:"fileSize"`
	AllSegmentCount  int8                 `json:"allSegmentCount"`
	FileExt          string               `json:"fileExt"`
	AllFileSegements []SegmentFileAddress `json:"fileSegments"`
}

type Tunnel struct {
}

var mu sync.Mutex

func (t *Tunnel) HealthCheck(w http.ResponseWriter, r *http.Request) {

	log.Println("Checking health")

	tr := TunnelResponse{
		W: w,
	}

	tr.ResponseWithJson(STATUS_RESPONSE_OK, struct {
		ServerOk string `json:"serverOkk"`
	}{
		ServerOk: "Its running",
	})

}

func (t *Tunnel) NewContentAnnounce(w http.ResponseWriter, r *http.Request) {

	body := struct {
		UAddress        string   `json:"uAddress"`
		FileName        string   `json:"fileName"`
		FileSize        int64    `json:"fileSize"`
		FileImage       string   `json:"fileImage"`
		FileDescription string   `json:"fileDescription"`
		FileHash        string   `json:"fileHash"`
		FileExt         string   `json:"fileExt"`
		AllSegmentCount int8     `json:"allSegmentCount"`
		FileSegments    []SegDet `json:"fileSegments"`
	}{}

	log.Println("New Announce requested")

	tr := TunnelResponse{
		W: w,
	}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		log.Println("Announce: ", err.Error())
		tr.ResponseWithError(STATUS_RESPONSE_ERROR, err.Error())
		return
	}

	mu.Lock()

	rFile, err := os.Open(fmt.Sprintf("%s/%s", TRACKFILE_FOLDER, TRACKER_FILE_NAME))

	if err != nil {
		log.Println("Announce: ", err.Error())

		tr.ResponseWithError(503, fmt.Sprintf("Something went wrong while1:%s", err.Error()))
		return
	}

	defer rFile.Close()

	var trackerDetails []TunnelTracerContent

	rDecoder := json.NewDecoder(rFile)

	err = rDecoder.Decode(&trackerDetails)

	if err != nil {
		log.Println("Announce: ", err.Error())

		tr.ResponseWithError(503, fmt.Sprintf("Something went wrong while2:%s", err.Error()))
		return
	}

	alreadyExcist := false

	for _, td := range trackerDetails {
		if td.FileHash == body.FileHash {
			alreadyExcist = true
			break
		}
	}

	if alreadyExcist {
		tr.ResponseWithError(503, fmt.Sprintf("%s :already excist", body.FileHash))
		return
	}

	tFile, err := os.Create(fmt.Sprintf("%s/%s", TRACKFILE_FOLDER, TRACKER_FILE_NAME))
	if err != nil {
		tr.ResponseWithError(503, fmt.Sprintf("Something went wrong while:%s", err.Error()))
		return
	}
	defer tFile.Close()

	segmentFileAddress := make([]SegmentFileAddress, len(body.FileSegments))
	for i, seg := range body.FileSegments {
		segmentFileAddress[i] = SegmentFileAddress{
			FileSegmentHash: seg.FileSegmentHash,
			SegContentSize:  seg.ContentSize,
			SegmentNumber:   seg.SegmentNumber,
			SegFileSize:     seg.SegFileSize,
			SegAddress:      []string{body.UAddress},
		}
	}

	trackerCon := &TunnelTracerContent{
		FileHash:         body.FileHash,
		FileName:         body.FileName,
		FileSize:         body.FileSize,
		FileImage:        body.FileImage,
		FileDescription:  body.FileDescription,
		AllSegmentCount:  body.AllSegmentCount,
		FileExt:          body.FileExt,
		AllFileSegements: segmentFileAddress,
	}

	log.Printf("%+v", trackerCon)

	trackerDetails = append(trackerDetails, *trackerCon)

	marshaled, err := json.MarshalIndent(trackerDetails, "", " ")

	if err != nil {
		log.Panic("NewContent:", err.Error())
		return
	}

	_, err = tFile.Write(marshaled)

	if err != nil {
		tr.ResponseWithError(503, fmt.Sprintf("Something went wrong while:%s", err.Error()))
		return
	}

	mu.Unlock()

	tr.ResponseWithJson(STATUS_RESPONSE_OK, struct {
		FileName string `json:"fileName"`
	}{
		FileName: body.FileHash,
	})

}

func (t *Tunnel) GetTrackerFile(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Disposition", "attachmet; filename=\"tracker.json\"")
	w.Header().Set("Content-type", "application/octet-stream")

	http.ServeFile(w, r, fmt.Sprintf("%s/%s", TRACKFILE_FOLDER, TRACKER_FILE_NAME))
}
