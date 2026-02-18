package desktop

import (
	"fmt"
	"os"

	"fyne.io/systray"
)

func StartTray(cfg AppConfig) {
	systray.Run(
		func() { onTrayReady(cfg) },
		func() {},
	)
}

func onTrayReady(cfg AppConfig) {
	systray.SetIcon(cfg.IconBytes)
	systray.SetTitle(cfg.ModuleName)
	systray.SetTooltip(fmt.Sprintf("%s — port %d", cfg.ModuleName, cfg.HTTPPort))

	mOpen := systray.AddMenuItem("Open "+cfg.ModuleName, "Open the dashboard window")
	mBrowser := systray.AddMenuItem("Open in Browser", "Open in default browser")

	systray.AddSeparator()

	var moduleClicks []<-chan struct{}
	for _, item := range cfg.TrayMenuItems {
		mi := systray.AddMenuItem(item.Label, "")
		moduleClicks = append(moduleClicks, mi.ClickedCh)
	}

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit "+cfg.ModuleName)

	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				ShowWindow(cfg)
			case <-mBrowser.ClickedCh:
				OpenBrowser(fmt.Sprintf("http://localhost:%d", cfg.HTTPPort))
			case <-mQuit.ClickedCh:
				if cfg.OnQuit != nil {
					cfg.OnQuit()
				}
				systray.Quit()
				os.Exit(0)
			}
		}
	}()

	for i, ch := range moduleClicks {
		idx := i
		go func() {
			for range ch {
				if idx < len(cfg.TrayMenuItems) {
					cfg.TrayMenuItems[idx].OnClick()
				}
			}
		}()
	}
}
