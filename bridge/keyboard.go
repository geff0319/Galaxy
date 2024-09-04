package bridge

import (
	hook "github.com/robotn/gohook"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) ExitKey() {
	hook.End()
}

func CreateHook(a *App) {
	go func() {
		//hook.Register(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
		//	fmt.Println("ctrl-shift-q")
		//	hook.End()
		//})
		hook.Register(hook.KeyDown, []string{"alt", "w"}, func(e hook.Event) {
			runtime.EventsEmit(a.Ctx, "keyChangeWidgets")
		})

		s := hook.Start()
		<-hook.Process(s)
	}()
}
