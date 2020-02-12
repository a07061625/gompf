package mpapp

import (
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/i18n"
)

// 设置应用附加数据
func (app *appBasic) SetConfOther(configs map[string]interface{}) {
    if app.initFlag {
        return
    }

    app.confOther = configs
}

// 设置国际化配置,配置文件路径只能是以./开始,代表从main.go文件所在目录开始,否则会报错
func (app *appBasic) SetConfI18n(conf *i18n.I18n) {
    if app.initFlag {
        return
    }

    app.instance.I18n = conf
}

func (app *appBasic) initConf() {
    for k, v := range app.confOther {
        app.confApp = append(app.confApp, iris.WithOtherValue(k, v))
    }
    app.instance.Configure(app.confApp...)
}
