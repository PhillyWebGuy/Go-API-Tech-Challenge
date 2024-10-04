package handlers

import (
	"gorm.io/gorm"
)

type RequestHandler struct {
	DB *gorm.DB
}

func NewRequestHandler(db *gorm.DB) *RequestHandler {
	return &RequestHandler{DB: db}
}
