package domain

import (
	"bvpn-prototype/internal/infrastructure/config"
	"bvpn-prototype/internal/infrastructure/di"
	common_errors "bvpn-prototype/internal/infrastructure/errors"
	"bvpn-prototype/internal/protocol"
	"bvpn-prototype/internal/protocol/entity/block_data"
	"bvpn-prototype/internal/vpn/errors"
	"bvpn-prototype/internal/vpn/storage"
	"bvpn-prototype/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/songgao/water"
	"net"
	"time"
)

type VpnService interface {
	CreateConnection(clientAddr string) (*PublicProfile, error)
	BreakConnection(id uuid.UUID) error
}

type VpnServiceImpl struct {
	profileRepo ProfileRepo

	conn     net.Conn
	listener net.Listener
	connMap  map[string]Connection
}

func (s *VpnServiceImpl) Init() error {
	iface, err := s.createTun()
	if err != nil {
		return err
	}

	s.listener, err = net.Listen(config.Get().VpnProto, utils.MyIP()+":"+config.Get().VpnPort)
	if err != nil {
		return errors.NetworkConfigurationError(err.Error())
	}

	errCh := make(chan error, 5)
	go s.listenConn(iface, errCh)
	go s.listenInterface(iface, errCh)
	go s.listenErrors(errCh)

	return nil
}

func (s *VpnServiceImpl) CreateConnection(clientAddr string) (*PublicProfile, error) {
	chainService := di.Get("chain_public").(ChainPublicService)
	offer, err := chainService.GetMyLastOffer()
	if err != nil {
		return nil, err
	}

	if offer == nil {
		return nil, errors.NoOfferError()
	}

	profile := protocol.GenerateVpnProfile(*offer, clientAddr)
	_, err = s.profileRepo.Save(profile)
	if err != nil {
		return nil, common_errors.StorageError(err.Error())
	}

	s.connMap[profile.Id.String()] = Connection{
		Profile: profile,
		Traffic: 0,
	}

	pub := ProfileToPub(profile)
	return &pub, nil
}

func (s *VpnServiceImpl) BreakConnection(id uuid.UUID) error {
	ok, err := s.profileRepo.IsExist(id)
	if err != nil {
		return common_errors.StorageError(err.Error())
	}

	if ok {
		err = s.profileRepo.Remove(id)
		if err != nil {
			return common_errors.StorageError(err.Error())
		}
	}

	connection, ok := s.connMap[id.String()]
	if !ok {
		return nil // todo
	}

	chainService := di.Get("chain_public").(ChainPublicService)
	result := block_data.Traffic{
		Timestamp: time.Now(),
		Node:      protocol.GetMyAddr(),
		Client:    connection.Profile.Client,
		Bytes:     connection.Traffic,
	}
	err = chainService.SaveTraffic(result)
	if err != nil {
		return err
	}

	return nil
}

func (s *VpnServiceImpl) createTun() (*water.Interface, error) {
	cnf := water.Config{
		DeviceType: water.TUN,
	}

	iface, err := water.New(cnf)
	if err != nil {
		return nil, errors.NetworkConfigurationError(err.Error())
	}
	_, err = utils.Exec(fmt.Sprintf("sudo ip addr add 0.0.0.0/0 dev %s", iface.Name()))
	if err != nil {
		return nil, errors.NetworkConfigurationError(err.Error())
	}

	_, err = utils.Exec(fmt.Sprintf("sudo ip link set dev %s up", iface.Name()))
	if err != nil {
		return nil, errors.NetworkConfigurationError(err.Error())
	}

	return iface, nil
}

func (s *VpnServiceImpl) listenConn(iface *water.Interface, errCh chan error) {
	var err error
	s.conn, err = s.listener.Accept()
	if err != nil {
		errCh <- errors.ConnectionListenerError(err.Error())
		return
	}

	message := make([]byte, 65535)
	for {
		n, err := s.conn.Read(message)
		if err != nil {
			errCh <- errors.ConnectionListenerError(err.Error())
			return
		}

		// todo: counting traffic
		// todo: encryption

		_, err = iface.Write(message[:n])
		if err != nil {
			errCh <- errors.ConnectionListenerError(err.Error())
			return
		}

	}
}

func (s *VpnServiceImpl) listenInterface(iface *water.Interface, errCh chan error) {
	packet := make([]byte, 65535)
	for {
		n, err := iface.Read(packet)
		if err != nil {
			errCh <- errors.InterfaceListenerError(err.Error())
			continue
		}

		// todo: decryption

		_, err = s.conn.Write(packet[:n])
		if err != nil {
			errCh <- errors.InterfaceListenerError(err.Error())
			continue
		}
	}
}

func (s *VpnServiceImpl) listenErrors(errCh chan error) {
	// todo
}

func NewVpnService() (*VpnServiceImpl, error) {
	profileRepo, err := storage.NewProfileRepo()
	if err != nil {
		return nil, err
	}

	return &VpnServiceImpl{
		profileRepo: profileRepo,
		connMap:     make(map[string]Connection),
	}, nil
}
