package main

import (
	"context"
	"embed"
	"wails_send/backend"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
    // 创建 Sender 和 Receiver 实例
    app := backend.NewSender()


    // 创建应用程序
    err := wails.Run(&options.App{
        Title:  "wails_send",
        Width:  640,
        Height: 730,
        Windows: &windows.Options{
            WebviewIsTransparent: true,
            WindowIsTranslucent:  true,
            BackdropType: windows.Acrylic,
        },
        DisableResize:true,
        // Frameless: true,
        AssetServer: &assetserver.Options{
            Assets: assets,
        },
        OnStartup: func(ctx context.Context) {
            // 初始化 Sender 和 Receiver 的 context
            app.Startup(ctx)
        },
        OnShutdown: func(ctx context.Context) {
            // 清理资源
            app.Shutdown(ctx)
        },
        Bind: []interface{}{
            app,
        },
    })

    if err != nil {
        println("Error:", err.Error())
    }
}
