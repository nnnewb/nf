//go:build mage
// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"runtime"
	"strings"
	"time"
)

const APP_VERSION = "0.1.0"

func makeBuildCommit() (string, error) {
	output, err := sh.Output("git", "--no-pager", "log", "--pretty=format:%H", "-1")
	if err != nil {
		return "", err
	}
	commit := strings.TrimSpace(output)

	output, err = sh.Output("git", "status", "--porcelain")
	if err != nil {
		return "", err
	}
	var dirty string
	if len(strings.TrimSpace(output)) > 0 {
		dirty = "-dirty"
	}
	return commit + dirty, nil
}

func makeBuildTime() string {
	return time.Now().Format(time.RFC3339)
}

func Build() error {
	commit, err := makeBuildCommit()
	if err != nil {
		return err
	}

	outputFilename := "bin/nf"
	if runtime.GOOS == "windows" {
		outputFilename += ".exe"
	}

	ldFlags := []string{"-s", "-w"}
	ldFlags = append(ldFlags, "-X", fmt.Sprintf("github.com/nnnewb/nf/internal/constants.BUILD_COMMIT=%s", commit))
	ldFlags = append(ldFlags, "-X", fmt.Sprintf("github.com/nnnewb/nf/internal/constants.BUILD_TIME=%s", makeBuildTime()))
	ldFlags = append(ldFlags, "-X", fmt.Sprintf("github.com/nnnewb/nf/internal/constants.VERSION=%s", APP_VERSION))

	err = sh.RunV("go", "build", "-ldflags", strings.Join(ldFlags, " "), "-o", outputFilename, "main.go")
	if err != nil {
		return err
	}
	return nil
}
