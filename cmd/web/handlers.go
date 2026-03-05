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
	pd := b.NewPageData(c)
	return c.Render(http.StatusOK, "home.gohtml", pd)
}

func (b *backend) notesView(c echo.Context) error {
	var notePage dummyNoteType
	rendered, err := b.renderer.Markdown(dummyNote)
	if err != nil {
		return c.String(http.StatusNotFound, "something broke mate")
	}
	notePage = dummyNoteType{
		Heading: "Some Heading", Rendered: rendered,
	}
	pd := b.NewPageData(c)
	pd.Form = notePage
	return c.Render(http.StatusOK, "notes.gohtml", pd)
}

func (b *backend) demoView(c echo.Context) error {
	return c.Render(http.StatusInternalServerError, "perma_error.gohtml", Error{
		ErrType: "Internal Server Error",
		ErrCode: 500,
		Message: "something went wrong",
	})
}
