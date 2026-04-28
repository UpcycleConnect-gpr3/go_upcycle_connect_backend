package step_models

import (
	"fmt"
	"go-upcycle_connect-backend/database"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"

	"github.com/google/uuid"
)

const TABLE = "STEPS"

type Step struct {
	Id          string `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	ImagePath   string `db:"image_path" json:"image_path"`
	UserId      string `db:"user_id" json:"user_id"`
	ProjectId   int    `db:"project_id" json:"project_id"`
	ScheduledAt string `db:"scheduled_at" json:"scheduled_at"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	UpdatedAt   string `db:"updated_at" json:"updated_at"`
}

type CreateStepDTO struct {
	Name        string
	Description string
	ImagePath   string
	UserId      string
	ProjectId   int
	ScheduledAt string
}

type UpdateStepDTO struct {
	Name        string
	Description string
	ImagePath   string
	UserId      string
	ProjectId   int
	ScheduledAt string
}

func (step *Step) Get(columns []string, by string, value any) error {
	return db.GetQuery[Step](database.UpcycleConnect, TABLE, columns, by, value, step)
}

func (step *Step) All(columns []string, dest *[]Step) error {
	return db.AllQuery[Step](database.UpcycleConnect, TABLE, columns, dest)
}

func CreateStep(dto CreateStepDTO) *Step {
	action := fmt.Sprintf("INSERT INTO %s: %s", TABLE, dto.Name)
	stepId := uuid.New().String()
	_, err := database.UpcycleConnect.Exec(
		"INSERT INTO "+TABLE+" (id, name, description, image_path, user_id, project_id, scheduled_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())",
		stepId, dto.Name, dto.Description, dto.ImagePath, dto.UserId, dto.ProjectId, dto.ScheduledAt,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &Step{Id: stepId}
}

func UpdateStep(id string, dto UpdateStepDTO) *Step {
	action := fmt.Sprintf("UPDATE %s WHERE ID: %s", TABLE, id)
	_, err := database.UpcycleConnect.Exec(
		"UPDATE "+TABLE+" SET name=?, description=?, image_path=?, user_id=?, project_id=?, scheduled_at=?, updated_at=NOW() WHERE id=?",
		dto.Name, dto.Description, dto.ImagePath, dto.UserId, dto.ProjectId, dto.ScheduledAt, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &Step{Id: id}
}

func DeleteStep(id string) {
	action := fmt.Sprintf("DELETE FROM %s WHERE ID: %s", TABLE, id)
	_, err := database.UpcycleConnect.Exec("DELETE FROM "+TABLE+" WHERE id=?", id)
	if err != nil {
		log.Database(action, err)
	}
}
