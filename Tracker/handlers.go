package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TunnelContent struct {
	Cid              uuid.UUID `json:"id"`
	Uid              uuid.UUID `json:"uid"`
	UAddress         string    `json:"uAddress"`
	FileName         string    `json:"fileName"`
	FileSize         string    `json:"fileSize"`
	FileImage        string    `json:"fileImage"`
	FileHash         string    `json:"fileHash"`
	FileSegmentsHash []string  `json:"fileSegmentsHash"`
}

type SegmentFileAddress struct {
	FileSegmentHash string   `json:"fileSegmentHash"`
	FileAddress     []string `json:"fileAddress"`
}

type TunnelTracerContent struct {
	FileHash         string               `json:"fileHash"`
	AllFileSegements []SegmentFileAddress `json:"fileSegments"`
}

type Tunnel struct {
	SqlDb *SqlDb
}

var mu sync.Mutex

func (t *Tunnel) HealthCheck(w http.ResponseWriter, r *http.Request) {

	log.Println("Checking health")

	tr := TunnerResponse{
		W: w,
	}

	tr.ResponseWithJson(STATUS_RESPONSE_OK, struct {
		ServerOk string `json:"serverOkk"`
	}{
		ServerOk: "Its running",
	})

}

func (t *Tunnel) SignupUser(w http.ResponseWriter, r *http.Request) {
	body := struct {
		UserName  string `json:"userName"`
		Password  string `json:"password"`
		UserImage string `json:"userImage"`
	}{}

	tr := TunnerResponse{
		W: w,
	}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		tr.ResponseWithError(STATUS_RESPONSE_ERROR, err.Error())
		return
	}

	encryptPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	encryptPasswordString := hex.EncodeToString(encryptPassword)
	if err != nil {
		tr.ResponseWithError(STATUS_RESPONSE_SERVER_ERROR, err.Error())
		return
	}

	tu := TunnelUser{
		Id:        uuid.New(),
		UserName:  body.UserName,
		UserImage: body.UserImage,
		Password:  encryptPasswordString,
	}

	err = t.SqlDb.InsertUser(tu)

	if err != nil {
		tr.ResponseWithError(STATUS_RESPONSE_SERVER_ERROR, err.Error())
		return
	}

	tr.ResponseWithJson(STATUS_RESPONSE_OK, struct {
		UserId string `json:"userId"`
	}{
		UserId: tu.Id.String(),
	})
}

func (t *Tunnel) LoginUser(w http.ResponseWriter, r *http.Request) {

	body := struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}{}
	tr := TunnerResponse{
		W: w,
	}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		tr.ResponseWithError(STATUS_RESPONSE_ERROR, err.Error())
		return
	}

	tUser, err := t.SqlDb.QueryUser(body.UserName)

	if err != nil {
		tr.ResponseWithError(STATUS_RESPONSE_ERROR, err.Error())
		return
	}

	passwordByte, err := hex.DecodeString(tUser.Password)

	if err != nil {
		log.Panic("Login User:", err.Error())
		return
	}

	if err := bcrypt.CompareHashAndPassword(passwordByte, []byte(body.Password)); err != nil {
		tr.ResponseWithError(STATUS_RESPONSE_ERROR, err.Error())
		return
	}

	tr.ResponseWithJson(STATUS_RESPONSE_OK, struct {
		Id        string `json:"uid"`
		UserName  string `json:"userName"`
		UserImage string `json:"userImage"`
	}{
		Id:        tUser.Id.String(),
		UserName:  tUser.UserName,
		UserImage: tUser.UserImage,
	})

}

func (t *Tunnel) NewContentAnnounce(w http.ResponseWriter, r *http.Request) {

	body := struct {
		Uid              string   `json:"uid"`
		UAddress         string   `json:"uAddress"`
		FileName         string   `json:"fileName"`
		FileSize         string   `json:"fileSize"`
		FileImage        string   `json:"fileImage"`
		FileHash         string   `json:"fileHash"`
		FileSegmentsHash []string `json:"fileSegmentsHash"`
	}{}

	tr := TunnerResponse{
		W: w,
	}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		tr.ResponseWithError(STATUS_RESPONSE_ERROR, err.Error())
		return
	}

	// userId, err := uuid.Parse(body.Uid)

	// if err != nil {
	// 	tr.ResponseWithError(STATUS_RESPONSE_ERROR, err.Error())
	// 	return
	// }

	tunnelContent := &TunnelContent{
		Cid:              uuid.New(),
		Uid:              uuid.New(),
		FileName:         body.FileName,
		FileSize:         body.FileSize,
		FileImage:        body.FileImage,
		FileHash:         body.FileHash,
		UAddress:         body.UAddress,
		FileSegmentsHash: body.FileSegmentsHash,
	}

	mu.Lock()

	rFile, err := os.Open(TRACKER_FILE_NAME)

	if err != nil {
		tr.ResponseWithError(503, fmt.Sprintf("Something went wrong while1:%s", err.Error()))
		return
	}

	defer rFile.Close()

	var trackerDetails []TunnelTracerContent

	rDecoder := json.NewDecoder(rFile)

	err = rDecoder.Decode(&trackerDetails)

	if err != nil {
		tr.ResponseWithError(503, fmt.Sprintf("Something went wrong while2:%s", err.Error()))
		return
	}

	tFile, err := os.Create(TRACKER_FILE_NAME)
	if err != nil {
		tr.ResponseWithError(503, fmt.Sprintf("Something went wrong while:%s", err.Error()))
		return
	}
	defer tFile.Close()

	segmentFileAddress := make([]SegmentFileAddress, len(body.FileSegmentsHash))
	for i, seg := range body.FileSegmentsHash {
		segmentFileAddress[i] = SegmentFileAddress{
			FileSegmentHash: seg,
			FileAddress:     []string{body.UAddress},
		}
	}

	trackerCon := &TunnelTracerContent{
		FileHash:         tunnelContent.FileHash,
		AllFileSegements: segmentFileAddress,
	}

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
		Cid               string   `json:"cid"`
		FileName          string   `json:"fileName"`
		FileSize          string   `json:"fileSize"`
		FileImage         string   `json:"fileImage"`
		FileHash          string   `json:"fileHash"`
		FileSegementsHash []string `json:"fileSegmentsHash"`
	}{
		FileSegementsHash: body.FileSegmentsHash,
	})

}

func (t *Tunnel) GetAllContent(w http.ResponseWriter, r *http.Request) {

	contents, err := t.SqlDb.GetAllContent()

	tr := TunnerResponse{
		W: w,
	}

	if err != nil {
		tr.ResponseWithError(STATUS_RESPONSE_ERROR, err.Error())
		return
	}

	tr.ResponseWithJson(STATUS_RESPONSE_OK, struct {
		AllContents []TunnelContent `json:"allContents"`
	}{
		AllContents: contents,
	})
}
