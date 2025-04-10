package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Start")

	sDb, err := initDataBase()

	if err != nil {
		log.Panic("Cannot initilize database:", err.Error())
		return
	}

	defer sDb.closeDataBase()

	tunnel := Tunnel{
		SqlDb: sDb,
	}

	tunnelMux := http.NewServeMux()

	tunnelMux.HandleFunc("/health", tunnel.HealthCheck)
	tunnelMux.HandleFunc("/signUp", tunnel.SignupUser)
	tunnelMux.HandleFunc("/login", tunnel.LoginUser)
	tunnelMux.HandleFunc("/relay", tunnel.RelaySocketHandler)
	tunnelMux.HandleFunc("/clients", tunnel.OnlineStatusHandler)
	tunnelMux.HandleFunc("/newContent", tunnel.NewContent)

	log.Println("Server started at:", SERVER_PORT)

	if err := http.ListenAndServe(SERVER_PORT, tunnelMux); err != nil {
		log.Panic("Cannot able to start server:", err.Error())
	}

}
