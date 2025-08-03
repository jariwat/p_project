package repository

import (
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jariwat/p_project/profile-service/models"
	"github.com/jariwat/p_project/profile-service/service/profile"
	"gorm.io/gorm"
)

type profileRepository struct {
	client *gorm.DB
}

// FetchProfiles implements profile.ProfileRepository.
func (p *profileRepository) FetchProfiles(params profile.GetProfilesParams, paginator *models.Paginator) ([]*models.Profile, error) {
	var profiles []*models.Profile
	var totalRows int64
	var limit = paginator.PerPage
	var offset = (paginator.Page - 1) * paginator.PerPage

	query := p.client.Model(&models.Profile{})

	if params.SearchWord != nil && *params.SearchWord != "" {
		likeQuery := "%" + strings.ToLower(strings.ReplaceAll(*params.SearchWord, " ", "")) + "%"
		query = query.Where(" LOWER(REPLACE(CONCAT_WS('', first_name, middle_name, last_name), ' ', '')) LIKE ?", likeQuery)
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return nil, err
	}

	if err := query.Preload("Skills").
		Limit(limit).
		Offset(offset).
		Find(&profiles).Error; err != nil {
		return nil, err
	}

	paginator.SetTotal(int(totalRows))

	return profiles, nil
}

// FetchProfileById implements profile.ProfileRepository.
func (p *profileRepository) FetchProfileById(profileId *uuid.UUID) (*models.Profile, error) {
	var profile models.Profile
	if err := p.client.Preload("Skills").First(&profile, "id = ?", profileId).Error; err != nil {
		return nil, err
	}

	return &profile, nil
}

// CreateProfile implements profile.ProfileRepository.
func (p *profileRepository) CreateProfile(profile *models.Profile) error {
	return p.client.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(profile).Error; err != nil {
			return err
		}

		return nil
	})
}

// UpdateProfile implements profile.ProfileRepository.
func (p *profileRepository) UpdateProfile(profile *models.Profile) error {
	return p.client.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Profile{}).Where("id = ?", profile.ID).Updates(map[string]interface{}{
			"first_name":  profile.FirstName,
			"middle_name": profile.MiddleName,
			"last_name":   profile.LastName,
			"gender":      profile.Gender,
			"class":       profile.Class,
			"updated_at":  profile.UpdatedAt,
		}).Error; err != nil {
			return err
		}

		if err := tx.Where("profile_id = ?", profile.ID).Delete(&models.Skill{}).Error; err != nil {
			return err
		}

		for _, skill := range profile.Skills {
			if err := tx.Create(skill).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// DeleteProfile implements profile.ProfileRepository.
func (p *profileRepository) DeleteProfile(profileId *uuid.UUID) error {
	return p.client.Delete(&models.Profile{}, profileId).Error
}

func NewPsqlProfileRepository(client *gorm.DB) profile.ProfileRepository {
	return &profileRepository{
		client: client,
	}
}
