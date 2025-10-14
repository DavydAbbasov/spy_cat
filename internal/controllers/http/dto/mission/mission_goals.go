package dto

import (
	"strings"
	"time"

	"github.com/DavydAbbasov/spy-cat/internal/domain"
)

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
type AssignMissionRequest struct {
	CatID *int64 `json:"catId" validate:"omitempty,gt=0"`
}
type CreateMissionResponse struct {
	ID int64 `json:"id"`
}
type GetMissionsQuery struct {
	Status *string `form:"status" binding:"omitempty,oneof=planned active completed"`
	CatID  *int64  `form:"catId"  binding:"omitempty,gt=0"`
	Q      *string `form:"q"      binding:"omitempty,min=1,max=128"`
	Limit  int     `form:"limit,default=10"  binding:"min=1,max=200"`
	Offset int     `form:"offset,default=0"  binding:"min=0"`
}
type GetMissionsResponse struct {
	Items  []MissionListItem `json:"items"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
	Total  int               `json:"total"`
}
type MissionListItem struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	CatID     *int64 `json:"catId,omitempty"`
	CreatedAt string `json:"createdAt"`
}
type UpdateMissionStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=planned active completed"`
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
func ToGetMissionsResponse(items []domain.MissionListItem, limit, offset, total int) GetMissionsResponse {
	out := make([]MissionListItem, 0, len(items))

	for _, it := range items {
		out = append(out, MissionListItem{
			ID:        it.ID,
			Title:     it.Title,
			Status:    string(it.Status),
			CatID:     it.CatID,
			CreatedAt: it.CreatedAt.Format(time.RFC3339),
		})
	}

	return GetMissionsResponse{
		Items:  out,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}
}
