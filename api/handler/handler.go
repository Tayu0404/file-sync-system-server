package handler

import (
	"github.com/Tayu0404/file-sync-system-server/api/model"
)

type handler struct {
	Model model.Database
}

func NewHandler(d model.Database) {
	return &handler {
		Model: d,
	}
}