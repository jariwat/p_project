package repository

import (
	"gorm.io/gorm"
)

type profileRepository struct {
	client *gorm.DB
}

func NewPsqlProfileRepository(client *gorm.DB) *profileRepository {
	return &profileRepository{
		client: client,
	}
}