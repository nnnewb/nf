package elevation

import (
	"github.com/pkg/errors"
	"log"
	"os"
	"os/user"
	"syscall"
)

func isElevated() bool {
	u, err := user.Current()
	if err != nil {
		log.Fatalf("get current user failed, error %+v", err)
	}

	return u.Uid == "0"
}

func runElevated() error {
	exe, err := os.Executable()
	if err != nil {
		return errors.WithStack(err)
	}
	return syscall.Exec("sudo", append([]string{exe}, os.Args...), os.Environ())
}
