package domain

import "bvpn-prototype/internal/vpn/storage"

type VpnService interface {
}

type VpnServiceImpl struct {
	profileRepo ProfileRepo
}

func (*VpnServiceImpl) Init() error {
	// todo
	return nil
}

func NewVpnService() (*VpnServiceImpl, error) {
	profileRepo, err := storage.NewProfileRepo()
	if err != nil {
		return nil, err
	}

	return &VpnServiceImpl{
		profileRepo: profileRepo,
	}, nil
}
