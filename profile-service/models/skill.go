package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Skill struct {
	ID        *uuid.UUID `json:"id"`
	ProfileID *uuid.UUID `json:"profile_id"`
	Skill     string     `json:"skill"`
	Detail    string     `json:"detail"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (Skill) TableName() string {
	return "skill"
}

func (s *Skill) GenUUID() {
	id, _ := uuid.NewV4()
	s.ID = &id
}

func (s *Skill) SetCreatedAt() {
	now := time.Now()
	s.CreatedAt = &now
}

func (s *Skill) SetUpdatedAt() {
	now := time.Now()
	s.UpdatedAt = &now
}