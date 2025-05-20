package main

//constants for file segementation
const (
	SEGMENT_SIZE             int8   = 6
	SEGEMENT_MIN_FILE_SIZE   int64  = 1024
	SEGMENT_MAGIC_BYTES      string = "BLACKBOX"
	SEGEMENT_STORE_DIRECTORY string = "segments"
	JOINT_STORE_DIRECTORY    string = "made"
	ORIGINAL_FILE_STORAGE    string = "storage"
	SEGEMENT_EXT             string = ".bl"
)

//constants for peer httpRequest

const (
	TRACKER_PATH string = "tracker.json"
	TRACKER_URL  string = "http://127.0.0.1:8080"
	PEER_ADDRESS string = "127.0.0.1:6000"
)
