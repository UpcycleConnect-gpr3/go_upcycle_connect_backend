package project_models

import (
	"fmt"
	"go-upcycle_connect-backend/database"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"

	"github.com/google/uuid"
)

const TABLE = "PROJECTS"

type Project struct {
	Id          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	ImagePath   string `db:"image_path" json:"image_path"`
	UserId      string `db:"user_id" json:"user_id"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	UpdatedAt   string `db:"updated_at" json:"updated_at"`
}

type CreateProjectDTO struct {
	Name        string
	Description string
	ImagePath   string
	UserId      string
}

type UpdateProjectDTO struct {
	Name        string
	Description string
	ImagePath   string
	UserId      string
}

type ObjectSummary struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type StepSummary struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImagePath   string `json:"image_path"`
	ScheduledAt string `json:"scheduled_at"`
}

func (project *Project) Get(columns []string, by string, value any) error {
	return db.GetQuery[Project](database.UpcycleConnect, TABLE, columns, by, value, project)
}

func (project *Project) All(columns []string, dest *[]Project) error {
	return db.AllQuery[Project](database.UpcycleConnect, TABLE, columns, dest)
}

func CreateProject(dto CreateProjectDTO) *Project {
	action := fmt.Sprintf("INSERT INTO %s: %s", TABLE, dto.Name)
	res, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (name, description, image_path, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())",
		dto.Name, dto.Description, dto.ImagePath, dto.UserId,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	id, _ := res.LastInsertId()
	return &Project{Id: int(id)}
}

func UpdateProject(id int, dto UpdateProjectDTO) *Project {
	action := fmt.Sprintf("UPDATE %s WHERE ID: %d", TABLE, id)
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET name=?, description=?, image_path=?, user_id=?, updated_at=NOW() WHERE id=?",
		dto.Name, dto.Description, dto.ImagePath, dto.UserId, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &Project{Id: id}
}

func DeleteProject(id int) {
	action := fmt.Sprintf("DELETE FROM %s WHERE ID: %d", TABLE, id)
	_, err := database.UpcycleConnect.Exec("DELETE FROM "+TABLE+" WHERE id=?", id)
	if err != nil {
		log.Database(action, err)
	}
}

func GetProjectObjects(projectID int) []ObjectSummary {
	result := []ObjectSummary{}
	rows, err := database.UpcycleConnect.Query(
		"SELECT o.id, o.name FROM OBJECTS o JOIN OBJECT_PROJECT op ON o.id=op.object_id WHERE op.project_id=?",
		projectID,
	)
	if err != nil {
		log.Database("SELECT PROJECT OBJECTS", err)
		return result
	}
	defer rows.Close()
	for rows.Next() {
		o := ObjectSummary{}
		_ = rows.Scan(&o.Id, &o.Name)
		result = append(result, o)
	}
	return result
}

func LinkObject(projectID int, objectID string, userID string) {
	_, err := database.UpcycleConnect.Exec(
		"INSERT IGNORE INTO OBJECT_PROJECT (object_id, project_id) VALUES (?, ?)",
		objectID, projectID,
	)
	if err != nil {
		log.Database("LINK OBJECT TO PROJECT", err)
	}
}

func UnlinkObject(projectID int, objectID string) {
	_, err := database.UpcycleConnect.Exec(
		"DELETE FROM OBJECT_PROJECT WHERE object_id=? AND project_id=?",
		objectID, projectID,
	)
	if err != nil {
		log.Database("UNLINK OBJECT FROM PROJECT", err)
	}
}

func GetProjectSteps(projectID int) []StepSummary {
	result := []StepSummary{}
	rows, err := database.UpcycleConnect.Query(
		"SELECT s.id, s.name, s.description, s.image_path, s.scheduled_at FROM STEPS s JOIN PROJECT_STEP ps ON s.id=ps.step_id WHERE ps.project_id=?",
		projectID,
	)
	if err != nil {
		log.Database("SELECT PROJECT STEPS", err)
		return result
	}
	defer rows.Close()
	for rows.Next() {
		s := StepSummary{}
		_ = rows.Scan(&s.Id, &s.Name, &s.Description, &s.ImagePath, &s.ScheduledAt)
		result = append(result, s)
	}
	return result
}

func CreateProjectStep(projectID int, name, description, imagePath string, scheduledAt string) *StepSummary {
	stepId := uuid.New().String()
	_, err := database.UpcycleConnect.Exec(
		"INSERT INTO STEPS (id, name, description, image_path, scheduled_at) VALUES (?, ?, ?, ?, ?)",
		stepId, name, description, imagePath, scheduledAt,
	)
	if err != nil {
		log.Database("CREATE PROJECT STEP", err)
		return nil
	}
	_, err = database.UpcycleConnect.Exec(
		"INSERT IGNORE INTO PROJECT_STEP (project_id, step_id) VALUES (?, ?)",
		projectID, stepId,
	)
	if err != nil {
		log.Database("LINK STEP TO PROJECT", err)
	}
	return &StepSummary{Id: 0, Name: name, Description: description, ImagePath: imagePath, ScheduledAt: scheduledAt}
}

func LinkStep(projectID int, stepID string) error {
	_, err := database.UpcycleConnect.Exec(
		"INSERT IGNORE INTO PROJECT_STEP (project_id, step_id) VALUES (?, ?)",
		projectID, stepID,
	)
	if err != nil {
		log.Database("LINK STEP TO PROJECT", err)
		return err
	}
	return nil
}

func UnlinkStep(projectID int, stepID string) error {
	_, err := database.UpcycleConnect.Exec(
		"DELETE FROM PROJECT_STEP WHERE project_id=? AND step_id=?",
		projectID, stepID,
	)
	if err != nil {
		log.Database("UNLINK STEP FROM PROJECT", err)
		return err
	}
	return nil
}
