package user_models

import (
	"database/sql"
	"fmt"
	"go-upcycle_connect-backend/database"
	"go-upcycle_connect-backend/utils/log"
	"time"

	"github.com/google/uuid"
)

const (
	TABLE = "USERS"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func GetUserByEmail(email string) *User {
	user := User{}
	action := fmt.Sprintf("SELECT USER WHERE USERNAME : %s", email)

	row := database.UpcycleConnect.QueryRow("SELECT id, email FROM "+TABLE+" WHERE email = ?", email)

	err := row.Scan(&user.Id, &user.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Database(action, err)
		return nil
	}

	if err = row.Err(); err != nil {
		log.Database(action, err)
		return nil
	}

	return &user
}

func CreateUser(user User) {
	action := fmt.Sprintf("INSERT INTO "+TABLE+" : %s", user.Email)

	_, err := database.UpcycleConnect.Exec("INSERT INTO "+TABLE+" (id, email) VALUES (?, ?)", "uuid.New().String()", user.Email)

	if err != nil {
		log.Database(action, err)
	}
}
