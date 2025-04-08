package main

const (
	SEGMENT_SIZE             int8   = 6
	SEGEMENT_MIN_FILE_SIZE   int64  = 1024
	SEGMENT_MAGIC_BYTES      string = "BLACKBOX"
	SEGEMENT_STORE_DIRECTORY string = "segments"
	JOINT_STORE_DIRECTORY    string = "made"
	SEGEMENT_EXT             string = ".bl"

	DATABASE_DIRECTORY   string = "database"
	DATABASE_DRIVER_NAME string = "sqlite3"
	DATABASE_PATH        string = "database/tunnel.db"
)
