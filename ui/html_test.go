package ui

import (
	"fmt"
	"strings"
	"testing"
)

func TestHtml(t *testing.T) {
	app := &App{
		Header: &AppHeader{},
		Body:   &AppBody{},
		Footer: &AppFooter{},
	}

	builder, err := NewHtmlBuilder()
	if err != nil {
		panic(err)
	}
	builder.App = app
	builder.BasePath = "/"
	page := &Page{}

	writer := &strings.Builder{}
	err = builder.OutHtml(writer, page)
	if err != nil {
		panic(err)
	}
	fmt.Println(writer.String())
}
