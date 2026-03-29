package main

import (
	"Continuity/internal/data"
	"Continuity/internal/validate"
	"errors"
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type createNoteForm struct {
	CollectionID       string `form:"collection_id"`
	Title              string `form:"title"`
	Content            string `form:"content"`
	validate.Validator `form:"_"`
}

type createNoteData struct {
	CollectionID uuid.UUID
	IsRoot       bool
}

type viewNoteData struct {
	Heading  string
	Rendered template.HTML
}

func (b *backend) viewRootCreateNoteForm(c echo.Context) error {
	root := b.contextGetRootCollection(c)

	noteData := &createNoteData{
		CollectionID: root.ID,
		IsRoot:       root.IsRoot,
	}
	noteForm := new(createNoteForm)

	pd := b.NewPageData(c)
	pd.Form = noteForm
	pd.Data = noteData
	return c.Render(http.StatusOK, "create_note.gohtml", pd)
}

func (b *backend) postCreateNoteFormDummy(c echo.Context) error {
	var noteForm createNoteForm

	if err := b.decodePostForm(c, &noteForm); err != nil {
		return b.renderBadRequestErr(c, "invalid form data")
	}
	pID, err := uuid.Parse(noteForm.CollectionID)
	if err != nil {
		return b.renderBadRequestErr(c, "invalid parent id")
	}
	userID := b.contextGetUser(c).ID

	note := &data.Note{
		ID: uuid.New(), UserId: userID,
		CollectionId: pID,
		Title:        noteForm.Title,
		Content:      noteForm.Content,
	}
	validate.ValidateNote(&noteForm.Validator, note)
	if !noteForm.Valid() {
		return b.renderWithFieldErr(c, "notes.gohtml", &noteForm)
	}
	rendered, err := b.renderer.Markdown(note.Content)
	if err != nil {
		return b.renderInternalServerErr(c)
	}
	noteView := &viewNoteData{
		Heading: note.Title, Rendered: rendered,
	}
	pd := b.NewPageData(c)
	pd.Form = noteView
	return c.Render(http.StatusOK, "notes.gohtml", pd)
}

func (b *backend) postCreateNoteForm(c echo.Context) error {
	var noteForm createNoteForm

	if err := b.decodePostForm(c, &noteForm); err != nil {
		return b.renderBadRequestErr(c, "invalid form data")
	}
	pID, err := uuid.Parse(noteForm.CollectionID)
	if err != nil {
		return b.renderBadRequestErr(c, "invalid parent id")
	}
	userID := b.contextGetUser(c).ID

	note := &data.Note{
		ID: uuid.New(), UserId: userID,
		CollectionId: pID,
		Title:        noteForm.Title,
		Content:      noteForm.Content,
	}
	validate.ValidateNote(&noteForm.Validator, note)
	if !noteForm.Valid() {
		return b.renderWithFieldErr(c, "notes.gohtml", &noteForm)
	}
	if err = b.models.Note.Insert(c.Request().Context(), note); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateNote):
			noteForm.AddFieldError("name", "title already exists")
			return b.renderWithFieldErr(c, "notes.gohtml", &noteForm)
		default:
			b.zlog.Err(err).
				Str("handler", "postCreateCollectionForm").
				Msg("insert collection into db")
			return b.renderInternalServerErr(c)
		}
	}
	return c.Redirect(http.StatusOK, "/v1/collection/"+noteForm.CollectionID)
}
