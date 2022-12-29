package pinger

import (
	"net"
	"net/netip"

	"github.com/pkg/errors"
)

// SendUDPPacket send udp msg
func SendUDPPacket(dst net.IP, port int) error {
	addr, _ := netip.AddrFromSlice(dst)
	addrPort := netip.AddrPortFrom(addr, uint16(port))
	c, err := net.DialUDP("udp4", nil, net.UDPAddrFromAddrPort(addrPort))
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = c.Write([]byte("PING"))
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
