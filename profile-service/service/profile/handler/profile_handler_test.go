package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jariwat/p_project/profile-service/constants"
	"github.com/jariwat/p_project/profile-service/models"
	_profile "github.com/jariwat/p_project/profile-service/service/profile"
	"github.com/jariwat/p_project/profile-service/service/profile/mocks"
	"github.com/oapi-codegen/runtime/types"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func ptrUUID() *uuid.UUID {
	id, _ := uuid.NewV4()
	return &id
}

func TestDeleteProfileId_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange
	mockUsecase := new(mocks.ProfileUsecase)

	profileID := ptrUUID()
	mockUsecase.
		On("DeleteProfile", mock.AnythingOfType("*uuid.UUID")).
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/profile/"+profileID.String(), nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: profileID.String()}}
	c.Request = req

	// Act
	handler := NewProfileHandler(mockUsecase)
	handler.DeleteProfileId(c, (types.UUID)(*profileID))

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestDeleteProfileId_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := new(mocks.ProfileUsecase)

	profileID := ptrUUID()
	mockUsecase.
		On("DeleteProfile", mock.AnythingOfType("*uuid.UUID")).
		Return(errors.New("delete error"))

	req := httptest.NewRequest(http.MethodDelete, "/profile/"+profileID.String(), nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: profileID.String()}}
	c.Request = req

	handler := NewProfileHandler(mockUsecase)
	handler.DeleteProfileId(c, (types.UUID)(*profileID))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestGetProfileId_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := new(mocks.ProfileUsecase)

	profileID := ptrUUID()
	expectedProfile := &models.Profile{
		ID:        profileID,
		FirstName: "SeiA",
		LastName:  "Phanes",
	}

	mockUsecase.
		On("FetchProfileById", mock.MatchedBy(func(id *uuid.UUID) bool {
			return id != nil && *id == *profileID
		})).
		Return(expectedProfile, nil)

	// Prepare HTTP request
	req := httptest.NewRequest(http.MethodGet, "/profile/"+profileID.String(), nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := NewProfileHandler(mockUsecase)
	handler.GetProfileId(c, (types.UUID)(*profileID))

	assert.Equal(t, http.StatusOK, w.Code)

	var response _profile.ProfileResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)
	assert.Equal(t, expectedProfile.FirstName, *response.Data.FirstName)
}

func TestGetProfileId_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := new(mocks.ProfileUsecase)

	profileID := ptrUUID()

	mockUsecase.
		On("FetchProfileById", mock.AnythingOfType("*uuid.UUID")).
		Return(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/profile/"+profileID.String(), nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := NewProfileHandler(mockUsecase)
	handler.GetProfileId(c, (types.UUID)(*profileID))

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Profile not found")
}

func TestGetProfileId_FetchError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := new(mocks.ProfileUsecase)

	profileID := ptrUUID()

	mockUsecase.
		On("FetchProfileById", mock.AnythingOfType("*uuid.UUID")).
		Return(nil, errors.New("fetch error"))

	req := httptest.NewRequest(http.MethodGet, "/profile/"+profileID.String(), nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := NewProfileHandler(mockUsecase)
	handler.GetProfileId(c, (types.UUID)(*profileID))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "fetch error")
}

func TestGetProfiles_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := new(mocks.ProfileUsecase)

	page := 1
	perPage := 10
	params := _profile.GetProfilesParams{
		Page:    &page,
		PerPage: &perPage,
	}

	mockProfiles := []*models.Profile{
		{
			ID:        ptrUUID(),
			FirstName: "SeiA",
			LastName:  "Phanes",
		},
	}

	mockUsecase.
		On("FetchProfiles", mock.Anything, mock.AnythingOfType("*models.Paginator")).
		Return(mockProfiles, nil)

	req := httptest.NewRequest(http.MethodGet, "/profiles?page=1&per_page=10", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := NewProfileHandler(mockUsecase)
	handler.GetProfiles(c, params)

	assert.Equal(t, http.StatusOK, w.Code)

	var response _profile.ProfilesPaginationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)
}

func TestGetProfiles_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := new(mocks.ProfileUsecase)

	page := 1
	perPage := 10
	params := _profile.GetProfilesParams{
		Page:    &page,
		PerPage: &perPage,
	}

	mockUsecase.
		On("FetchProfiles", mock.Anything, mock.AnythingOfType("*models.Paginator")).
		Return([]*models.Profile{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/profiles?page=1&per_page=10", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := NewProfileHandler(mockUsecase)
	handler.GetProfiles(c, params)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetProfiles_FetchError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := new(mocks.ProfileUsecase)

	page := 1
	perPage := 10
	params := _profile.GetProfilesParams{
		Page:    &page,
		PerPage: &perPage,
	}

	mockUsecase.
		On("FetchProfiles", mock.Anything, mock.AnythingOfType("*models.Paginator")).
		Return(nil, errors.New("fetch error"))

	req := httptest.NewRequest(http.MethodGet, "/profiles?page=1&per_page=10", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := NewProfileHandler(mockUsecase)
	handler.GetProfiles(c, params)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPostProfile_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock request payload
	newProfile := _profile.UpsertProfile{
		FirstName: "ทดสอบ",
		LastName:  "ทดสอบ",
		Class:     "ทดสอบ",
		Gender:    "MALE",
		Skills: []_profile.UpsertSkill{
			{
				Skill:  "ทดสอบ",
				Detail: "ทดสอบ",
			},
		},
	}
	body, _ := json.Marshal(newProfile)

	req, err := http.NewRequest(http.MethodPost, "/profile", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Setup mock and handler
	mockUsecase := new(mocks.ProfileUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Prepare profile (with generated UUID)
	expectedProfile := &models.Profile{}
	expectedProfile.GenUUID()

	// Mock CreateProfile call
	mockUsecase.
		On("CreateProfile", mock.AnythingOfType("*models.Profile"), newProfile).
		Run(func(args mock.Arguments) {
			// Copy generated UUID into the argument profile
			argProfile := args.Get(0).(*models.Profile)
			*argProfile = *expectedProfile
		}).
		Return(nil)

	// Call handler
	handler := NewProfileHandler(mockUsecase)
	handler.PostProfile(c)

	// Assertions
	require.Equal(t, http.StatusOK, w.Code)

	var resp _profile.Success
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "Profile created successfully", resp.Message)
	assert.Equal(t, expectedProfile.ID.String(), resp.Id.String())

	mockUsecase.AssertExpectations(t)
}

func TestPostProfile_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{invalid-json}`)

	req, _ := http.NewRequest(http.MethodPost, "/profile", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Setup mock and handler
	mockUsecase := new(mocks.ProfileUsecase)
	handler := NewProfileHandler(mockUsecase)
	handler.PostProfile(c)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Invalid input", resp["error"])
}

func TestPostProfile_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock request payload
	newProfile := _profile.UpsertProfile{
		FirstName: "ทดสอบ",
		LastName:  "ทดสอบ",
		Class:     "ทดสอบ",
		Gender:    "MALE",
		Skills: []_profile.UpsertSkill{
			{
				Skill:  "ทดสอบ",
				Detail: "ทดสอบ",
			},
		},
	}
	body, _ := json.Marshal(newProfile)

	req, _ := http.NewRequest(http.MethodPost, "/profile", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	mockUsecase := new(mocks.ProfileUsecase)

	mockUsecase.
		On("CreateProfile", mock.AnythingOfType("*models.Profile"), newProfile).
		Return(errors.New("create error"))

	handler := NewProfileHandler(mockUsecase)
	handler.PostProfile(c)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "create error", resp["error"])
}

func TestPutProfileId_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	profileId := ptrUUID()
	updateProfile := _profile.UpsertProfile{
		FirstName: "ทดสอบ",
		LastName:  "ทดสอบ",
		Class:     "ทดสอบ",
		Gender:    "MALE",
		Skills: []_profile.UpsertSkill{
			{
				Skill:  "ทดสอบ",
				Detail: "ทดสอบ",
			},
		},
	}
	body, _ := json.Marshal(updateProfile)

	req, _ := http.NewRequest(http.MethodPut, "/profile/"+profileId.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: profileId.String()}}

	mockUsecase := new(mocks.ProfileUsecase)

	mockUsecase.On("UpdateProfile", mock.MatchedBy(func(pID *uuid.UUID) bool {
		return *pID == *profileId
	}), updateProfile).Return(nil)

	handler := NewProfileHandler(mockUsecase)
	handler.PutProfileId(c, (types.UUID)(*profileId))

	require.Equal(t, http.StatusOK, w.Code)

	var resp _profile.Success
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "Profile updated successfully", resp.Message)
	assert.Equal(t, profileId.String(), resp.Id.String())

	mockUsecase.AssertExpectations(t)
}

func TestPutProfileId_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)

	profileId := ptrUUID()
	body := []byte(`{invalid-json}`)

	req, _ := http.NewRequest(http.MethodPut, "/profile/"+profileId.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	mockUsecase := new(mocks.ProfileUsecase)
	handler := NewProfileHandler(mockUsecase)
	handler.PutProfileId(c, (types.UUID)(*profileId))

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Invalid input", resp["error"])
}

func TestPutProfileId_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	profileId := ptrUUID()
	updateProfile := _profile.UpsertProfile{
		FirstName: "ทดสอบ",
		LastName:  "ทดสอบ",
		Class:     "ทดสอบ",
		Gender:    "MALE",
		Skills: []_profile.UpsertSkill{
			{
				Skill:  "ทดสอบ",
				Detail: "ทดสอบ",
			},
		},
	}
	body, _ := json.Marshal(updateProfile)

	req, _ := http.NewRequest(http.MethodPut, "/profile/"+profileId.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: profileId.String()}}

	mockUsecase := new(mocks.ProfileUsecase)

	mockUsecase.
		On("UpdateProfile", mock.AnythingOfType("*uuid.UUID"), updateProfile).
		Return(constants.ErrProfileNotFound)

	handler := NewProfileHandler(mockUsecase)
	handler.PutProfileId(c, (types.UUID)(*profileId))

	require.Equal(t, http.StatusConflict, w.Code)

	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, constants.ErrProfileNotFound.Error(), resp["error"])
}

func TestPutProfileId_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	profileId := ptrUUID()
	updateProfile := _profile.UpsertProfile{
		FirstName: "ทดสอบ",
		LastName:  "ทดสอบ",
		Class:     "ทดสอบ",
		Gender:    "MALE",
		Skills: []_profile.UpsertSkill{
			{
				Skill:  "ทดสอบ",
				Detail: "ทดสอบ",
			},
		},
	}
	body, _ := json.Marshal(updateProfile)

	req, _ := http.NewRequest(http.MethodPut, "/profile/"+profileId.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: profileId.String()}}

	mockUsecase := new(mocks.ProfileUsecase)

	mockUsecase.
		On("UpdateProfile", mock.AnythingOfType("*uuid.UUID"), updateProfile).
		Return(errors.New("unexpected DB error"))

	handler := NewProfileHandler(mockUsecase)
	handler.PutProfileId(c, (types.UUID)(*profileId))

	require.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "unexpected DB error", resp["error"])
}
