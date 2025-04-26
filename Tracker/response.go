package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type TunnerResponse struct {
	W http.ResponseWriter
}

func (tr *TunnerResponse) ResponseWithJson(
	statusCode int,
	body interface{}) {

	data, err := json.Marshal(body)

	if err != nil {
		log.Panic("Cannot marshal json:", err.Error())
		return
	}

	tr.W.Header().Add("Content-type", "application/json")
	tr.W.WriteHeader(statusCode)
	tr.W.Write(data)
}

func (tr *TunnerResponse) ResponseWithError(statusCode int, errorMsg string) {

	err := struct {
		ErrorMsg string `json:"errorMsg"`
	}{
		ErrorMsg: errorMsg,
	}

	tr.ResponseWithJson(statusCode, err)
}
