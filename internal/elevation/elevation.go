package elevation

// IsElevated is current user privileged. On windows, test by open \\.\PHYSICALDRIVE0.
// On Linux/Unix like platform, test by check user.Current Uid.
func IsElevated() bool {
	return isElevated()
}

// RunElevated use `sudo` command execute current command by privileged user.
func RunElevated() error {
	return runElevated()
}
