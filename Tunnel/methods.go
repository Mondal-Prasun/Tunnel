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

func getTrackerFile(trackerUrl string) error {

	url := trackerUrl

	res, err := http.Get(fmt.Sprintf("%s/getTracker", url))

	if err != nil {
		log.Println("GetTrackerFile:", err.Error())
		return err
	}

	defer res.Body.Close()

	tFile, err := os.Create(TRACKER_PATH)

	if err != nil {
		log.Println("GetTrackerFile:", err.Error())
		return err
	}

	defer tFile.Close()

	n, err := io.Copy(tFile, res.Body)

	if err != nil {
		log.Println("GetTrackerFile:", err.Error())
		return err
	}

	log.Println("Written file byte: ", n)

	return err

}

func announceFile(filePath string, trackerUrl string) {

	tMeta, err := SegmentFile(filePath)

	if err != nil {
		log.Println("Announce: ", err.Error())
		return
	}

	var segFileHash []string

	for _, seg := range tMeta.AllSegments {

		segFileHash = append(segFileHash, seg.Filehash)

	}

	pFileInfo, err := os.Stat(filePath)

	if err != nil {
		log.Println("Announce: ", err.Error())
		return
	}

	// log.Println("Announce: ", tMeta.AllSegments[0].ContentSize)
	//

	var allSegments []SegDet

	for _, s := range tMeta.AllSegments {
		allSegments = append(allSegments, SegDet{
			FileSegmentHash: s.Filehash,
			SegmentNumber:   s.SegmentNumber,
			ContentSize:     s.ContentSize,
			SegFileSize:     s.FileSize,
		})
	}

	data := struct {
		UAddress        string   `json:"uAddress"`
		FileName        string   `json:"fileName"`
		FileSize        int64    `json:"fileSize"`
		FileImage       string   `json:"fileImage"`
		FileHash        string   `json:"fileHash"`
		FileExt         string   `json:"fileExt"`
		AllSegmentCount int8     `json:"allSegmentCount"`
		FileSegments    []SegDet `json:"fileSegments"`
	}{
		UAddress:        PEER_ADDRESS,
		FileName:        tMeta.ParentFileName,
		FileSize:        pFileInfo.Size(),
		FileImage:       "", //fix it later
		FileHash:        tMeta.ParentFilehash,
		FileExt:         tMeta.ParentFileExtention,
		AllSegmentCount: tMeta.SegmentCount,
		FileSegments:    allSegments,
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

	go getTrackerFile(trackerUrl)
}

func listenToTheOtherPeers(peerPort string) error {

	peerAddrss := fmt.Sprintf("http://127.0.0.1:%s", peerPort)

	listner, err := net.Listen("tcp", peerAddrss)

	if err != nil {
		log.Panic("ListenToPeers:", err.Error())
		return err
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

type SegDet struct {
	FileSegmentHash string `json:"fileSegmentHash"`
	SegmentNumber   int8   `json:"segNum"`
	ContentSize     int64  `json:"contentSize"`
	SegFileSize     int64  `json:"segFileSize"`
}

type SegmentFileAddress struct {
	FileSegmentHash string   `json:"fileSegmentHash"`
	SegContentSize  int64    `json:"segContentSize"`
	SegFileSize     int64    `json:"segFileSize"`
	SegmentNumber   int8     `json:"SegmentNumber"`
	SegAddress      []string `json:"segAddress"`
}

type TunnelTracerContent struct {
	FileHash         string               `json:"fileHash"`
	FileName         string               `json:"fileName"`
	FileSize         int64                `json:"fileSize"`
	AllSegmentCount  int8                 `json:"allSegmentCount"`
	FileExt          string               `json:"fileExt"`
	AllFileSegements []SegmentFileAddress `json:"fileSegments"`
}

func getTrackerContent() (trackerContent []TunnelTracerContent, err error) {

	rFile, err := os.Open(TRACKER_PATH)

	if err != nil {
		log.Println("GetTrackerContent: ", err.Error())
		return nil, err
	}

	defer rFile.Close()

	var trackerDetails []TunnelTracerContent

	rDecoder := json.NewDecoder(rFile)

	err = rDecoder.Decode(&trackerDetails)

	if err != nil {
		log.Println("GetTrackerContent: ", err.Error())
		return nil, err
	}

	return trackerDetails, nil
}

func requestSegments(parentFileHash string) error {

	trackerDetails, err := getTrackerContent()

	if err != nil {
		log.Println("RequestSegments: ", err.Error())
		return err
	}

	var segAdd []SegmentFileAddress

	for _, trackerDet := range trackerDetails {
		if trackerDet.FileHash == parentFileHash {
			for _, segHash := range trackerDet.AllFileSegements {
				segAdd = append(segAdd, segHash)
			}
		}
	}

	if len(segAdd) == 0 {
		return fmt.Errorf("there is no file in the tracker")
	}

	neededSeg := checkIfPeerHasSeg(segAdd)

	if len(neededSeg) == 0 {
		log.Println("All segments already excists")
		return nil
	}

	for _, nS := range neededSeg {
		if nS.SegAddress[0] != PEER_ADDRESS {
			go connecTopeer(nS)
		}
	}

	return nil
}

func checkIfPeerHasSeg(segMentsName []SegmentFileAddress) (neededSeg []SegmentFileAddress) {

	dirEntry, err := os.ReadDir(SEGEMENT_STORE_DIRECTORY)

	if err != nil {

		log.Panic("CheckPeerSeg:", err.Error())

	}

	var needSeg []SegmentFileAddress
	// dirE.Name()[:(len(dirE.Name())-3)]

	var fileName []string

	for _, dirE := range dirEntry {
		fileName = append(fileName, dirE.Name()[:(len(dirE.Name())-3)])
	}

	for _, seg := range segMentsName {
		if !slices.Contains(fileName, seg.FileSegmentHash) {
			needSeg = append(needSeg, seg)
		}
	}
	return needSeg

}

func connecTopeer(segDet SegmentFileAddress) {

	add := segDet.SegAddress[0]

	conn, err := net.Dial("tcp", add)

	if err != nil {
		log.Println("ConnecToPeer:", add, " is not online")
		return
	}

	defer conn.Close()

	req := fmt.Sprintf("REQUEST:%s\n", segDet.FileSegmentHash)

	conn.Write([]byte(req))

	requestFeedBackReader := bufio.NewReader(conn)

	headerFeedBack, err := requestFeedBackReader.ReadString('\n')

	if err != nil {
		fmt.Println("ConnectoPeer:", err.Error())
		return
	}

	// log.Println("peer: ", headerFeedBack)

	idStr := strings.TrimPrefix(headerFeedBack, "CHUNK:")
	log.Println("peer:", idStr)

	segmentReader := bufio.NewReader(conn)

	// segHeader, err := segmentReader.ReadString('\n')

	// if err != nil {
	// 	fmt.Println("ConnectoPeer:", err.Error())
	// 	return
	// }

	// log.Println("peer: ", segHeader)

	log.Println("ConnectTopeer: ", "SegSize: ", segDet.SegFileSize)

	segReadBytes := make([]byte, segDet.SegFileSize)

	_, err = io.ReadFull(segmentReader, segReadBytes)

	if err != nil {
		log.Println("ConnectTopeer: ", err.Error())
		return
	}

	// log.Printf("ConnectPeer :%x", segReadBytes[segLenght-1])

	segPath := fmt.Sprintf("%s/%s.bl", SEGEMENT_STORE_DIRECTORY, segDet.FileSegmentHash)

	seg, err := os.Create(segPath)

	if err != nil {
		log.Println("ConnetPeer :", err.Error())
		return
	}

	defer seg.Close()

	_, err = seg.Write(segReadBytes)

	if err != nil {
		log.Println("ConnectToPeer: ", err.Error())
		return
	}

	size, _ := seg.Stat()

	log.Println("ConnectToPeer :", segPath, " is written", "and size is: ", size.Size())
}

//[X] make a separate method for tracker json parse
//[ ] take all the .bl files and make the original file
//[ ] bind with the wails framework

func makeOriginalFile(parentFileHash string) error {

	allTracerContent, err := getTrackerContent()

	if err != nil {
		log.Println("MakeOriginalFile: ", err.Error())
		return err
	}

	var getParentFileDetails TunnelTracerContent

	for _, t := range allTracerContent {
		if t.FileHash == parentFileHash {
			getParentFileDetails = t
		}
	}

	var segments []TunnelSegment

	for _, s := range getParentFileDetails.AllFileSegements {
		segments = append(segments, TunnelSegment{
			FileSize:        s.SegFileSize,
			ContentSize:     s.SegContentSize,
			Filehash:        s.FileSegmentHash,
			FileDestination: fmt.Sprintf("%s/%s.bl", SEGEMENT_STORE_DIRECTORY, s.FileSegmentHash),
			SegmentNumber:   s.SegmentNumber,
		})
	}

	metaData := &TunnelSegmentFileMetadata{
		ParentFileName:      getParentFileDetails.FileName,
		ParentFileExtention: getParentFileDetails.FileExt,
		ParentFilehash:      getParentFileDetails.FileHash,
		ParentFileSize:      getParentFileDetails.FileSize,
		SegmentCount:        SEGMENT_SIZE,
		AllSegments:         segments,
	}

	err = jointBLFiles(metaData)

	if err != nil {
		return err
	}

	return nil
}
