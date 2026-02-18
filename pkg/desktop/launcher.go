package desktop

import (
	"os/exec"
	"runtime"
)

type AppConfig struct {
	ModuleID      string
	ModuleName    string
	HTTPPort      int
	IconBytes     []byte
	Headless      bool
	TrayOnly      bool
	OnQuit        func()
	TrayMenuItems []TrayMenuItem
}

type TrayMenuItem struct {
	Label   string
	OnClick func()
}

func Launch(cfg AppConfig) {
	if cfg.Headless {
		select {}
	}

	if cfg.TrayOnly {
		StartTray(cfg)
		return
	}

	go StartTray(cfg)
	OpenWindow(cfg)

	select {}
}

func OpenBrowser(url string) {
	switch runtime.GOOS {
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	default:
		exec.Command("xdg-open", url).Start()
	}
}
