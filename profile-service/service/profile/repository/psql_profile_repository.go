package repository

import (
	"gorm.io/gorm"
)

type profileRepository struct {
	client *gorm.DB
}