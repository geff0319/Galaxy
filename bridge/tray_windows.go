//go:build windows

package bridge

import (
	"embed"
	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"os"
	sysruntime "runtime"
)

var assets embed.FS

//func CreateTray(a *App, icon []byte, fs embed.FS) {
//	assets = fs
//	go func() {
//		fmt.Println("Tray start")
//		defer fmt.Println("Tray out")
//		systray.Run(func() {
//			systray.SetIcon([]byte(icon))
//			systray.SetTitle("Galaxy")
//			systray.SetTooltip("Galaxy")
//			systray.SetOnDClick(func(menu systray.IMenu) {
//				fmt.Println("SetOnDClick")
//				runtime.WindowShow(a.Ctx)
//			})
//			systray.SetOnRClick(func(menu systray.IMenu) {
//				fmt.Println("SetOnRClick")
//				menu.ShowMenu()
//			})
//		}, nil)
//	}()
//}

func InitTray(a *App, icon []byte, fs embed.FS) {
	go func() {
		sysruntime.LockOSThread()
		defer sysruntime.UnlockOSThread()
		systray.Run(func() {
			systray.SetIcon([]byte(icon))
			systray.SetTitle("Galaxy")
			systray.SetTooltip("Galaxy")
			systray.SetOnClick(func(menu systray.IMenu) { runtime.WindowShow(a.Ctx) })
			systray.SetOnRClick(func(menu systray.IMenu) { menu.ShowMenu() })

			// Ensure the tray is still available if rolling-release fails
			mRestart := systray.AddMenuItem("重启", "重启")
			mExit := systray.AddMenuItem("退出", "退出")
			mRestart.Click(func() { a.RestartApp() })
			mExit.Click(func() { a.ExitApp() })
		}, nil)
	}()
}

func (a *App) UpdateTray(tray TrayContent) {
	if tray.Icon != "" {
		ico, _ := assets.ReadFile("frontend/dist/" + tray.Icon)
		systray.SetIcon(ico)
	}
	if tray.Title != "" {
		systray.SetTitle(tray.Title)
	}
	if tray.Tooltip != "" {
		systray.SetTooltip(tray.Tooltip)
	}
}

func (a *App) UpdateTrayMenus(menus []MenuItem) {
	log.Printf("UpdateTrayMenus")

	systray.ResetMenu()

	for _, menu := range menus {
		createMenuItem(menu, a, nil)
	}
}

func createMenuItem(menu MenuItem, a *App, parent *systray.MenuItem) {
	if menu.Hidden {
		return
	}
	switch menu.Type {
	case "item":
		var m *systray.MenuItem
		if parent == nil {
			m = systray.AddMenuItem(menu.Text, menu.Tooltip)
		} else {
			m = parent.AddSubMenuItem(menu.Text, menu.Tooltip)
		}
		m.Click(func() { runtime.EventsEmit(a.Ctx, menu.Event) })

		if menu.Checked {
			m.Check()
		}

		for _, child := range menu.Children {
			createMenuItem(child, a, m)
		}
	case "separator":
		systray.AddSeparator()
	}
}

func (a *App) ExitApp() {
	systray.Quit()
	runtime.Quit(a.Ctx)
	os.Exit(0)
}

//func (a *App) RestartApp() FlagResult {
//	exePath := Env.BasePath + "\\" + Env.AppName
//
//	cmd := exec.Command(exePath)
//	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
//
//	err := cmd.Start()
//	if err != nil {
//		return FlagResult{false, err.Error()}
//	}
//
//	a.ExitApp()
//
//	return FlagResult{true, "Success"}
//}
