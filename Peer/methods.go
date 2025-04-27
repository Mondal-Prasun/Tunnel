package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"slices"
	"strings"
)

func makeRequiredDirs() {

	allDir := []string{
		JOINT_STORE_DIRECTORY,
		ORIGINAL_FILE_STORAGE,
		SEGEMENT_STORE_DIRECTORY,
	}

	for _, dir := range allDir {

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.Mkdir(dir, os.ModeDir)
			if err != nil {
				log.Panic("Create Directory: ", err.Error())
			}

		}

	}

}

func getTrackerFile() {

	res, err := http.Get(fmt.Sprintf("%s/getTracker", TRACKER_URL))

	if err != nil {
		log.Println("GetTrackerFile:", err.Error())
		return
	}

	defer res.Body.Close()

	tFile, err := os.Create(TRACKER_PATH)

	if err != nil {
		log.Println("GetTrackerFile:", err.Error())
		return
	}

	defer tFile.Close()

	n, err := io.Copy(tFile, res.Body)

	if err != nil {
		log.Println("GetTrackerFile:", err.Error())
		return
	}

	log.Println("Written file byte: ", n)

}

func announceFile(filePath string) {

	tMeta, err := SegmentFile(filePath)

	if err != nil {
		log.Println("Announce: ", err.Error())
		return
	}

	var segFileHash []string

	for _, seg := range tMeta.AllSegments {

		segFileHash = append(segFileHash, seg.Filehash)

	}

	data := struct {
		Uid              string   `json:"uid"`
		UAddress         string   `json:"uAddress"`
		FileName         string   `json:"fileName"`
		FileSize         string   `json:"fileSize"`
		FileImage        string   `json:"fileImage"`
		FileHash         string   `json:"fileHash"`
		FileSegmentsHash []string `json:"fileSegmentsHash"`
	}{
		Uid:              "",
		UAddress:         PEER_ADDRESS,
		FileName:         tMeta.ParentFileName,
		FileSize:         "", //fix it later
		FileImage:        "", //fix it later
		FileHash:         tMeta.ParentFilehash,
		FileSegmentsHash: segFileHash,
	}

	marshaledata, err := json.Marshal(data)

	if err != nil {
		log.Println("Announce: ", err.Error())
		for _, seg := range tMeta.AllSegments {
			err := os.Remove(seg.FileDestination)

			if err != nil {
				log.Println("Announce: ", err.Error())
				break
			}

		}
		return
	}

	res, err := http.Post(fmt.Sprintf("%s/announce", TRACKER_URL), "application-json", bytes.NewBuffer(marshaledata))

	if err != nil {
		log.Println("Announce: ", err.Error())
		return
	}

	defer res.Body.Close()

	log.Println(res.Body)

	go getTrackerFile()

}

func listenToTheOtherPeers() {

	listner, err := net.Listen("tcp", PEER_ADDRESS)

	if err != nil {
		log.Panic("ListenToPeers:", err.Error())
		return
	}

	log.Println("Peer listening on: ", PEER_ADDRESS)

	defer listner.Close()

	for {
		conn, err := listner.Accept()

		if err != nil {
			log.Println("ListenToPeers: ", err.Error())
			continue
		}

		go HandlePeerContections(conn)

	}
}

func HandlePeerContections(conn net.Conn) {

	defer conn.Close()

	connReader := bufio.NewReader(conn)

	for {

		reqString, err := connReader.ReadString('\n')

		log.Println("HandlePeer: ", reqString)

		if err != nil {
			return
		}

		if strings.HasPrefix(reqString, "REQUEST:") {
			segIdStr := strings.TrimPrefix(reqString, "REQUEST:")

			msg := fmt.Sprintf("CHUNK:%s\n", segIdStr)
			conn.Write([]byte(msg))

			{

				segBytes, err := os.ReadFile(fmt.Sprintf("%s/%s.bl", SEGEMENT_STORE_DIRECTORY,
					strings.TrimSpace(segIdStr)))

				if err != nil {
					log.Println("Handlepeer: ", err.Error())
					return
				}

				log.Println(len(segBytes))

				conn.Write(segBytes)

			}

		}

	}

}

type SegmentFileAddress struct {
	FileSegmentHash string   `json:"fileSegmentHash"`
	FileAddress     []string `json:"fileAddress"`
}

type TunnelTracerContent struct {
	FileHash         string               `json:"fileHash"`
	AllFileSegements []SegmentFileAddress `json:"fileSegments"`
}

func requestSegments(parentFileHash string) error {

	rFile, err := os.Open(TRACKER_PATH)

	if err != nil {
		log.Println("RequestSegments: ", err.Error())
		return err
	}

	defer rFile.Close()

	var trackerDetails []TunnelTracerContent

	rDecoder := json.NewDecoder(rFile)

	err = rDecoder.Decode(&trackerDetails)

	if err != nil {
		log.Println("RequestSegments: ", err.Error())
		return err
	}

	var allSegmentsHash []string

	var segAdd []SegmentFileAddress

	for _, trackerDet := range trackerDetails {

		if trackerDet.FileHash == parentFileHash {
			for _, segHash := range trackerDet.AllFileSegements {
				allSegmentsHash = append(allSegmentsHash, segHash.FileSegmentHash)
				segAdd = append(segAdd, segHash)
			}
		}
	}

	if len(allSegmentsHash) == 0 {
		return fmt.Errorf("there is no file in the tracker")
	}

	neededSeg := checkIfPeerHasSeg(allSegmentsHash)

	if len(neededSeg) == 0 {
		log.Println("All segments already excists")
		return nil
	}
	// for _, nSeg := range neededSeg {
	// 	log.Println(nSeg)
	// }

	var neededSegAdd []SegmentFileAddress

	for _, seg := range neededSeg {
		for _, segA := range segAdd {
			if seg == segA.FileSegmentHash {
				neededSegAdd = append(neededSegAdd, segA)
			}
		}
	}

	// for _, a := range neededSegAdd {
	// 	log.Println(a.FileAddress)
	// }

	// for {

	for _, nS := range neededSegAdd {
		if nS.FileAddress[0] != PEER_ADDRESS {
			go connecTopeer(nS.FileAddress[0], nS.FileSegmentHash)
		}
	}
	// time.Sleep(10 * time.Second)
	// }

	return nil
}

func checkIfPeerHasSeg(segMentsName []string) (neededSeg []string) {

	dirEntry, err := os.ReadDir(SEGEMENT_STORE_DIRECTORY)

	if err != nil {

		log.Panic("CheckPeerSeg:", err.Error())

	}

	var needSeg []string
	// dirE.Name()[:(len(dirE.Name())-3)]

	var fileName []string

	for _, dirE := range dirEntry {
		fileName = append(fileName, dirE.Name()[:(len(dirE.Name())-3)])
	}

	for _, segName := range segMentsName {
		if !slices.Contains(fileName, segName) {
			needSeg = append(needSeg, segName)
		}
	}
	return needSeg

}

func connecTopeer(add string, segId string) {

	conn, err := net.Dial("tcp", add)

	if err != nil {
		log.Println("ConnecToPeer:", add, " is not online")
		return
	}

	defer conn.Close()

	req := fmt.Sprintf("REQUEST:%s\n", segId)

	conn.Write([]byte(req))

	reader := bufio.NewReader(conn)

	header, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("ConnectoPeer:", err.Error())
		return
	}

	log.Println("peer: ", header)

	idStr := strings.TrimPrefix(header, "CHUNK:")
	log.Println("peer:", idStr)

	reader2 := bufio.NewReader(conn)

	header2, err := reader2.ReadString('\n')

	if err != nil {
		fmt.Println("ConnectoPeer:", err.Error())
		return
	}

	log.Println("peer: ", header2)

}
