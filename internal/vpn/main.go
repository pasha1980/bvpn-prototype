package vpn

import (
	"bvpn-prototype/internal/storage/vpn_profile"
	utils2 "bvpn-prototype/utils"
	"fmt"
	"github.com/songgao/water"
	"net"
)

var listener net.Listener
var conn net.Conn

/*
todo: Encryption - decryption
*/

func Init(port string, proto string) error {
	err := vpn_profile.InitStorage()
	if err != nil {
		return err // todo
	}

	iface, err := createTun()
	if err != nil {
		return err // todo
	}

	listener, err = net.Listen(proto, utils2.MyIP()+":"+port)
	if err != nil {
		return err // todo
	}

	go listenConn(iface)
	go listenInterface(iface)

	return nil
}

func createTun() (*water.Interface, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}

	iface, err := water.New(config)
	if err != nil {
		return nil, err
	}
	_, err = utils2.Exec(fmt.Sprintf("sudo ip addr add 0.0.0.0/0 dev %s", iface.Name()))
	if err != nil {
		return nil, err
	}

	_, err = utils2.Exec(fmt.Sprintf("sudo ip link set dev %s up", iface.Name()))
	if err != nil {
		return nil, err
	}
	return iface, nil
}

func listenConn(iface *water.Interface) {
	var err error
	conn, err = listener.Accept()
	if err != nil {
		// todo
	}

	for {
		message := make([]byte, 65535)
		for {
			n, err := conn.Read(message)
			if err != nil {
				// todo
			}
			if iface != nil {
				_, err = iface.Write(message[:n])
				if err != nil {
					// todo
				} else {
					// todo
				}
			}
		}

		// todo: count traffic

	}
}

func listenInterface(iface *water.Interface) {
	packet := make([]byte, 65535)
	for {
		n, err := iface.Read(packet)
		if err != nil {
			// todo
		}

		if conn != nil {
			_, err = conn.Write(packet[:n])
			if err != nil {
				// todo
			}
		}

		// todo: count traffic
	}
}
