package usecase

import (
	"github.com/gofrs/uuid"
	"github.com/jariwat/p_project/profile-service/constants"
	"github.com/jariwat/p_project/profile-service/models"
	"github.com/jariwat/p_project/profile-service/service/profile"
)

type profileUsecase struct {
	profileRepo profile.ProfileRepository
}

// FetchProfiles implements profile.ProfileUsecase.
func (p *profileUsecase) FetchProfiles(params profile.GetProfilesParams, paginator *models.Paginator) ([]*models.Profile, error) {
	return p.profileRepo.FetchProfiles(params, paginator)
}

// FetchProfileById implements profile.ProfileUsecase.
func (p *profileUsecase) FetchProfileById(profileId *uuid.UUID) (*models.Profile, error) {
	return p.profileRepo.FetchProfileById(profileId)
}

// CreateProfile implements profile.ProfileUsecase.
func (p *profileUsecase) CreateProfile(profile *models.Profile, newProfile profile.UpsertProfile) error {
	profile.FirstName = newProfile.FirstName
	profile.MiddleName = newProfile.MiddleName
	profile.LastName = newProfile.LastName
	profile.Gender = models.Gender(newProfile.Gender)
	profile.Class = newProfile.Class
	profile.SetCreatedAt()
	profile.SetUpdatedAt()
	if newProfile.Skills != nil && len(newProfile.Skills) > 0 {
		profile.Skills = make([]*models.Skill, 0)
		for _, skill := range newProfile.Skills {
			skills := &models.Skill{
				ProfileID: profile.ID,
				Skill:     skill.Skill,
				Detail:    skill.Detail,
			}
			skills.GenUUID()
			skills.SetCreatedAt()
			skills.SetUpdatedAt()

			profile.Skills = append(profile.Skills, skills)
		}
	}

	return p.profileRepo.CreateProfile(profile)
}

// UpdateProfile implements profile.ProfileUsecase.
func (p *profileUsecase) UpdateProfile(profileId *uuid.UUID, updateProfile profile.UpsertProfile) error {
	profile, err := p.profileRepo.FetchProfileById(profileId)
	if err != nil {
		return err
	}

	if profile == nil {
		return constants.ErrProfileNotFound
	}

	profile.FirstName = updateProfile.FirstName
	profile.MiddleName = updateProfile.MiddleName
	profile.LastName = updateProfile.LastName
	profile.Gender = models.Gender(updateProfile.Gender)
	profile.Class = updateProfile.Class
	profile.SetUpdatedAt()
	if updateProfile.Skills != nil && len(updateProfile.Skills) > 0 {
		profile.Skills = make([]*models.Skill, 0)
		for _, skill := range updateProfile.Skills {
			skills := &models.Skill{
				ProfileID: profile.ID,
				Skill:     skill.Skill,
				Detail:    skill.Detail,
			}
			skills.GenUUID()
			skills.SetCreatedAt()
			skills.SetUpdatedAt()

			profile.Skills = append(profile.Skills, skills)
		}
	}

	return p.profileRepo.UpdateProfile(profile)
}

// DeleteProfile implements profile.ProfileUsecase.
func (p *profileUsecase) DeleteProfile(profileId *uuid.UUID) error {
	return p.profileRepo.DeleteProfile(profileId)
}

func NewProfileUsecase(profileRepo profile.ProfileRepository) profile.ProfileUsecase {
	return &profileUsecase{
		profileRepo: profileRepo,
	}
}
