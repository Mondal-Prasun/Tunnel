package main

import (
	"context"
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

	go listenToTheOtherPeers(port)

	return nil
}
