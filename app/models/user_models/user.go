package user_models

import (
	"fmt"
	"go-upcycle_connect-backend/database"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"

	"github.com/google/uuid"
)

const TABLE = "USERS"

type User struct {
	Id        uuid.UUID `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Firstname string    `db:"firstname" json:"firstname"`
	Lastname  string    `db:"lastname" json:"lastname"`
	Email     string    `db:"email" json:"email"`
	CreatedAt string    `db:"created_at" json:"created_at"`
	UpdatedAt string    `db:"updated_at" json:"updated_at"`
}

type CreateUserDTO struct {
	Username  string
	Firstname string
	Lastname  string
	Email     string
}

type UpdateUserDTO struct {
	Username  string
	Firstname string
	Lastname  string
	Email     string
}

func (user *User) Get(columns []string, by string, value any) error {
	return db.GetQuery[User](database.UpcycleConnect, TABLE, columns, by, value, user)
}

func (user *User) All(columns []string, dest *[]User) error {
	return db.AllQuery[User](database.UpcycleConnect, TABLE, columns, dest)
}

func CreateUser(dto CreateUserDTO) *User {
	action := fmt.Sprintf("INSERT INTO %s: %s", TABLE, dto.Username)
	id := uuid.New()
	_, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (id, username, firstname, lastname, email) VALUES (?, ?, ?, ?, ?)",
		id, dto.Username, dto.Firstname, dto.Lastname, dto.Email,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &User{Id: id}
}

func UpdateUser(id uuid.UUID, dto UpdateUserDTO) *User {
	action := fmt.Sprintf("UPDATE %s WHERE id: %s", TABLE, id.String())
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET username=?, firstname=?, lastname=?, email=? WHERE id=?",
		dto.Username, dto.Firstname, dto.Lastname, dto.Email, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &User{Id: id}
}

func DeleteUser(id uuid.UUID) {
	action := fmt.Sprintf("DELETE FROM %s WHERE id: %s", TABLE, id.String())
	_, err := database.UpcycleConnect.Exec("DELETE FROM "+TABLE+" WHERE id=?", id)
	if err != nil {
		log.Database(action, err)
	}
}
