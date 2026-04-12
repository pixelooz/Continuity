package main

import (
	"Continuity/internal/data"
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
)

type dummyNoteType struct {
	Heading  string
	Rendered template.HTML
}

// todo: it's corresponding type which controls actions for the home page will be a diff data type.
// homeViewData will contain data that should be displayed on the home page.
type homeViewData struct {
	Notes []*data.Note
}

func (b *backend) recentNotesView(c echo.Context) error {
	recentNotes, err := b.models.Note.GetRecentlyModified(c.Request().Context(), 10, 0)
	if err != nil {
		b.zlog.Err(err).
			Str("handler", "homeView").
			Msg("failed to get recently modified notes")
		return b.renderInternalServerErr(c)
	}
	homeData := &homeViewData{Notes: recentNotes}
	pd := b.NewPageData(c)
	pd.Data = homeData
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
