package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type SqlDb struct {
	Db *sql.DB
}

func initDataBase() (*SqlDb, error) {

	if _, err := os.Stat(DATABASE_DIRECTORY); err != nil {
		err = os.Mkdir(DATABASE_DIRECTORY, os.ModeDir)

		if err != nil {
			return nil, err
		}

	}

	db, err := sql.Open(DATABASE_DRIVER_NAME, DATABASE_PATH)

	if err != nil {
		log.Println("Cannot open sql driver:", err.Error())
		return nil, err
	}

	err = isTableExcist(db, "user")

	if err != nil {

		if _, err := db.Exec(DATABASE_USER_TABLE); err != nil {
			return nil, err
		}
	}

	err = isTableExcist(db, "content")

	if err != nil {

		if _, err := db.Exec(DATABASE_CONTENT_TABLE); err != nil {
			return nil, err
		}
	}

	err = isTableExcist(db, "segment")

	if err != nil {

		if _, err := db.Exec(DATABASE_SEGMENT_TABLE); err != nil {
			return nil, err
		}
	}

	return &SqlDb{
		Db: db,
	}, nil

}

func (sqlDb *SqlDb) closeDataBase() {
	sqlDb.Db.Close()
}

func isTableExcist(db *sql.DB, tableName string) error {

	query := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s';", tableName)

	_, err := db.Exec(query)

	return err
}

type TunnelUser struct {
	Id        uuid.UUID `json:"id"`
	UserName  string    `json:"userName"`
	Password  string    `json:"password"`
	UserImage string    `json:"userImage"`
}

func (sqlDb *SqlDb) InsertUser(user TunnelUser) error {

	userInsert := `INSERT INTO user (id,name, password, image) VALUES (?,?,?,?);`

	_, err := sqlDb.Db.Exec(userInsert, user.Id, user.UserName, user.Password, user.UserImage)

	if err != nil {
		return err
	}

	return nil
}

func (sqlDb *SqlDb) QueryUser(userName string) (*TunnelUser, error) {
	userQuery := `SELECT * FROM user WHERE name = ?;`

	res, err := sqlDb.Db.Query(userQuery, userName)

	if err != nil {
		return nil, err
	}

	var id, name, password, image string

	for res.Next() {
		err = res.Scan(&id, &name, &password, &image)
		if err != nil {
			log.Panic("QueryUser:", err.Error())
			return nil, err
		}
	}

	uid, err := uuid.Parse(id)

	if err != nil {
		return nil, err
	}

	return &TunnelUser{
		Id:        uid,
		UserName:  name,
		Password:  password,
		UserImage: image,
	}, nil

}

// `CREATE TABLE content (
//     id TEXT NOT NULL,
//     uid TEXT NOT NULL,
//     fileName TEXT NOT NULL,
//     fileSize INTEGER NOT NULL,
//     fileImage TEXT NOT NULL,
// 	fileHash TEXT NOT NULL,
//     FOREIGN KEY (uid) REFERENCES user(id));`

type TunnelContent struct {
	Cid       uuid.UUID `json:"id"`
	Uid       uuid.UUID `json:"uid"`
	FileName  string    `json:"fileName"`
	FileSize  string    `json:"fileSize"`
	FileImage string    `json:"fileImage"`
	FileHash  string    `json:"fileHash"`
}

func (sqlDb *SqlDb) InsertNewContentInformation(tunnelContent *TunnelContent) error {

	contentInsert := `INSERT INTO content (id,uid, fileName, fileSize, fileImage, fileHash) 
	VALUES (?,?,?,?,?,?);`

	_, err := sqlDb.Db.Exec(contentInsert,
		tunnelContent.Cid,
		tunnelContent.Uid,
		tunnelContent.FileName,
		tunnelContent.FileSize,
		tunnelContent.FileImage,
		tunnelContent.FileHash)

	if err != nil {
		return err
	}

	return nil

}
