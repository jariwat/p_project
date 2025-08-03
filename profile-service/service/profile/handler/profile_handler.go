package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jariwat/p_project/profile-service/constants"
	"github.com/jariwat/p_project/profile-service/models"
	_profile "github.com/jariwat/p_project/profile-service/service/profile"
	"github.com/oapi-codegen/runtime/types"
)

type profileHandler struct {
	profileUs _profile.ProfileUsecase
}

// DeleteProfileId implements profile.ServerInterface.
func (p *profileHandler) DeleteProfileId(c *gin.Context, id types.UUID) {
	var profileId = uuid.FromStringOrNil(id.String())

	if err := p.profileUs.DeleteProfile(&profileId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := _profile.Success{
		Message: "Profile deleted successfully",
	}

	c.JSON(http.StatusOK, response)
}

// GetProfileId implements profile.ServerInterface.
func (p *profileHandler) GetProfileId(c *gin.Context, id types.UUID) {
	var profileId = uuid.FromStringOrNil(id.String())

	profile, err := p.profileUs.FetchProfileById(&profileId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	var data _profile.Profile
	bu, err := json.Marshal(profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal profile"})
		return
	}

	if err := json.Unmarshal(bu, &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal profile"})
		return
	}

	response := _profile.ProfileResponse{
		Data: &data,
	}

	c.JSON(http.StatusOK, response)
}

// GetProfiles implements profile.ServerInterface.
func (p *profileHandler) GetProfiles(c *gin.Context, params _profile.GetProfilesParams) {
	var page, perPage int
	if params.Page != nil && params.PerPage != nil {
		page = *params.Page
		perPage = *params.PerPage
	}
	var paginator = models.NewPaginator(page, perPage)

	profiles, err := p.profileUs.FetchProfiles(params, paginator)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if profiles == nil || len(profiles) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No profiles found"})
		return
	}

	var data []_profile.Profiles
	bu, err := json.Marshal(profiles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal profiles"})
		return
	}

	if err := json.Unmarshal(bu, &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal profiles"})
		return
	}

	response := _profile.ProfilesPaginationResponse{
		Data:       &data,
		Page:       &paginator.Page,
		PerPage:    &paginator.PerPage,
		TotalPages: &paginator.TotalPages,
		TotalRows:  &paginator.TotalRows,
	}

	c.JSON(http.StatusOK, response)
}

// PostProfile implements profile.ServerInterface.
func (p *profileHandler) PostProfile(c *gin.Context) {
	var newProfile _profile.UpsertProfile
	if err := c.ShouldBindJSON(&newProfile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var profile = new(models.Profile)
	profile.GenUUID()
	if err := p.profileUs.CreateProfile(profile, newProfile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := _profile.Success{
		Message: "Profile created successfully",
		Id:      (*types.UUID)(profile.ID),
	}

	c.JSON(http.StatusOK, response)
}

// PutProfileId implements profile.ServerInterface.
func (p *profileHandler) PutProfileId(c *gin.Context, id types.UUID) {
	var profileId = uuid.FromStringOrNil(id.String())

	var updateProfile _profile.UpsertProfile
	if err := c.ShouldBindJSON(&updateProfile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := p.profileUs.UpdateProfile(&profileId, updateProfile); err != nil {
		if errors.Is(err, constants.ErrProfileNotFound) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := _profile.Success{
		Message: "Profile updated successfully",
		Id:      (*types.UUID)(&profileId),
	}

	c.JSON(http.StatusOK, response)
}

func NewProfileHandler(profileUs _profile.ProfileUsecase) _profile.ServerInterface {
	return &profileHandler{
		profileUs: profileUs,
	}
}
