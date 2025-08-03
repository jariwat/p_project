package profile

import (
	"github.com/gofrs/uuid"
	"github.com/jariwat/p_project/profile-service/models"
)

type ProfileUsecase interface {
	FetchProfiles(params GetProfilesParams, paginator *models.Paginator) ([]*models.Profile, error)
	FetchProfileById(profileId *uuid.UUID) (*models.Profile, error)
	CreateProfile(profile *models.Profile, newProfile UpsertProfile) error
	UpdateProfile(profileId *uuid.UUID, updateProfile UpsertProfile) error
	DeleteProfile(profileId *uuid.UUID) error
}
