package handlers

import (
	"github.com/alireza-akbarzadeh/ginflow/internal/repository"
)

// Handler holds all dependencies for HTTP handlers
type Handler struct {
	Repos     *repository.Models
	JWTSecret string
}

// NewHandler creates a new Handler instance
func NewHandler(repos *repository.Models, jwtSecret string) *Handler {
	return &Handler{
		Repos:     repos,
		JWTSecret: jwtSecret,
	}
}
