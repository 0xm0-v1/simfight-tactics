package units

import (
	"github.com/google/uuid"
)

type Unit struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Cost   int       `json:"cost"`
	Traits []string  `json:"traits"`
	Roles  []string  `json:"roles"`
	Stats  Stats     `json:"stats"`
}
