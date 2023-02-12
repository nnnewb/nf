package pinger

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/routing"
	"github.com/nnnewb/nf/internal/protocol"
	"github.com/pkg/errors"
)

// HandShakeTCP try tcp hand shake
func HandShakeTCP(dst net.IP, port int, timeout time.Duration) error {
	c, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", dst, port), timeout)
	if err != nil {
		return errors.WithStack(err)
	}
	defer c.Close()

	return nil
}

var (
	ErrPortClosed = errors.New("port closed")
	ErrTimeout    = errors.New("timeout")
)

// SynScan SYN 半开放扫描
func SynScan(dst net.IP, port int, timeout time.Duration) error {
	var (
		src       net.IP
		gw        net.IP
		iface     *net.Interface
		dstHwAddr net.HardwareAddr
		router    routing.Router
		err       error
		handle    *pcap.Handle
	)

	// determine preferred route options
	router, err = routing.New()
	if err != nil {
		return errors.WithStack(err)
	}

	iface, gw, src, err = router.Route(dst)
	if err != nil {
		return errors.WithStack(err)
	}

	// initialize pcap
	handle, err = pcap.OpenLive(iface.Name, 0, false, timeout)
	if err != nil {
		return errors.WithStack(err)
	}

	if gw == nil {
		dstHwAddr, err = protocol.AddressResolveSys(iface, dst)
		if err != nil {
			return err
		}
	} else {
		dstHwAddr, err = protocol.AddressResolveSys(iface, gw)
		if err != nil {
			return err
		}
	}

	// construct SYN packet
	eth := layers.Ethernet{
		SrcMAC:       iface.HardwareAddr,
		DstMAC:       dstHwAddr,
		EthernetType: layers.EthernetTypeIPv4,
	}

	ipv4 := layers.IPv4{
		Version:  4,
		Protocol: layers.IPProtocolTCP,
		SrcIP:    src,
		DstIP:    dst,
		TTL:      64,
	}

	tcp := layers.TCP{
		SrcPort: 54321,
		DstPort: layers.TCPPort(port),
		SYN:     true,
		Window:  1000,
		Seq:     123456789,
	}

	err = tcp.SetNetworkLayerForChecksum(&ipv4)
	if err != nil {
		return err
	}

	result := make(chan error)
	o := sync.Once{}
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		ddl := time.After(30 * time.Second)
		for {
			select {
			case <-ddl:
				// may be packet was dropped
				result <- errors.Errorf("timeout")
				return
			default:
				// 确保开始抓包之后才发送包
				o.Do(func() { wg.Done() })

				data, _, err := handle.ZeroCopyReadPacketData()
				if err != nil {
					if errors.Is(err, pcap.NextErrorTimeoutExpired) {
						result <- ErrTimeout
						return
					}
					result <- errors.WithStack(err)
					return
				}

				packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.NoCopy)

				ipv4Layer := packet.Layer(layers.LayerTypeIPv4)
				if ipv4Layer == nil {
					continue
				}

				ipv4Packet, ok := ipv4Layer.(*layers.IPv4)
				if !ok {
					continue
				}

				tcpLayer := packet.Layer(layers.LayerTypeTCP)
				if tcpLayer == nil {
					continue
				}

				tcpPacket, ok := tcpLayer.(*layers.TCP)
				if !ok {
					continue
				}

				if ipv4Packet.SrcIP.Equal(dst) && ipv4Packet.DstIP.Equal(src) {
					if tcpPacket.DstPort == 54321 && tcpPacket.SrcPort == layers.TCPPort(port) {
						// open
						if tcpPacket.ACK && tcpPacket.SYN {
							// log.Printf("%+v", tcpPacket)
							result <- nil
							return
						}

						// close
						if tcpPacket.ACK && tcp.RST {
							result <- ErrPortClosed
							return
						}
					}
				}
			}
		}
	}()

	// 确保开始抓包之后才发送包
	wg.Wait()
	err = protocol.SendEthernetPacket(handle, &eth, &ipv4, &tcp)
	if err != nil {
		return err
	}

	return <-result
}
