package pinger

import (
	"fmt"
	"net"
	"time"

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
