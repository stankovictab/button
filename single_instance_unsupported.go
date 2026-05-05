//go:build !linux

package main

func handleExistingLinuxInstance(uniqueID string, action launchAction, args []string) bool {
	return false
}
