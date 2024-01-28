package testing

import (
	"fmt"
	"net"
)

// GetFreeTCPPort returns a free TCP port to use.
func GetFreeTCPPort() (int, error) {
	// Get a free TCP port
	addr, err := net.ResolveTCPAddr("tcp", "[::]:0")
	if err != nil {
		return 0, fmt.Errorf("%w: %s", err, "unable to resolve TCP address")
	}

	// Listen on the free TCP port
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", err, "unable to resolve TCP address")
	}

	// Get the free TCP port
	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return 0, fmt.Errorf("%s", "unable to get a free TCP port")
	}
	defer listener.Close()

	return addr.Port, nil
}
