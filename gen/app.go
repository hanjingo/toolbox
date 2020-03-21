package gen

import (
	"strings"
	"sync"
)

type App struct {
	Geners []GenerI
	Config *MainConfig
}

var app *App
var appOnce = new(sync.Once)

func GetApp() *App {
	appOnce.Do(func() {
		app = &App{
			Geners: []GenerI{},
			Config: NewMainConfig(),
		}
	})
	return app
}

func (app *App) Load(addr string) error {
	if err := app.Config.Load(addr); err != nil {
		return err
	}
	for _, cfg := range app.Config.configMap {
		switch strings.ToUpper(cfg.Lang) {
		case LANG_GO:
			app.Geners = append(app.Geners, NewGoGenerator1(cfg))
		case LANG_CSHARP:
			app.Geners = append(app.Geners, NewCsGenerator1(cfg))
		}
	}
	return nil
}

func (app *App) Gen() error {
	for _, gener := range app.Geners {
		if err := gener.Gen(); err != nil {
			return err
		}
	}
	return nil
}
