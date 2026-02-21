package desktop

import (
	"fmt"
	"sync"

	webview "github.com/webview/webview_go"
)

var (
	currentWindow webview.WebView
	windowMu      sync.Mutex
)

func OpenWindow(cfg AppConfig) {
	windowMu.Lock()
	if currentWindow != nil {
		windowMu.Unlock()
		return
	}

	w := webview.New(false)
	currentWindow = w
	windowMu.Unlock()

	defer func() {
		w.Destroy()
		windowMu.Lock()
		currentWindow = nil
		windowMu.Unlock()
	}()

	w.SetTitle(cfg.ModuleName)
	w.SetSize(1280, 800, webview.HintNone)
	w.SetSize(800, 500, webview.HintMin)

	url := fmt.Sprintf("http://localhost:%d", cfg.HTTPPort)
	w.Navigate(url)

	w.Run()
	// Окно закрыто — приложение остаётся в трее; показываем уведомление.
	NotifyTrayMinimized(cfg.ModuleName)
}

func ShowWindow(cfg AppConfig) {
	windowMu.Lock()
	hasWindow := currentWindow != nil
	windowMu.Unlock()

	if hasWindow {
		return
	}

	go OpenWindow(cfg)
}
