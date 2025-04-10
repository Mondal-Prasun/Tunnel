package main

//constants for file segementation
const (
	SEGMENT_SIZE             int8   = 6
	SEGEMENT_MIN_FILE_SIZE   int64  = 1024
	SEGMENT_MAGIC_BYTES      string = "BLACKBOX"
	SEGEMENT_STORE_DIRECTORY string = "segments"
	JOINT_STORE_DIRECTORY    string = "made"
	SEGEMENT_EXT             string = ".bl"
)

//contstant for database
const (
	DATABASE_DIRECTORY   string = "database"
	DATABASE_DRIVER_NAME string = "sqlite3"
	DATABASE_PATH        string = "database/tunnel.db"

	DATABASE_USER_TABLE string = `CREATE TABLE user(
		id TEXT NOT NULL,
		name TEXT NOT NULL,
		password TEXT NOT NULL,
		image TEXT NOT NULL);`

	DATABASE_CONTENT_TABLE string = `CREATE TABLE content (
    id TEXT NOT NULL,
    uid TEXT NOT NULL,
    fileName TEXT NOT NULL,
    fileSize INTEGER NOT NULL,
    fileImage TEXT NOT NULL,
	fileHash TEXT NOT NULL,
    FOREIGN KEY (uid) REFERENCES user(id));`

	DATABASE_SEGMENT_TABLE string = `CREATE TABLE segment (
		id TEXT NOT NULL,
		cid TEXT NOT NULL,
		segmentNumber INTEGER NOT NULL,
		segmentsSize INTEGER NOT NULL,
		segmentName TEXT NOT NULL,
		uid TEXT NOT NULL,
		parentFileHash TEXT NOT NULL,
		FOREIGN KEY (cid) REFERENCES content(id),
		FOREIGN KEY (uid) REFERENCES user(id));`
)

//contstants for http handlers and server

const (
	STATUS_RESPONSE_OK           int    = 200
	STATUS_RESPONSE_ERROR        int    = 403
	STATUS_RESPONSE_NOT_FOUNR    int    = 404
	STATUS_RESPONSE_SERVER_ERROR int    = 500
	SERVER_PORT                  string = ":8080"
)
