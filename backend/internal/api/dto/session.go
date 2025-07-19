package dto

import "time"

type CreateSessionRequest struct {
	ID          string `json:"id" binding:"required"`
	User        string `json:"user" binding:"required"`
	Environment string `json:"environment" binding:"required"`
	Ticker      string `json:"ticker" binding:"required"`
}

type CreateSessionResponse struct {
	Status string `json:"status"`
}

type EndSessionRequest struct {
	ID string `json:"id" binding:"required"`
}

type EndSessionResponse struct {
	Status string `json:"status"`
}

type SessionItemResponse struct {
	ID          string     `json:"id"`
	User        string     `json:"user"`
	Environment string     `json:"environment"`
	Ticker      string     `json:"ticker"`
	StartedAt   time.Time  `json:"started_at"`
	EndedAt     *time.Time `json:"ended_at,omitempty"`
}

type SessionListResponse struct {
	Sessions []SessionItemResponse `json:"sessions"`
}
