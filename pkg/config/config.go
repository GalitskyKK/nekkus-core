package config

import (
	"os"
	"path/filepath"
	"runtime"
)

func GetDataDir(moduleID string) string {
	var base string
	switch runtime.GOOS {
	case "windows":
		base = os.Getenv("APPDATA")
	case "darwin":
		base = filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
	default:
		base = filepath.Join(os.Getenv("HOME"), ".config")
	}
	dir := filepath.Join(base, "nekkus", moduleID)
	os.MkdirAll(dir, 0755)
	return dir
}

func GetLogDir(moduleID string) string {
	dir := filepath.Join(GetDataDir(moduleID), "logs")
	os.MkdirAll(dir, 0755)
	return dir
}
