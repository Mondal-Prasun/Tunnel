package main

import "github.com/google/uuid"

type TunnelUser struct {
	Id        uuid.UUID `json:"id"`
	UserName  string    `json:"userName"`
	Password  string    `json:"password"`
	UserImage string    `json:"userImage"`
}

func (sqlDb *SqlDb) insertUser(user TunnelUser) error {

	userInsert := `INSERT INTO user (id,name, password, image) VALUES (?,?,?,?);`

	_, err := sqlDb.Db.Exec(userInsert, user.Id, user.UserName, user.Password, user.UserImage)

	if err != nil {
		return err
	}

	return nil

}
