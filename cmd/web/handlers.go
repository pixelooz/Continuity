package main

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
)

type dummyNoteType struct {
	Heading  string
	Rendered template.HTML
}

func (b *backend) homeView(c echo.Context) error {
	data := struct {
		path string
	}{
		path: c.Path(),
	}
	b.zlog.Info().Str("path", data.path)
	return c.Render(200, "home.gohtml", data)
}

func (b *backend) notesView(c echo.Context) error {
	rendered, err := b.renderer.Markdown(dummyNote)
	if err != nil {
		return c.String(http.StatusNotFound, "something broke mate")
	}
	return c.Render(200, "notes.gohtml", dummyNoteType{
		Heading:  "Some Heading",
		Rendered: rendered,
	})
}
