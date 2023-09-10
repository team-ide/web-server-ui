package ui

import "io"

func NewHtmlBuilder() (builder *HtmlBuilder, err error) {
	builder = &HtmlBuilder{
		Links: []*HtmlLink{
			{Href: `static/layui-v2.8.16/css/layui.css`},
			{Href: `static/commons/index.css`},
		},
		Scripts: []*HtmlScript{
			{Src: `static/monaco-editor/min/vs/loader.js`},
			{Src: `static/monaco-editor/main.js`},
			{Src: `static/layui-v2.8.16/layui.js`},
			{Src: `static/commons/index.js`},
			{Src: `static/commons/tool.js`},
		},
	}
	return
}

type HtmlBuilder struct {
	App      *App          `json:"app"`
	BasePath string        `json:"basePath"`
	Links    []*HtmlLink   `json:"links"`
	Scripts  []*HtmlScript `json:"scripts"`
}

type HtmlLink struct {
	Href string `json:"href"`
}

type HtmlScript struct {
	Src string `json:"src"`
}

func (this_ *HtmlBuilder) OutHtml(writer io.Writer, page *Page) (err error) {
	var title string
	if this_.App != nil {
		title = this_.App.Title
	}
	if page != nil && page.Title != "" {
		title = page.Title
	}
	var tab = 0
	if err = writeHtml(writer, tab, `<!DOCTYPE html>`+"\n"); err != nil {
		return
	}
	if err = writeHtml(writer, tab, `<html lang="">`+"\n"); err != nil {
		return
	}
	if err = writeHtml(writer, tab+1, `<head>`+"\n"); err != nil {
		return
	}
	if err = writeHtml(writer, tab+2, `<meta charset="utf-8">`+"\n"); err != nil {
		return
	}
	if err = writeHtml(writer, tab+2, `<meta http-equiv="X-UA-Compatible" content="IE=edge">`+"\n"); err != nil {
		return
	}
	if err = writeHtml(writer, tab+2, `<meta name="viewport" content="width=device-width,initial-scale=1.0">`+"\n"); err != nil {
		return
	}
	if err = writeHtml(writer, tab+2, `<title>`+title+`</title>`+"\n"); err != nil {
		return
	}
	if err = writeHtml(writer, tab+2, `<link rel="icon" href="`+this_.BasePath+`favicon.png">`+"\n"); err != nil {
		return
	}
	if err = this_.AppendLinks(writer, tab+2, this_.Links); err != nil {
		return
	}
	if err = this_.AppendScripts(writer, tab+2, this_.Scripts); err != nil {
		return
	}

	if err = writeHtml(writer, tab+1, `</head>`+"\n"); err != nil {
		return
	}
	if err = writeHtml(writer, tab+1, `<body>`+"\n"); err != nil {
		return
	}
	if this_.App != nil {
		if err = this_.App.OutHtml(writer, tab+2); err != nil {
			return
		}
	}

	if err = writeHtml(writer, tab+1, `</body>`+"\n"); err != nil {
		return
	}
	if err = writeHtml(writer, tab, `</html>`+"\n"); err != nil {
		return
	}
	return
}

func (this_ *HtmlBuilder) AppendLinks(writer io.Writer, tab int, links []*HtmlLink) (err error) {
	for _, one := range links {
		if one.Href == "" {
			continue
		}
		if err = writeHtml(writer, tab, `<link rel="stylesheet" type="text/css" href="`+this_.BasePath+one.Href+`" />`+"\n"); err != nil {
			return
		}
	}
	return
}

func (this_ *HtmlBuilder) AppendScripts(writer io.Writer, tab int, scripts []*HtmlScript) (err error) {
	for _, one := range scripts {
		if one.Src == "" {
			continue
		}
		if err = writeHtml(writer, tab, `<script type="text/javascript" src="`+this_.BasePath+one.Src+`"></script>`+"\n"); err != nil {
			return
		}
	}
	return
}
