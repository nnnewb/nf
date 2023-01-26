package protocol

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"
)

func SendEthernetPacket(handle *pcap.Handle, l ...gopacket.SerializableLayer) error {
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	err := gopacket.SerializeLayers(buf, opts, l...)
	if err != nil {
		return errors.WithStack(err)
	}

	return handle.WritePacketData(buf.Bytes())
}
