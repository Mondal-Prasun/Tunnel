package main

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Tunnel struct {
	SqlDb *SqlDb
}

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

func (t *Tunnel) NewContent(w http.ResponseWriter, r *http.Request) {

	body := struct {
		Uid       string `json:"uid"`
		FileName  string `json:"fileName"`
		FileSize  string `json:"fileSize"`
		FileImage string `json:"fileImage"`
	}{}

	tr := TunnerResponse{
		W: w,
	}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		tr.ResponseWithError(STATUS_RESPONSE_ERROR, err.Error())
		return
	}

	userId, err := uuid.Parse(body.Uid)

	if err != nil {
		tr.ResponseWithError(STATUS_RESPONSE_ERROR, err.Error())
		return
	}

	tunnelContent := &TunnelContent{
		Cid:       uuid.New(),
		Uid:       userId,
		FileName:  body.FileName,
		FileSize:  body.FileSize,
		FileImage: body.FileImage,
		FileHash:  "",
	}

	err = t.SqlDb.InsertNewContentInformation(tunnelContent)

	if err != nil {
		tr.ResponseWithError(STATUS_RESPONSE_ERROR, err.Error())
		return
	}

	tr.ResponseWithJson(STATUS_RESPONSE_OK, struct {
		Cid       string `json:"cid"`
		FileName  string `json:"fileName"`
		FileSize  string `json:"fileSize"`
		FileImage string `json:"fileImage"`
	}{
		Cid:       tunnelContent.Cid.String(),
		FileName:  tunnelContent.FileName,
		FileImage: tunnelContent.FileImage,
		FileSize:  tunnelContent.FileSize,
	})

}
