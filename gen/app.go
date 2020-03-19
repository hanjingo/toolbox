package gen

import (
	"sync"
)

type App struct {
	Geners map[string]GenerI
	Config *MainConfig
}

var app *App
var appOnce = new(sync.Once)

func GetApp() *App {
	appOnce.Do(func(){
		app = &App{
			Geners: make(map[string]GenerI),
			Config: NewMainConfig(),
		}
	})
	return app
}