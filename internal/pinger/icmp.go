package pinger

import (
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// SendICMPEcho send ICMP echo request and wait for reply.
func SendICMPEcho(id, seq int, addr net.UDPAddr, timeout time.Duration) (net.Addr, *icmp.Message, error) {
	deadline := time.After(timeout)
	conn, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	defer conn.Close()

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   id,
			Seq:  seq,
			Data: []byte("nf-icmp-ping"),
		},
	}

	mb, err := msg.Marshal(nil)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	_, err = conn.WriteTo(mb, &addr)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	var result = make(chan struct {
		addr net.Addr
		msg  *icmp.Message
		err  error
	}, 1)
	go func() {
		defer close(result)

		rb := make([]byte, 1500)
		_, peer, err := conn.ReadFrom(rb)
		if err != nil {
			result <- struct {
				addr net.Addr
				msg  *icmp.Message
				err  error
			}{
				nil, nil, errors.WithStack(err),
			}
		}

		rm, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), rb)
		if err != nil {
			result <- struct {
				addr net.Addr
				msg  *icmp.Message
				err  error
			}{
				nil, nil, errors.WithStack(err),
			}
		}

		result <- struct {
			addr net.Addr
			msg  *icmp.Message
			err  error
		}{
			peer, rm, nil,
		}
	}()

	select {
	case r := <-result:
		if r.err != nil {
			return nil, nil, err
		}
		return r.addr, r.msg, nil
	case <-deadline:
		return nil, nil, fmt.Errorf("read timeout")
	}
}
