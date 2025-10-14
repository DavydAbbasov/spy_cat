package servieserrors

import "errors"

var (
	ErrMissionNotFound = errors.New("mission not found")
	ErrGoalNotFound    = errors.New("mission goal not found")
)
var (
	ErrMissionAlreadyCompleted = errors.New("mission already completed")
	ErrMissionNotPlanned       = errors.New("mission is not in planned status")
	ErrMissionNotActive        = errors.New("mission is not in active status")
	ErrMissionHasAssignee      = errors.New("mission is assigned to a cat")
	ErrMissionAlreadyExists    = errors.New("mission with same title already exists")
)
var (
	ErrGoalAlreadyDone     = errors.New("goal already done")
	ErrGoalDeleteForbidden = errors.New("cannot delete a completed goal")
)
var (
	ErrInvalidCreateMission = errors.New("create mission invalid")
	ErrInvalidGoalName      = errors.New("goals name is invalid")
	ErrInvalidCountry       = errors.New("counrty is invalid")
	ErrInvalidStatus        = errors.New("invalid status")
	ErrInvalidTransition    = errors.New("invalid transition")
	ErrConflict             = errors.New("ststus conflict")
)
