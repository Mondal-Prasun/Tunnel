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

	// log.Println("GetRequiredConent: ", content)

	return content, err

}

func (a *App) MakeRequiredFile() error {

	errChan := make(chan error)
	defer close(errChan)
	go func() {
		err := makeRequiredDirs()
		errChan <- err
	}()

	if err := <-errChan; err != nil {
		return err
	}

	return nil
}

func (a *App) AnnounceCurrentFile(filePath string, fileImage string, fileDescription string, trackerUrl string, port string) error {

	err := announceFile(filePath, trackerUrl, fileImage, fileDescription, port)

	if err != nil {
		return err
	}

	return nil
}

func (a *App) RequestRequiredSegments(parentFileHash string) error {

	err := requestSegments(parentFileHash)

	if err != nil {
		return err
	}

	return nil
}

func (a *App) CheckIfAllSegmentAreAvaliable(allSeg []SegmentFileAddress) ([]SegmentFileAddress, error) {

	needSeg, err := checkIfPeerHasSeg(allSeg)

	if err != nil {
		log.Println("CheckIF", err.Error())
		return nil, err
	}

	log.Println("CheckIF", needSeg)

	return needSeg, nil
}

func (a *App) MakeOriginaleFile(parentFileHash string) error {
	err := makeOriginalFile(parentFileHash)

	if err != nil {
		return err
	}

	return nil
}
