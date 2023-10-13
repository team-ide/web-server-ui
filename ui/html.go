package ui

import "io"

type HtmlHead struct {
	title   string
	links   []*HtmlLink
	scripts []*HtmlScript
}

func (this_ *HtmlHead) SetTitle(title string) {
	this_.title = title
}

func (this_ *HtmlHead) AddLink(href string) {
	this_.links = append(this_.links, &HtmlLink{
		href: href,
	})
}

func (this_ *HtmlHead) AddScript(src string) {
	this_.scripts = append(this_.scripts, &HtmlScript{
		src: src,
	})
}

func (this_ *HtmlHead) appendLinks(options *BuildOptions) (err error) {
	for _, one := range this_.links {
		err = one.append(options)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *HtmlHead) appendScripts(options *BuildOptions) (err error) {
	for _, one := range this_.scripts {
		err = one.append(options)
		if err != nil {
			return
		}
	}
	return
}

type HtmlLink struct {
	href string
}

func (this_ *HtmlLink) append(options *BuildOptions) (err error) {
	app := options.app
	if this_.href != "" {
		if err = writeHtml(options, `<link rel="stylesheet" type="text/css" href="`+app.basePath+this_.href+`"></link>`+"\n"); err != nil {
			return
		}
	}
	return
}

type HtmlScript struct {
	src string
}

func (this_ *HtmlScript) append(options *BuildOptions) (err error) {
	app := options.app
	if this_.src != "" {
		if err = writeHtml(options, `<script type="text/javascript" src="`+app.basePath+this_.src+`"></script>`+"\n"); err != nil {
			return
		}
	}
	return
}

type BuildOptions struct {
	writer   io.Writer
	tab      int
	app      *App
	page     *Page
	onlyPage bool // 只输出 page 的 html
	tabChar  []byte
}

func (this_ *Page) NewPageBuilder(writer io.Writer) (builder *BuildOptions) {
	builder = &BuildOptions{
		app:     this_.app,
		page:    this_,
		writer:  writer,
		tabChar: []byte("\t"),
	}
	return
}

func (this_ *App) BuildHtml(options *BuildOptions) (err error) {
	options.tab = 0
	defer func() {
		options.tab = 0
	}()

	page := options.page
	app := options.app

	var title = app.title
	if page != nil && page.title != "" {
		title = page.title
	}
	if title == "" {
		title = "Web Server UI"
	}

	if err = writeHtml(options, `<!DOCTYPE html>`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `<html lang="zh-cn">`+"\n"); err != nil {
		return
	}

	options.tab++
	if err = writeHtml(options, `<head>`+"\n"); err != nil {
		return
	}

	options.tab++
	if err = writeHtml(options, `<meta charset="utf-8">`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `<meta http-equiv="X-UA-Compatible" content="IE=edge">`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `<meta name="viewport" content="width=device-width,initial-scale=1.0">`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `<title>`+title+`</title>`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `<link rel="icon" href="`+app.basePath+`favicon.png">`+"\n"); err != nil {
		return
	}

	if err = writeHtml(options, `<script type="text/javascript">`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `window.basePath = "`+this_.basePath+`"`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `</script>`+"\n"); err != nil {
		return
	}

	if err = this_.appendLinks(options); err != nil {
		return
	}
	if err = this_.appendScripts(options); err != nil {
		return
	}

	if err = page.appendLinks(options); err != nil {
		return
	}
	if err = page.appendScripts(options); err != nil {
		return
	}

	options.tab--
	if err = writeHtml(options, `</head>`+"\n"); err != nil {
		return
	}
	if err = writeHtml(options, `<body>`+"\n"); err != nil {
		return
	}

	options.tab++

	if options.onlyPage {
		if err = page.Append(options); err != nil {
			return
		}
	} else {
		if err = app.Append(options); err != nil {
			return
		}
	}

	options.tab--
	if err = writeHtml(options, `</body>`+"\n"); err != nil {
		return
	}

	options.tab--
	if err = writeHtml(options, `</html>`+"\n"); err != nil {
		return
	}
	return
}

func writeHtml(options *BuildOptions, html string) (err error) {
	//tab := options.tab
	//for tab > 0 {
	//	tab--
	//	if _, err = options.writer.Write(options.tabChar); err != nil {
	//		return
	//	}
	//}
	if _, err = options.writer.Write([]byte(html)); err != nil {
		return
	}
	return
}
