package user_models

import (
	"fmt"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/var/database"

	"github.com/google/uuid"
)

const TABLE = "USERS"

type User struct {
	Id        uuid.UUID `db:"id" json:"id"`
	Firstname *string   `db:"firstname" json:"firstname"`
	Lastname  *string   `db:"lastname" json:"lastname"`
	Email     *string   `db:"email" json:"email"`
	CreatedAt string    `db:"created_at" json:"created_at"`
	UpdatedAt string    `db:"updated_at" json:"updated_at"`
}

type CreateUserDTO struct {
	Id        string
	Firstname string
	Lastname  string
	Email     string
}
type UpdateUserDTO struct {
	Firstname string
	Lastname  string
	Email     string
}

func (user *User) Get(columns []string, by string, value ...any) error {
	return db.GetQuery[User](database.UpcycleConnect, TABLE, user, columns, by, value...)
}

func (user *User) All(columns []string, dest *[]User) error {
	return db.AllQuery[User](database.UpcycleConnect, TABLE, columns, dest)
}

func (user *User) Create(dto CreateUserDTO) error {
	user.Id, _ = uuid.Parse(dto.Id)
	user.Email = &dto.Email
	user.Firstname = &dto.Firstname
	user.Lastname = &dto.Lastname
	columns := []string{"id", "email", "firstname", "lastname"}
	return db.CreateQuery(database.UpcycleConnect, TABLE, columns, user.Id.String(), user.Email, user.Firstname, user.Lastname)
}

func (user *User) Update(userDto UpdateUserDTO, userId string) error {
	parsedId, err := uuid.Parse(userId)
	if err != nil {
		return err
	}
	user.Id = parsedId

	user.Firstname = &userDto.Firstname
	user.Lastname = &userDto.Lastname
	user.Email = &userDto.Email

	return db.UpdateQuery[*User](database.UpcycleConnect, TABLE, user, "id")
}

func DeleteUser(id uuid.UUID) {
	action := fmt.Sprintf("DELETE FROM %s WHERE id: %s", TABLE, id.String())
	_, err := database.UpcycleConnect.Exec("DELETE FROM "+TABLE+" WHERE id=?", id)
	if err != nil {
		log.Database(action, err)
	}
}
