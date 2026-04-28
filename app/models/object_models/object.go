package object_models

import (
	"fmt"
	"go-upcycle_connect-backend/database"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"

	"github.com/google/uuid"
)

const TABLE = "OBJECTS"

type Object struct {
	Id                    string  `db:"id" json:"id"`
	Name                  string  `db:"name" json:"name"`
	Description           string  `db:"description" json:"description"`
	Price                 float64 `db:"price" json:"price"`
	ImagePath             string  `db:"image_path" json:"image_path"`
	ColumnForCalcTheScore string  `db:"column_for_calc_the_score" json:"column_for_calc_the_score"`
	Quantity              int     `db:"quantity" json:"quantity"`
	UserId                string  `db:"user_id" json:"user_id"`
	Score                 float64 `db:"score" json:"score"`
	CreatedAt             string  `db:"created_at" json:"created_at"`
	UpdatedAt             string  `db:"updated_at" json:"updated_at"`
}

type CreateObjectDTO struct {
	Id                    string
	Name                  string
	Description           string
	Price                 float64
	ImagePath             string
	ColumnForCalcTheScore string
	Quantity              int
	UserId                string
	Score                 float64
}

type UpdateObjectDTO struct {
	Name                  string
	Description           string
	Price                 float64
	ImagePath             string
	ColumnForCalcTheScore string
	Quantity              int
	UserId                string
	Score                 float64
}

type DeliveryMethodSummary struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ProjectSummary struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type UserSummary struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type ScoreResponse struct {
	ObjectId string  `json:"object_id"`
	Score    float64 `json:"score"`
}

func (object *Object) Get(columns []string, by string, value any) error {
	return db.GetQuery[Object](database.UpcycleConnect, TABLE, columns, by, value, object)
}

func (object *Object) All(columns []string, dest *[]Object) error {
	return db.AllQuery[Object](database.UpcycleConnect, TABLE, columns, dest)
}

func CreateObject(dto CreateObjectDTO) *Object {
	action := fmt.Sprintf("INSERT INTO %s: %s", TABLE, dto.Name)
	objectId := dto.Id
	if objectId == "" {
		objectId = uuid.New().String()
	}
	_, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (id, name, description, price, image_path, column_for_calc_the_score, quantity, user_id, score, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())",
		objectId, dto.Name, dto.Description, dto.Price, dto.ImagePath, dto.ColumnForCalcTheScore, dto.Quantity, dto.UserId, dto.Score,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &Object{Id: objectId}
}

func UpdateObject(id string, dto UpdateObjectDTO) *Object {
	action := fmt.Sprintf("UPDATE %s WHERE ID: %s", TABLE, id)
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET name=?, description=?, price=?, image_path=?, column_for_calc_the_score=?, quantity=?, user_id=?, score=? WHERE id=?",
		dto.Name, dto.Description, dto.Price, dto.ImagePath, dto.ColumnForCalcTheScore, dto.Quantity, dto.UserId, dto.Score, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &Object{Id: id}
}

func DeleteObject(id string) {
	action := fmt.Sprintf("DELETE FROM %s WHERE ID: %s", TABLE, id)
	_, err := database.UpcycleConnect.Exec("DELETE FROM "+TABLE+" WHERE id=?", id)
	if err != nil {
		log.Database(action, err)
	}
}

func GetObjectScore(id string) *ScoreResponse {
	var o Object
	if err := o.Get([]string{"id", "score"}, "id", id); err != nil {
		return nil
	}
	return &ScoreResponse{
		ObjectId: o.Id,
		Score:    o.Score,
	}
}

func GetObjectDeliveryMethods(objectID string) []DeliveryMethodSummary {
	result := []DeliveryMethodSummary{}
	rows, err := database.UpcycleConnect.Query(
		"SELECT dm.id, dm.name FROM DELIVERY_METHODS dm JOIN OBJECT_DELIVERY_METHOD odm ON dm.id=odm.delivery_method_id WHERE odm.object_id=?",
		objectID,
	)
	if err != nil {
		log.Database("SELECT OBJECT DELIVERY METHODS", err)
		return result
	}
	defer rows.Close()
	for rows.Next() {
		d := DeliveryMethodSummary{}
		_ = rows.Scan(&d.Id, &d.Name)
		result = append(result, d)
	}
	return result
}

func LinkDeliveryMethod(objectID string, deliveryMethodID int) {
	_, err := database.UpcycleConnect.Exec(
		"INSERT IGNORE INTO OBJECT_DELIVERY_METHOD (object_id, delivery_method_id) VALUES (?, ?)",
		objectID, deliveryMethodID,
	)
	if err != nil {
		log.Database("LINK DELIVERY METHOD TO OBJECT", err)
	}
}

func UnlinkDeliveryMethod(objectID string, deliveryMethodID int) {
	_, err := database.UpcycleConnect.Exec(
		"DELETE FROM OBJECT_DELIVERY_METHOD WHERE object_id=? AND delivery_method_id=?",
		objectID, deliveryMethodID,
	)
	if err != nil {
		log.Database("UNLINK DELIVERY METHOD FROM OBJECT", err)
	}
}

func GetObjectProjects(objectID string) []ProjectSummary {
	result := []ProjectSummary{}
	rows, err := database.UpcycleConnect.Query(
		"SELECT p.id, p.name FROM PROJECTS p JOIN OBJECT_PROJECT op ON p.id=op.project_id WHERE op.object_id=?",
		objectID,
	)
	if err != nil {
		log.Database("SELECT OBJECT PROJECTS", err)
		return result
	}
	defer rows.Close()
	for rows.Next() {
		p := ProjectSummary{}
		_ = rows.Scan(&p.Id, &p.Name)
		result = append(result, p)
	}
	return result
}

func LinkProject(objectID string, projectID int) {
	_, err := database.UpcycleConnect.Exec(
		"INSERT IGNORE INTO OBJECT_PROJECT (object_id, project_id) VALUES (?, ?)",
		objectID, projectID,
	)
	if err != nil {
		log.Database("LINK PROJECT TO OBJECT", err)
	}
}

func UnlinkProject(objectID string, projectID int) {
	_, err := database.UpcycleConnect.Exec(
		"DELETE FROM OBJECT_PROJECT WHERE object_id=? AND project_id=?",
		objectID, projectID,
	)
	if err != nil {
		log.Database("UNLINK PROJECT FROM OBJECT", err)
	}
}

func GetObjectUsers(objectID string) []UserSummary {
	result := []UserSummary{}
	rows, err := database.UpcycleConnect.Query(
		"SELECT u.id, u.username FROM USERS u JOIN OBJECT_USER ou ON u.id=ou.user_id WHERE ou.object_id=?",
		objectID,
	)
	if err != nil {
		log.Database("SELECT OBJECT USERS", err)
		return result
	}
	defer rows.Close()
	for rows.Next() {
		u := UserSummary{}
		_ = rows.Scan(&u.Id, &u.Username)
		result = append(result, u)
	}
	return result
}

func LinkUser(objectID string, userID string) error {
	_, err := database.UpcycleConnect.Exec(
		"INSERT IGNORE INTO OBJECT_USER (object_id, user_id) VALUES (?, ?)",
		objectID, userID,
	)
	if err != nil {
		log.Database("LINK USER TO OBJECT", err)
		return err
	}
	return nil
}

func UnlinkUser(objectID string, userID string) error {
	_, err := database.UpcycleConnect.Exec(
		"DELETE FROM OBJECT_USER WHERE object_id=? AND user_id=?",
		objectID, userID,
	)
	if err != nil {
		log.Database("UNLINK USER FROM OBJECT", err)
		return err
	}
	return nil
}
