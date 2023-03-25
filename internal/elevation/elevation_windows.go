package elevation

import (
	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
	"log"
	"os"
	"strings"
	"syscall"
)

func isElevated() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		log.Printf("privilege test failed, assume current user not have privilege, error %+v", err)
		return false
	}
	return true
}

func runElevated() error {
	verb := "runas"
	exe, err := os.Executable()
	if err != nil {
		return errors.WithStack(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		return errors.WithStack(err)
	}
	args := strings.Join(os.Args[1:], " ")

	verbPtr, err := syscall.UTF16PtrFromString(verb)
	if err != nil {
		return errors.WithStack(err)
	}
	exePtr, err := syscall.UTF16PtrFromString(exe)
	if err != nil {
		return errors.WithStack(err)
	}
	cwdPtr, err := syscall.UTF16PtrFromString(cwd)
	if err != nil {
		return errors.WithStack(err)
	}
	argPtr, err := syscall.UTF16PtrFromString(args)
	if err != nil {
		return errors.WithStack(err)
	}

	var showCmd int32 = 1 //SW_NORMAL

	err = windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
