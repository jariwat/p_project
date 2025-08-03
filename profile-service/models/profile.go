package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Gender string

const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
)

type Profile struct {
	ID         *uuid.UUID `json:"id"`
	FirstName  string     `json:"first_name"`
	MiddleName *string    `json:"middle_name"`
	LastName   string     `json:"last_name"`
	Gender     Gender     `json:"gender"`
	Class      string     `json:"class"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`

	Skills []*Skill `json:"skills"`
}

func (Profile) TableName() string {
	return "profile"
}

func (p *Profile) GenUUID() {
	id, _ := uuid.NewV4()
	p.ID = &id

}

func (p *Profile) SetCreatedAt() {
	now := time.Now()
	p.CreatedAt = &now
}

func (p *Profile) SetUpdatedAt() {
	now := time.Now()
	p.UpdatedAt = &now
}
