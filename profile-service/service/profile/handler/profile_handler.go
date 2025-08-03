package handler

import (
	"github.com/gin-gonic/gin"
	_profile "github.com/jariwat/p_project/profile-service/service/profile"
)

type profileHandler struct {
	profileUs _profile.ProfileUsecase
}

// DeleteProfileId implements profile.ServerInterface.
func (p *profileHandler) DeleteProfileId(c *gin.Context, id string) {
	panic("unimplemented")
}

// GetProfileId implements profile.ServerInterface.
func (p *profileHandler) GetProfileId(c *gin.Context, id string) {
	panic("unimplemented")
}

// GetProfiles implements profile.ServerInterface.
func (p *profileHandler) GetProfiles(c *gin.Context, params _profile.GetProfilesParams) {
	panic("unimplemented")
}

// PostProfile implements profile.ServerInterface.
func (p *profileHandler) PostProfile(c *gin.Context) {
	panic("unimplemented")
}

// PutProfileId implements profile.ServerInterface.
func (p *profileHandler) PutProfileId(c *gin.Context, id string) {
	panic("unimplemented")
}

func NewProfileHandler(profileUs _profile.ProfileUsecase) _profile.ServerInterface {
	return &profileHandler{
		profileUs: profileUs,
	}
}
