package protocol

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"
)

func readLine(reader *bufio.Reader) ([]byte, error) {
	line, isPrefix, err := reader.ReadLine()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return []byte{}, nil
		}
		return nil, err
	}

	for isPrefix {
		var b []byte
		b, isPrefix, err = reader.ReadLine()
		if err != nil {
			return nil, err
		}
		line = append(line, b...)
	}

	return line, nil
}

// EthernetAddressResolver ARP 实现
type EthernetAddressResolver struct {
	cache map[string]net.HardwareAddr
}

func NewEthernetAddressResolver() EthernetAddressResolver {
	cacheMap := make(map[string]net.HardwareAddr)
	file, err := os.Open("/proc/net/arp")
	if err == nil {
		defer file.Close()
		reader := bufio.NewReader(file)
		_, _ = readLine(reader)
		for line, err := readLine(reader); err == nil && len(line) != 0; line, err = readLine(reader) {
			fields := strings.Fields(string(line))
			hardwareAddr, err := net.ParseMAC(string(fields[3]))
			if err != nil {
				continue
			}

			cacheMap[string(fields[0])] = hardwareAddr
		}
	}

	return EthernetAddressResolver{
		cache: cacheMap,
	}
}

// sendARPRequest 发送 ARP 请求
func sendARPRequest(handle *pcap.Handle, iface *net.Interface, src, dst net.IP) error {
	// 以太网协议需要对端的 MAC 地址，如果需要经过路由，则需要路由的 MAC 地址。
	// 先通过 ARP 协议广播获取对端的 MAC 地址。
	eth := layers.Ethernet{
		SrcMAC:       iface.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, // broadcast MAC address
		EthernetType: layers.EthernetTypeARP,
	}

	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(iface.HardwareAddr),
		SourceProtAddress: []byte(src.To4()),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
		DstProtAddress:    []byte(dst.To4()),
	}

	return SendEthernetPacket(handle, &eth, &arp)
}

// recvARPReply 接收 ARP 响应
func recvARPReply(handle *pcap.Handle, src, dst net.IP) (net.HardwareAddr, error) {
	// wait 3 seconds for arp reply
	ddl := time.After(3 * time.Second)

	for {
		select {
		case <-ddl:
			return nil, errors.Errorf("timeout getting ARP reply")
		default:
			data, _, err := handle.ZeroCopyReadPacketData()
			if err != nil {
				return nil, errors.WithStack(err)
			}

			packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.NoCopy)
			if arpLayer := packet.Layer(layers.LayerTypeARP); arpLayer != nil {
				arp := arpLayer.(*layers.ARP)
				if net.IP(arp.SourceProtAddress).Equal(dst) {
					log.Printf("arp reply: IP %s ether %s", net.IP(arp.SourceProtAddress), net.HardwareAddr(arp.SourceHwAddress))
					return arp.SourceHwAddress, nil
				}
			}
		}
	}
}

// AddressResolve 指定网卡获取目标IP的MAC地址
func AddressResolve(handle *pcap.Handle, iface *net.Interface, src, dst net.IP) (net.HardwareAddr, error) {
	err := sendARPRequest(handle, iface, src, dst)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return recvARPReply(handle, src, dst)
}

// resolveHardwareAddressSys 用 arp 命令解析地址
func resolveHardwareAddressSys(iface *net.Interface, dst net.IP) error {
	cmd := exec.Command("arp", "-i", iface.Name, dst.String())
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func readArpCache(iface *net.Interface, dst net.IP) (net.HardwareAddr, error) {
	file, err := os.Open("/proc/net/arp")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var isPrefix = true

	for isPrefix {
		_, isPrefix, err = reader.ReadLine()
		if err != nil {
			return nil, err
		}
	}

	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, errors.Errorf("HwAddr of IP %s not found in /proc/net/arp", dst)
			}
			return nil, err
		}

		for isPrefix {
			var b []byte
			b, isPrefix, err = reader.ReadLine()
			if err != nil {
				return nil, err
			}
			line = append(line, b...)
		}

		fields := strings.Fields(string(line))
		if fields[5] == iface.Name && fields[0] == dst.String() {
			mac, err := net.ParseMAC(fields[3])
			if err != nil {
				return nil, err
			}
			return mac, nil
		}
	}
}

func AddressResolveSys(iface *net.Interface, dst net.IP) (net.HardwareAddr, error) {
	hwAddr, err := readArpCache(iface, dst)
	if err == nil {
		return hwAddr, nil
	}

	err = resolveHardwareAddressSys(iface, dst)
	if err != nil {
		return nil, err
	}

	hwAddr, err = readArpCache(iface, dst)
	if err != nil {
		return nil, err
	}

	return hwAddr, nil
}
