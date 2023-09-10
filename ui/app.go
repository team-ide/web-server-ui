package ui

import (
	"io"
)

type App struct {
	Style
	Title  string     `json:"title"`
	Themes []*Theme   `json:"themes"`
	Header *AppHeader `json:"header"`
	Body   *AppBody   `json:"body"`
	Footer *AppFooter `json:"footer"`
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

func (this_ *App) OutHtml(writer io.Writer, tab int) (err error) {
	if err = writeHtml(writer, tab, `<div class="app-pack">`+"\n"); err != nil {
		return
	}
	if err = writeHtml(writer, tab+1, `<div class="app">`+"\n"); err != nil {
		return
	}

	if this_.Header != nil {
		if err = this_.Header.OutHtml(writer, tab+2); err != nil {
			return
		}
	}
	if this_.Body != nil {
		if err = this_.Body.OutHtml(writer, tab+2); err != nil {
			return
		}
	}
	if this_.Footer != nil {
		if err = this_.Footer.OutHtml(writer, tab+2); err != nil {
			return
		}
	}

	if err = writeHtml(writer, tab+1, `</div>`+"\n"); err != nil {
		return
	}
	if err = writeHtml(writer, tab, `</div>`+"\n"); err != nil {
		return
	}
	return
}

func (this_ *AppHeader) OutHtml(writer io.Writer, tab int) (err error) {
	if err = writeHtml(writer, tab, `<div class="app-header-pack">`+"\n"); err != nil {
		return
	}

	if err = writeHtml(writer, tab+1, `<div class="app-header">`+"\n"); err != nil {
		return
	}

	if err = writeHtml(writer, tab+1, `</div>`+"\n"); err != nil {
		return
	}

	if err = writeHtml(writer, tab, `</div>`+"\n"); err != nil {
		return
	}
	return
}

func (this_ *AppBody) OutHtml(writer io.Writer, tab int) (err error) {
	if err = writeHtml(writer, tab, `<div class="app-body-pack">`+"\n"); err != nil {
		return
	}

	if err = writeHtml(writer, tab+1, `<div class="app-body">`+"\n"); err != nil {
		return
	}

	if err = writeHtml(writer, tab+1, `</div>`+"\n"); err != nil {
		return
	}

	if err = writeHtml(writer, tab, `</div>`+"\n"); err != nil {
		return
	}
	return
}

func (this_ *AppFooter) OutHtml(writer io.Writer, tab int) (err error) {
	if err = writeHtml(writer, tab, `<div class="app-footer-pack">`+"\n"); err != nil {
		return
	}

	if err = writeHtml(writer, tab+1, `<div class="app-footer">`+"\n"); err != nil {
		return
	}

	if err = writeHtml(writer, tab+1, `</div>`+"\n"); err != nil {
		return
	}

	if err = writeHtml(writer, tab, `</div>`+"\n"); err != nil {
		return
	}
	return
}

func writeHtml(writer io.Writer, tab int, html string) (err error) {
	for tab > 0 {
		tab--
		if _, err = writer.Write([]byte("\t")); err != nil {
			return
		}
	}
	if _, err = writer.Write([]byte(html)); err != nil {
		return
	}
	return
}
