package usecase

import (
	"errors"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jariwat/p_project/profile-service/constants"
	"github.com/jariwat/p_project/profile-service/models"
	_profile "github.com/jariwat/p_project/profile-service/service/profile"
	"github.com/jariwat/p_project/profile-service/service/profile/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func ptrUUID() *uuid.UUID {
	id, _ := uuid.NewV4()
	return &id
}

func TestFetchProfiles_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.ProfileRepository)
	usecase := NewProfileUsecase(mockRepo)

	var page = 1
	var perPage = 10

	params := _profile.GetProfilesParams{
		SearchWord: nil,
		Page:       &page,
		PerPage:    &perPage,
	}
	paginator := &models.Paginator{Page: 1, PerPage: 10}

	expectedProfiles := []*models.Profile{
		{FirstName: "SeiA", LastName: "Phanes"},
		{FirstName: "AliZe", LastName: "Phanes"},
	}

	// Act
	mockRepo.
		On("FetchProfiles", params, paginator).
		Return(expectedProfiles, nil)

	result, err := usecase.FetchProfiles(params, paginator)

	// Assert
	require.NoError(t, err)
	require.Equal(t, expectedProfiles, result)

	mockRepo.AssertExpectations(t)
}

func TestFetchProfiles_Error(t *testing.T) {
	mockRepo := new(mocks.ProfileRepository)
	usecase := NewProfileUsecase(mockRepo)

	params := _profile.GetProfilesParams{}
	paginator := &models.Paginator{Page: 1, PerPage: 10}

	expectedErr := errors.New("db error")

	mockRepo.
		On("FetchProfiles", params, paginator).
		Return(nil, expectedErr)

	result, err := usecase.FetchProfiles(params, paginator)

	require.Error(t, err)
	require.Equal(t, expectedErr, err)
	require.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestFetchProfileById_Success(t *testing.T) {
	mockRepo := new(mocks.ProfileRepository)
	usecase := NewProfileUsecase(mockRepo)

	profileID := ptrUUID()
	expected := &models.Profile{
		ID:        profileID,
		FirstName: "John",
		LastName:  "Doe",
	}

	mockRepo.
		On("FetchProfileById", profileID).
		Return(expected, nil)

	result, err := usecase.FetchProfileById(profileID)

	require.NoError(t, err)
	require.Equal(t, expected, result)

	mockRepo.AssertExpectations(t)
}

func TestFetchProfileById_Error(t *testing.T) {
	mockRepo := new(mocks.ProfileRepository)
	usecase := NewProfileUsecase(mockRepo)

	profileID := ptrUUID()
	expectedErr := errors.New("not found")

	mockRepo.
		On("FetchProfileById", profileID).
		Return(nil, expectedErr)

	result, err := usecase.FetchProfileById(profileID)

	require.Error(t, err)
	require.Equal(t, expectedErr, err)
	require.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestCreateProfile_Success(t *testing.T) {
	// Mock repository
	mockRepo := new(mocks.ProfileRepository)
	usecase := NewProfileUsecase(mockRepo)

	// Prepare input
	profile := &models.Profile{}

	middle := "T"
	newProfile := _profile.UpsertProfile{
		FirstName:  "SeiA",
		MiddleName: &middle,
		LastName:   "Phanes",
		Gender:     "MALE",
		Class:      "King",
		Skills: []_profile.UpsertSkill{
			{
				Skill:  "Swordsmanship",
				Detail: "Expert in sword fighting techniques",
			},
		},
	}

	// Expect CreateProfile to be called with a filled Profile
	mockRepo.
		On("CreateProfile", mock.MatchedBy(func(p *models.Profile) bool {
			return p.FirstName == "SeiA" &&
				p.MiddleName == &middle &&
				p.LastName == "Phanes" &&
				p.Gender == models.Gender("MALE") &&
				p.Class == "King" &&
				len(p.Skills) == 1 &&
				p.Skills[0].Skill == "Swordsmanship" &&
				p.Skills[0].Detail == "Expert in sword fighting techniques"
		})).
		Return(nil)

	// Act
	err := usecase.CreateProfile(profile, newProfile)

	// Assert
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateProfile_RepoError(t *testing.T) {
	mockRepo := new(mocks.ProfileRepository)
	usecase := NewProfileUsecase(mockRepo)

	profile := &models.Profile{}
	newProfile := _profile.UpsertProfile{
		FirstName: "SeiA",
		LastName:  "Phanes",
		Gender:    "MALE",
		Class:     "King",
	}

	mockRepo.
		On("CreateProfile", mock.Anything).
		Return(errors.New("db error"))

	err := usecase.CreateProfile(profile, newProfile)

	require.Error(t, err)
	require.EqualError(t, err, "db error")
	mockRepo.AssertExpectations(t)
}

func TestUpdateProfile_Success(t *testing.T) {
	mockRepo := new(mocks.ProfileRepository)
	usecase := NewProfileUsecase(mockRepo)

	profileID := ptrUUID()
	middle := "F"
	update := _profile.UpsertProfile{
		FirstName:  "SeiA",
		MiddleName: &middle,
		LastName:   "Phanes",
		Gender:     "MALE",
		Class:      "Yuusha",
		Skills: []_profile.UpsertSkill{
			{Skill: "Swordsmanship", Detail: "Strong in sword fighting techniques"},
		},
	}

	existingProfile := &models.Profile{
		ID: profileID,
	}

	mockRepo.On("FetchProfileById", profileID).Return(existingProfile, nil)
	mockRepo.On("UpdateProfile", mock.MatchedBy(func(p *models.Profile) bool {
		return p.FirstName == "SeiA" && p.LastName == "Phanes" && p.Gender == models.Gender("MALE") && len(p.Skills) == 1
	})).Return(nil)

	err := usecase.UpdateProfile(profileID, update)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProfile_ProfileNotFound(t *testing.T) {
	mockRepo := new(mocks.ProfileRepository)
	usecase := NewProfileUsecase(mockRepo)

	profileID := ptrUUID()
	mockRepo.On("FetchProfileById", profileID).Return(nil, nil)

	err := usecase.UpdateProfile(profileID, _profile.UpsertProfile{})

	require.Error(t, err)
	require.Equal(t, constants.ErrProfileNotFound, err)
}

func TestUpdateProfile_FetchError(t *testing.T) {
	mockRepo := new(mocks.ProfileRepository)
	usecase := NewProfileUsecase(mockRepo)

	profileID := ptrUUID()
	mockRepo.On("FetchProfileById", profileID).Return(nil, errors.New("db error"))

	err := usecase.UpdateProfile(profileID, _profile.UpsertProfile{})

	require.EqualError(t, err, "db error")
}

func TestUpdateProfile_UpdateError(t *testing.T) {
	mockRepo := new(mocks.ProfileRepository)
	usecase := NewProfileUsecase(mockRepo)

	profileID := ptrUUID()
	existingProfile := &models.Profile{ID: profileID}
	mockRepo.On("FetchProfileById", profileID).Return(existingProfile, nil)
	mockRepo.On("UpdateProfile", mock.Anything).Return(errors.New("update failed"))

	err := usecase.UpdateProfile(profileID, _profile.UpsertProfile{FirstName: "Test"})

	require.EqualError(t, err, "update failed")
}

func TestDeleteProfile_Success(t *testing.T) {
	mockRepo := new(mocks.ProfileRepository)
	usecase := NewProfileUsecase(mockRepo)

	profileID := ptrUUID()

	// Setup expectation
	mockRepo.On("DeleteProfile", profileID).Return(nil)

	// Act
	err := usecase.DeleteProfile(profileID)

	// Assert
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProfile_Error(t *testing.T) {
	mockRepo := new(mocks.ProfileRepository)
	usecase := NewProfileUsecase(mockRepo)

	profileID := ptrUUID()
	mockRepo.On("DeleteProfile", profileID).Return(errors.New("delete failed"))

	err := usecase.DeleteProfile(profileID)

	require.EqualError(t, err, "delete failed")
	mockRepo.AssertExpectations(t)
}