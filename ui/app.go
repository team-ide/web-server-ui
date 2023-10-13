package ui

import "github.com/team-ide/go-tool/util"

var staticVersion = util.GetUUID()

func NewApp() (app *App) {
	app = &App{}
	app.AddLink(`static/layui-v2.8.16/css/layui.css`)
	app.AddLink(`static/commons/app.css?v=` + staticVersion)

	app.AddScript(`static/monaco-editor/min/vs/loader.js`)
	app.AddScript(`static/monaco-editor/main.js`)
	app.AddScript(`static/layui-v2.8.16/layui.js`)
	app.AddScript(`static/commons/app.js?v=` + staticVersion)
	app.AddScript(`static/commons/tool.js?v=` + staticVersion)
	return
}

type App struct {
	Style
	title    string
	themes   []*Theme
	header   AppHeader
	body     AppBody
	footer   AppFooter
	basePath string

	HtmlHead
}

func (this_ *App) SetBasePath(basePath string) {
	this_.basePath = basePath
}

type AppHeader struct {
	Style
	Left   *AppHeaderPack `json:"left"`
	Center *AppHeaderPack `json:"center"`
	Right  *AppHeaderPack `json:"right"`
}

type AppHeaderPack struct {
	Style
	Menus []*Menu `json:"menus"`
}

type AppBody struct {
	Style
	Left   *AppBodyPack `json:"left"`
	Center *AppBodyPack `json:"center"`
	Right  *AppBodyPack `json:"right"`
}

type AppBodyPack struct {
	Style
	Menus []*Menu `json:"menus"`
}

type AppFooter struct {
	Style
	Left   *AppFooterPack `json:"left"`
	Center *AppFooterPack `json:"center"`
	Right  *AppFooterPack `json:"right"`
}

type AppFooterPack struct {
	Style
	Menus []*Menu `json:"menus"`
}

type Menu struct {
	Style
	Text  string  `json:"text"`
	Href  string  `json:"href"`
	Size  string  `json:"size"`
	Menus []*Menu `json:"menus"`
}

func (this_ *App) Append(options *BuildOptions) (err error) {

	if err = writeHtml(options, `<div class="app-pack">`+"\n"); err != nil {
		return
	}

	options.tab++
	if err = writeHtml(options, `<div class="app">`+"\n"); err != nil {
		return
	}

	options.tab++
	if err = this_.header.Append(options); err != nil {
		return
	}
	if err = this_.body.Append(options); err != nil {
		return
	}
	if err = this_.footer.Append(options); err != nil {
		return
	}

	options.tab--
	if err = writeHtml(options, `</div>`+"\n"); err != nil {
		return
	}

	options.tab--
	if err = writeHtml(options, `</div>`+"\n"); err != nil {
		return
	}
	return
}

func (this_ *AppHeader) Append(options *BuildOptions) (err error) {
	if err = writeHtml(options, `<div class="app-header-pack">`+"\n"); err != nil {
		return
	}

	options.tab++
	if err = writeHtml(options, `<div class="app-header">`+"\n"); err != nil {
		return
	}

	options.tab++
	if err = writeHtml(options, `<div class="app-header-left">`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `</div>`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `<div class="app-header-center">`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `</div>`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `<div class="app-header-right">`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `</div>`+"\n"); err != nil {
		return
	}
	options.tab--

	if err = writeHtml(options, `</div>`+"\n"); err != nil {
		return
	}

	options.tab--
	if err = writeHtml(options, `</div>`+"\n"); err != nil {
		return
	}
	return
}

func (this_ *AppBody) Append(options *BuildOptions) (err error) {
	if err = writeHtml(options, `<div class="app-body-pack">`+"\n"); err != nil {
		return
	}

	options.tab++
	if err = writeHtml(options, `<div class="app-body">`+"\n"); err != nil {
		return
	}

	options.tab++
	if err = writeHtml(options, `<div class="app-body-left">`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `</div>`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `<div class="app-body-center">`+"\n"); err != nil {
		return
	}

	options.tab++
	if err = options.page.Append(options); err != nil {
		return
	}
	options.tab--

	if err = writeHtml(options, `</div>`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `<div class="app-body-right">`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `</div>`+"\n"); err != nil {
		return
	}
	options.tab--

	if err = writeHtml(options, `</div>`+"\n"); err != nil {
		return
	}

	options.tab--
	if err = writeHtml(options, `</div>`+"\n"); err != nil {
		return
	}
	return
}

func (this_ *AppFooter) Append(options *BuildOptions) (err error) {

	if err = writeHtml(options, `<div class="app-footer-pack">`+"\n"); err != nil {
		return
	}

	options.tab++
	if err = writeHtml(options, `<div class="app-footer">`+"\n"); err != nil {
		return
	}

	options.tab++
	if err = writeHtml(options, `<div class="app-footer-left">`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `</div>`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `<div class="app-footer-center">`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `</div>`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `<div class="app-footer-right">`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `</div>`+"\n"); err != nil {
		return
	}
	options.tab--

	if err = writeHtml(options, `</div>`+"\n"); err != nil {
		return
	}

	options.tab--
	if err = writeHtml(options, `</div>`+"\n"); err != nil {
		return
	}
	return
}
