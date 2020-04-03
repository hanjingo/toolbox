package gen

import (
	"strings"
	"sync"
)

type App struct {
	Geners      []GenerI    //生成器
	Config      *MainConfig //总配置
	startIdx    int         //起始id
	startErrIdx int         //起始错误id
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
	//读单个文件
	for _, cfg := range app.Config.FileMap {
		temp := NewConfig()
		if err := LoadJsonConfig(cfg.Path, temp); err != nil {
			continue
		}
		classes := []*ClassConfig{}
		for _, v := range temp.Classes {
			if value, ok := app.Config.classMap[v.Name]; ok { //防止重复
				classes = append(classes, value)
			} else {
				app.autoIncr(cfg.PathMap, v)
				app.Config.classMap[v.Name] = v
				classes = append(classes, v)
			}
		}
		switch strings.ToUpper(cfg.Lang) {
		case LANG_GO_V1:
			app.Geners = append(app.Geners, NewGoGenerator1(classes, cfg.PathMap, cfg.NameSpaceMap))
		case LANG_CSHARP_V1:
			app.Geners = append(app.Geners, NewCsGenerator1(classes, cfg.PathMap, cfg.NameSpaceMap))
		case LANG_DOC_V1:
			app.Geners = append(app.Geners, NewDocGenerator1(classes, cfg.PathMap, cfg.NameSpaceMap))
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

func (app *App) autoIncr(fileMap map[string]string, cfgs ...*ClassConfig) {
	for _, cfg := range cfgs {
		if isPrintId(fileMap) {
			if cfg.Id != 0 {
				app.startIdx = cfg.Id
			} else {
				app.startIdx++
			}
			cfg.Id = app.startIdx
		}
		if isPrintErr(fileMap) {
			if cfg.Id != 0 {
				app.startErrIdx = cfg.Id
			} else {
				app.startErrIdx++
			}
			cfg.Id = app.startErrIdx
		}
	}
}
