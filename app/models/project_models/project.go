package project_models

import (
	"fmt"
	"go-upcycle_connect-backend/database"
	"go-upcycle_connect-backend/utils/db"
	"go-upcycle_connect-backend/utils/log"
	"time"
)

const TABLE = "PROJECTS"

type Project struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Score       int       `json:"score"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type CreateProjectDTO struct {
	Name        string
	Description string
}

type UpdateProjectDTO struct {
	Name        string
	Description string
}

type ObjectSummary struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type StepSummary struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Order int    `json:"order"`
}

type ScoreBreakdown struct {
	Complexity int `json:"complexity"`
	Impact     int `json:"impact"`
}

type ScoreResponse struct {
	ProjectId int            `json:"project_id"`
	Score     int            `json:"score"`
	Breakdown ScoreBreakdown `json:"breakdown"`
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
		"INSERT INTO "+TABLE+" (name, description) VALUES (?, ?)",
		dto.Name, dto.Description,
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
		"UPDATE "+TABLE+" SET name=?, description=? WHERE id=?",
		dto.Name, dto.Description, id,
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

func GetProjectScore(id int) *ScoreResponse {
	var p Project
	if err := p.Get([]string{"id", "score"}, "id", id); err != nil {
		return nil
	}
	complexity := p.Score / 2
	impact := p.Score - complexity
	return &ScoreResponse{
		ProjectId: p.Id,
		Score:     p.Score,
		Breakdown: ScoreBreakdown{Complexity: complexity, Impact: impact},
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

func LinkObject(projectID, objectID int) {
	_, err := database.UpcycleConnect.Exec(
		"INSERT IGNORE INTO OBJECT_PROJECT (object_id, project_id) VALUES (?, ?)",
		objectID, projectID,
	)
	if err != nil {
		log.Database("LINK OBJECT TO PROJECT", err)
	}
}

func GetProjectSteps(projectID int) []StepSummary {
	result := []StepSummary{}
	rows, err := database.UpcycleConnect.Query(
		"SELECT s.id, s.title, s.`order` FROM STEPS s JOIN PROJECT_STEP ps ON s.id=ps.step_id WHERE ps.project_id=?",
		projectID,
	)
	if err != nil {
		log.Database("SELECT PROJECT STEPS", err)
		return result
	}
	defer rows.Close()
	for rows.Next() {
		s := StepSummary{}
		_ = rows.Scan(&s.Id, &s.Title, &s.Order)
		result = append(result, s)
	}
	return result
}

func CreateProjectStep(projectID int, title string, order int) *StepSummary {
	res, err := database.UpcycleConnect.Exec(
		"INSERT INTO STEPS (title, `order`) VALUES (?, ?)",
		title, order,
	)
	if err != nil {
		log.Database("CREATE PROJECT STEP", err)
		return nil
	}
	stepID, _ := res.LastInsertId()
	_, err = database.UpcycleConnect.Exec(
		"INSERT IGNORE INTO PROJECT_STEP (project_id, step_id) VALUES (?, ?)",
		projectID, int(stepID),
	)
	if err != nil {
		log.Database("LINK STEP TO PROJECT", err)
	}
	return &StepSummary{Id: int(stepID), Title: title, Order: order}
}
