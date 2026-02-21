package desktop

import (
	"github.com/gen2brain/beeep"
)

// NotifyTrayMinimized показывает системное уведомление о том, что приложение продолжает работать в трее.
func NotifyTrayMinimized(appName string) {
	_ = beeep.Notify(appName, "Приложение свернуто в трей. Иконка в области уведомлений.", "")
}
