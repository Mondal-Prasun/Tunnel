package main

import (
	"context"
	"log"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
// func (a *App) Greet(name string) string {
// 	return fmt.Sprintf("Hello %s, It's show time!", name)
// }
//

func (a *App) FetchTrackerFile(url string) error {

	err := getTrackerFile(url)

	if err != nil {
		return err
	}
	return nil
}

func (a *App) ListenToPeers(port string) error {
	log.Println("Listen: ", port)
	go listenToTheOtherPeers(port)

	return nil
}

func (a *App) GetRequiredContent() ([]TunnelTracerContent, error) {

	log.Println("GetRequiredContent: ", "Called")

	content, err := getTrackerContent()

	if err != nil {
		log.Println("GetrequiredContent: ", err.Error())
		return nil, err
	}

	log.Println("GetRequiredConent: ", content)

	return content, nil

}
