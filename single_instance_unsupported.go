//go:build !linux || bindings

package main

func handleExistingLinuxInstance(uniqueID string, action launchAction, args []string) (bool, error) {
	return false, nil
}
