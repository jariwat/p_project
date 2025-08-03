package handler

import _profile "jariwat/p_project/profile-service/service/profile"

type profileHandler struct {
	profileUs _profile.ProfileUsecase
}

func NewProfileHandler(profileUs _profile.ProfileUsecase) _profile.ServerInterface {
	return &profileHandler{
		profileUs: profileUs,
	}
}