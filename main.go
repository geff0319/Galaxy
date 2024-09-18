package main

import (
	"context"
	"embed"
	"galaxy/bridge"
	"github.com/ge-fei-fan/gefflog"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed frontend/dist/favicon.ico
var icon []byte

func main() {
	// Create an instance of the app structure
	bridge.InitBridge()

	app := bridge.NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  bridge.AppTitle,
		Width:  1000,
		Height: 640,
		//MinWidth:      50,
		//MinHeight:     50,
		Debug:         options.Debug{true},
		Frameless:     bridge.Env.OS == "windows",
		DisableResize: true,
		StartHidden: func() bool {
			if bridge.Env.FromTaskSch {
				return bridge.Config.WindowStartState == 2
			}
			return false
		}(),
		WindowStartState: func() options.WindowStartState {
			if bridge.Env.FromTaskSch {
				return options.WindowStartState(bridge.Config.WindowStartState)
			}
			return 0
		}(),
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
			//BackdropType:         windows.Acrylic,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			Appearance:           mac.DefaultAppearance,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
			About: &mac.AboutInfo{
				Title:   "Galaxy",
				Message: "© 2024 Galaxy",
				Icon:    icon,
			},
		},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 100},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId:               "Galaxy",
			OnSecondInstanceLaunch: app.OnSecondInstanceLaunch,
		},
		OnStartup: func(ctx context.Context) {
			app.Ctx = ctx
			bridge.CreateHook(app)
			bridge.MqNotifyConsumer(app)
			bridge.InitTray(app, icon, assets)
			bridge.InitScheduledTasks()
		},
		//OnBeforeClose: func(ctx context.Context) (prevent bool) {
		//	runtime.EventsEmit(ctx, "beforeClose")
		//	return true
		//},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		gefflog.Err("Error:", err.Error())
	}
}
