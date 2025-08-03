package profile

import (
	"github.com/gofrs/uuid"
	"github.com/jariwat/p_project/profile-service/models"
)

type ProfileRepository interface {
	FetchProfiles(params GetProfilesParams, paginator *models.Paginator) ([]*models.Profile, error)
	FetchProfileById(profileId *uuid.UUID) (*models.Profile, error)
	CreateProfile(profile *models.Profile) error
	UpdateProfile(profile *models.Profile) error
	DeleteProfile(profileId *uuid.UUID) error
}
