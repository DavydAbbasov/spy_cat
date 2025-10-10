package dto

import (
	"strings"
	"time"

	"github.com/DavydAbbasov/spy-cat/internal/domain"
)

type CreateMissionRequest struct {
	Title       string              `json:"title" validate:"required,min=3,max=128"`
	Description string              `json:"description"`
	Status      string              `json:"status" validate:"omitempty,oneof=planned active completed"`
	Goals       []CreateGoalRequest `json:"goals" validate:"dive"`
}

type CreateGoalRequest struct {
	Name    string `json:"name"    validate:"required,min=2,max=64"`
	Country string `json:"country" validate:"required,len=2"`
	Notes   string `json:"notes"   validate:"max=1000"`
}
type CreateMissionResponse struct {
	ID int64 `json:"id"`
}
type AssignMissionRequest struct {
	CatID *int64 `json:"catId" validate:"omitempty,gt=0"`
}
type MissionResponse struct {
	ID          int64          `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	CatID       *int64         `json:"catId,omitempty"`
	Goals       []GoalResponse `json:"goals,omitempty"`
	CreatedAt   string         `json:"createdAt"`
	UpdatedAt   string         `json:"updatedAt"`
}
type GoalResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Country   string `json:"country"`
	Notes     string `json:"notes"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// mapping
func ToCreateMissionParams(req CreateMissionRequest) domain.CreateMissionParams {
	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = string(domain.StatusPlanned)
	}

	return domain.CreateMissionParams{
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		Status:      domain.MissionStatus(status),
		Goals:       toCreateGoalParams(req.Goals),
	}
}

func toCreateGoalParams(in []CreateGoalRequest) []domain.CreateGoalParams {
	if len(in) == 0 {
		return nil
	}

	out := make([]domain.CreateGoalParams, 0, len(in))
	seen := make(map[string]struct{}, len(in))
	for _, g := range in {
		name := strings.TrimSpace(g.Name)
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, domain.CreateGoalParams{
			Name:    name,
			Country: strings.ToUpper(strings.TrimSpace(g.Country)),
			Notes:   strings.TrimSpace(g.Notes),
		})
	}
	return out
}
func ToMissionResponse(m domain.Mission, goals []domain.MissionGoal) MissionResponse {
	resp := MissionResponse{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Status:      string(m.Status),
		CatID:       m.CatID,
		CreatedAt:   m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   m.UpdatedAt.Format(time.RFC3339),
	}
	resp.Goals = make([]GoalResponse, 0, len(goals))
	for _, g := range goals {
		resp.Goals = append(resp.Goals, GoalResponse{
			ID:        g.ID,
			Name:      g.Name,
			Status:    string(g.Status),
			Country:   g.Country,
			Notes:     g.Notes,
			CreatedAt: g.CreatedAt.Format(time.RFC3339),
			UpdatedAt: g.UpdatedAt.Format(time.RFC3339),
		})
	}
	return resp
}
