package object_models

import (
	"fmt"
	"go-upcycle_connect-backend/database"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"
)

const TABLE = "OBJECTS"

type Object struct {
	Id             int    `db:"id" json:"id"`
	Name           string `db:"name" json:"name"`
	Material       string `db:"material" json:"material"`
	Condition      string `db:"condition" json:"condition"`
	Description    string `db:"description" json:"description"`
	UpcyclingScore int    `db:"upcycling_score" json:"upcycling_score"`
	CreatedAt      string `db:"created_at" json:"created_at"`
	UpdatedAt      string `db:"updated_at" json:"updated_at"`
}

type CreateObjectDTO struct {
	Name        string
	Material    string
	Condition   string
	Description string
}

type UpdateObjectDTO struct {
	Name        string
	Material    string
	Condition   string
	Description string
}

type DeliveryMethodSummary struct {
	Id   int     `json:"id"`
	Name string  `json:"name"`
	Cost float64 `json:"cost"`
}

type ProjectSummary struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type UserSummary struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type ScoreBreakdown struct {
	Material    int `json:"material"`
	Condition   int `json:"condition"`
	Reusability int `json:"reusability"`
}

type ScoreResponse struct {
	ObjectId       int            `json:"object_id"`
	UpcyclingScore int            `json:"upcycling_score"`
	Breakdown      ScoreBreakdown `json:"breakdown"`
}

func (object *Object) Get(columns []string, by string, value any) error {
	return db.GetQuery[Object](database.UpcycleConnect, TABLE, columns, by, value, object)
}

func (object *Object) All(columns []string, dest *[]Object) error {
	return db.AllQuery[Object](database.UpcycleConnect, TABLE, columns, dest)
}

func CreateObject(dto CreateObjectDTO) *Object {
	action := fmt.Sprintf("INSERT INTO %s: %s", TABLE, dto.Name)
	score := 0
	if dto.Material != "" {
		score += 30
	}
	if dto.Condition == "good" || dto.Condition == "excellent" {
		score += 25
	}
	score += 23 // reusability base
	res, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (name, material, `condition`, description, upcycling_score) VALUES (?, ?, ?, ?, ?)",
		dto.Name, dto.Material, dto.Condition, dto.Description, score,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	id, _ := res.LastInsertId()
	return &Object{Id: int(id)}
}

func UpdateObject(id int, dto UpdateObjectDTO) *Object {
	action := fmt.Sprintf("UPDATE %s WHERE ID: %d", TABLE, id)
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET name=?, material=?, `condition`=?, description=? WHERE id=?",
		dto.Name, dto.Material, dto.Condition, dto.Description, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &Object{Id: id}
}

func DeleteObject(id int) {
	action := fmt.Sprintf("DELETE FROM %s WHERE ID: %d", TABLE, id)
	_, err := database.UpcycleConnect.Exec("DELETE FROM "+TABLE+" WHERE id=?", id)
	if err != nil {
		log.Database(action, err)
	}
}

func GetObjectScore(id int) *ScoreResponse {
	var o Object
	if err := o.Get([]string{"id", "material", "`condition`", "upcycling_score"}, "id", id); err != nil {
		return nil
	}
	matScore := 0
	condScore := 0
	if o.Material != "" {
		matScore = 30
	}
	if o.Condition == "good" || o.Condition == "excellent" {
		condScore = 25
	}
	reuse := o.UpcyclingScore - matScore - condScore
	if reuse < 0 {
		reuse = 0
	}
	return &ScoreResponse{
		ObjectId:       o.Id,
		UpcyclingScore: o.UpcyclingScore,
		Breakdown:      ScoreBreakdown{Material: matScore, Condition: condScore, Reusability: reuse},
	}
}

func GetObjectDeliveryMethods(objectID int) []DeliveryMethodSummary {
	result := []DeliveryMethodSummary{}
	rows, err := database.UpcycleConnect.Query(
		"SELECT dm.id, dm.name, dm.cost FROM DELIVERY_METHODS dm JOIN OBJECT_DELIVERY_METHOD odm ON dm.id=odm.delivery_method_id WHERE odm.object_id=?",
		objectID,
	)
	if err != nil {
		log.Database("SELECT OBJECT DELIVERY METHODS", err)
		return result
	}
	defer rows.Close()
	for rows.Next() {
		d := DeliveryMethodSummary{}
		_ = rows.Scan(&d.Id, &d.Name, &d.Cost)
		result = append(result, d)
	}
	return result
}

func LinkDeliveryMethod(objectID, deliveryMethodID int) {
	_, err := database.UpcycleConnect.Exec(
		"INSERT IGNORE INTO OBJECT_DELIVERY_METHOD (object_id, delivery_method_id) VALUES (?, ?)",
		objectID, deliveryMethodID,
	)
	if err != nil {
		log.Database("LINK DELIVERY METHOD TO OBJECT", err)
	}
}

func UnlinkDeliveryMethod(objectID, deliveryMethodID int) {
	_, err := database.UpcycleConnect.Exec(
		"DELETE FROM OBJECT_DELIVERY_METHOD WHERE object_id=? AND delivery_method_id=?",
		objectID, deliveryMethodID,
	)
	if err != nil {
		log.Database("UNLINK DELIVERY METHOD FROM OBJECT", err)
	}
}

func GetObjectProjects(objectID int) []ProjectSummary {
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

func LinkProject(objectID, projectID int) {
	_, err := database.UpcycleConnect.Exec(
		"INSERT IGNORE INTO OBJECT_PROJECT (object_id, project_id) VALUES (?, ?)",
		objectID, projectID,
	)
	if err != nil {
		log.Database("LINK PROJECT TO OBJECT", err)
	}
}

func UnlinkProject(objectID, projectID int) {
	_, err := database.UpcycleConnect.Exec(
		"DELETE FROM OBJECT_PROJECT WHERE object_id=? AND project_id=?",
		objectID, projectID,
	)
	if err != nil {
		log.Database("UNLINK PROJECT FROM OBJECT", err)
	}
}

func GetObjectUsers(objectID int) []UserSummary {
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

func LinkUser(objectID int, userID string) {
	_, err := database.UpcycleConnect.Exec(
		"INSERT IGNORE INTO OBJECT_USER (object_id, user_id) VALUES (?, ?)",
		objectID, userID,
	)
	if err != nil {
		log.Database("LINK USER TO OBJECT", err)
	}
}

func UnlinkUser(objectID int, userID string) {
	_, err := database.UpcycleConnect.Exec(
		"DELETE FROM OBJECT_USER WHERE object_id=? AND user_id=?",
		objectID, userID,
	)
	if err != nil {
		log.Database("UNLINK USER FROM OBJECT", err)
	}
}
