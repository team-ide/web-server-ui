package ui

import (
	"fmt"
	"strings"
	"testing"
)

func TestHtml(t *testing.T) {
	var err error
	app := NewApp()
	app.basePath = "/"

	page := app.NewPage()

	writer := &strings.Builder{}
	builder := page.NewPageBuilder(writer)

	err = app.BuildHtml(builder)
	if err != nil {
		panic(err)
	}
	fmt.Println(writer.String())

	writer = &strings.Builder{}
	builder = page.NewPageBuilder(writer)
	builder.onlyPage = true

	err = app.BuildHtml(builder)
	if err != nil {
		panic(err)
	}
	fmt.Println(writer.String())
}
