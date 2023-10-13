package ui

func (this_ *App) NewPage() (page *Page) {
	page = &Page{
		app: this_,
	}
	return
}

type Page struct {
	app *App
	HtmlHead
}

func (this_ *Page) Append(options *BuildOptions) (err error) {
	if err = writeHtml(options, `<div class="app-page-pack">`+"\n"); err != nil {
		return
	}

	options.tab++
	if err = writeHtml(options, `<div class="app-page">`+"\n"); err != nil {
		return
	}

	options.tab++
	if err = this_.AppendHeader(options); err != nil {
		return
	}
	if err = this_.AppendBody(options); err != nil {
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

func (this_ *Page) AppendHeader(options *BuildOptions) (err error) {
	if err = writeHtml(options, `<div class="page-header-pack">`+"\n"); err != nil {
		return
	}

	options.tab++
	if err = writeHtml(options, `<div class="page-header">`+"\n"); err != nil {
		return
	}

	if err = writeHtml(options, `这是标题`+"\n"); err != nil {
		return
	}

	if err = writeHtml(options, `</div>`+"\n"); err != nil {
		return
	}

	options.tab--
	if err = writeHtml(options, `</div>`+"\n"); err != nil {
		return
	}
	return
}

func (this_ *Page) AppendBody(options *BuildOptions) (err error) {
	if err = writeHtml(options, `<div class="page-body-pack">`+"\n"); err != nil {
		return
	}

	options.tab++
	if err = writeHtml(options, `<div class="page-body">`+"\n"); err != nil {
		return
	}

	if err = writeHtml(options, `这是内容`+"\n"); err != nil {
		return
	}

	if err = writeHtml(options, `</div>`+"\n"); err != nil {
		return
	}

	options.tab--
	if err = writeHtml(options, `</div>`+"\n"); err != nil {
		return
	}
	return
}
