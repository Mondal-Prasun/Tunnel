package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type TunnelClient struct {
	Id   string
	Conn *websocket.Conn
}

var (
	tunnelClients = make(map[string]*TunnelClient)
	mu            sync.Mutex
	relayUgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (t *Tunnel) RelaySocketHandler(w http.ResponseWriter, r *http.Request) {

	uid := r.URL.Query().Get("Uid")

	relayConn, err := relayUgrader.Upgrade(w, r, nil)
	tr := TunnerResponse{
		W: w,
	}

	if err != nil {
		tr.ResponseWithError(STATUS_RESPONSE_ERROR, err.Error())
		return
	}

	defer relayConn.Close()

	mu.Lock()
	tunnelClients[uid] = &TunnelClient{
		Id:   uid,
		Conn: relayConn,
	}
	mu.Unlock()

	for {
		// Read message
		_, msg, err := relayConn.ReadMessage()
		if err != nil {
			fmt.Println("Client disconnected")
			break
		}

		// Broadcast message to all other clients
		mu.Lock()
		for c := range tunnelClients {
			if c != uid { // Don't echo back to sender

				data := struct {
					Id  string `json:"id"`
					Msg string `json:"msg"`
				}{
					Id:  uid,
					Msg: string(msg),
				}

				err := tunnelClients[c].Conn.WriteJSON(data)
				if err != nil {
					fmt.Println("Write error:", err)
					tunnelClients[c].Conn.Close()
					delete(tunnelClients, c)
				}
			}
		}
		mu.Unlock()
	}
}

func (t *Tunnel) OnlineStatusHandler(w http.ResponseWriter, r *http.Request) {
	var onlineClients []string

	for c := range tunnelClients {
		onlineClients = append(onlineClients, c)
	}

	tr := TunnerResponse{
		W: w,
	}

	tr.ResponseWithJson(STATUS_RESPONSE_OK, struct {
		OnlineIds []string `json:"onlineIds"`
	}{
		OnlineIds: onlineClients,
	})

}

func (t *Tunnel) SendFileRequestToSeeder(w http.ResponseWriter, r *http.Request) {

	uid := r.URL.Query().Get("Uid")

	cid := r.URL.Query().Get("Cid")

	tr := TunnerResponse{
		W: w,
	}

	if uid == "" || cid == "" {
		tr.ResponseWithError(STATUS_RESPONSE_ERROR, "Please give uid and cid")
		return
	}

	seederChannel := make(map[string]*TunnelClient)

	relayConn, err := relayUgrader.Upgrade(w, r, nil)

	if err != nil {
		tr.ResponseWithError(STATUS_RESPONSE_ERROR, err.Error())
		return
	}

	mu.Lock()
	seederChannel[uid] = &TunnelClient{
		Id:   uid,
		Conn: relayConn,
	}
	mu.Unlock()

	defer relayConn.Close()

	tunnelContentDetails, err := t.SqlDb.GetContentDetails(cid)

	if err != nil {
		tr.ResponseWithError(STATUS_RESPONSE_ERROR, err.Error())
		return
	}

	for c := range tunnelClients {

		if c == tunnelContentDetails.Uid.String() {

			mu.Lock()
			seederChannel[tunnelContentDetails.Uid.String()] = &TunnelClient{
				Id:   tunnelContentDetails.Uid.String(),
				Conn: relayConn,
			}
			mu.Unlock()

			//TODO : implement the whole thing

		} else {
			tr.ResponseWithError(STATUS_RESPONSE_ERROR, fmt.Sprintf("%s is not online", uid))
			return
		}
	}

}
