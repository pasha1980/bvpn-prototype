package dto

import "bvpn-prototype/internal/vpn/domain"

type ProfileDto struct {
	ID    string `json:"id"`
	Proto string `json:"proto"`
	Port  string `json:"port"`
	Key   string `json:"key"`
}

func PublicProfileToDTO(profile domain.PublicProfile) ProfileDto {
	return ProfileDto{
		ID:    profile.Id.String(),
		Proto: profile.Proto,
		Port:  profile.Port,
		Key:   profile.PubKey,
	}
}
