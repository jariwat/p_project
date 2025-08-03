package usecase

import "jariwat/p_project/profile-service/service/profile"

type profileUsecase struct {
	profileRepo profile.ProfileRepository
}

func NewProfileUsecase(profileRepo profile.ProfileRepository) profile.ProfileUsecase {
	return &profileUsecase{
		profileRepo: profileRepo,
	}
}
