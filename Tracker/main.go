package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Start")

	// testFunction()
	requiredFiles()
	tunnel := Tunnel{}

	tunnelMux := http.NewServeMux()

	tunnelMux.HandleFunc("/health", tunnel.HealthCheck)
	tunnelMux.HandleFunc("/announce", tunnel.NewContentAnnounce)
	tunnelMux.HandleFunc("/getTracker", tunnel.GetTrackerFile)

	log.Println("Server started at:", SERVER_PORT)

	if err := http.ListenAndServe(SERVER_PORT, tunnelMux); err != nil {
		log.Panic("Cannot able to start server:", err.Error())
	}

}

func requiredFiles() {
	if _, err := os.Stat("Track"); os.IsNotExist(err) {
		err := os.Mkdir(TRACKFILE_FOLDER, os.ModeDir)
		if err != nil {
			log.Panic("RequiredFile: ", err.Error())
		} else {
			f, err := os.Create(fmt.Sprintf("%s/%s", TRACKFILE_FOLDER, TRACKER_FILE_NAME))
			if err != nil {
				log.Panic("RequiredFile: ", err.Error())
			}
			defer f.Close()
			f.WriteString("[]")
		}
	}
}
