package domain

import (
	"time"
)

type MissionStatus string

const (
	StatusPlanned   MissionStatus = "planned"
	StatusActive    MissionStatus = "active"
	StatusCompleted MissionStatus = "completed"
)

type Mission struct {
	ID          int64
	Title       string
	Description string
	Status      MissionStatus
	CatID       *int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
type MissionGoalStatus string

const (
	GoalTodo MissionGoalStatus = "todo"
	GoalDone MissionGoalStatus = "done"
)

type MissionGoal struct {
	ID        int64
	MissionID int64
	Name      string
	Country   string
	Notes     string
	Status    MissionGoalStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}
type CreateMissionParams struct {
	Title       string
	Description string
	Status      MissionStatus
	Goals       []CreateGoalParams
}
type CreateGoalParams struct {
	Name    string
	Country string
	Notes   string
}
type MissionFilter struct {
	Status *MissionStatus
	CatID  *int64
	Q      *string
	//pagination
	Limit  int
	Offset int
}

type MissionListItem struct {
	ID        int64
	Title     string
	Status    MissionStatus
	CatID     *int64
	CreatedAt time.Time
}
type UpdateMissionStatusParams struct {
	ID     int64
	Status MissionStatus
}

func CanTransition(from, to MissionStatus) bool {
	switch from {

	case StatusPlanned:
		return to == StatusActive

	case StatusActive:
		return to == StatusCompleted

	case StatusCompleted:
		return false

	default:
		return false
	}
}
