package main

import (
	"Continuity/internal/data"
	"Continuity/internal/validate"
	"errors"
	"html/template"
	"net/http"
	"time"

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
	ID           uuid.UUID
	CollectionID uuid.UUID
	Heading      string
	Rendered     template.HTML
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// todo: make this work for other pages as well, listen on the same url which
//
//	also asks for container id but check if no id was provided through
//	the url, if it wasn't then use the root url and other wise use the url provided.

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

func (b *backend) viewCreateNoteForm(c echo.Context) error {
	colxnID, err := uuid.Parse(c.Param("collection-id"))
	if err != nil {
		b.zlog.Err(err).
			Str("handler", "viewCreateNoteForm").
			Msg("failed to parse collection id")
		return b.renderInternalServerErr(c)
	}
	root := b.contextGetRootCollection(c)
	noteData := &createNoteData{
		CollectionID: colxnID, IsRoot: root.ID == colxnID,
	}
	noteForm := new(createNoteForm)

	pd := b.NewPageData(c)
	pd.Form = noteForm
	pd.Data = noteData
	return c.Render(http.StatusOK, "create_note.gohtml", pd)
}

func (b *backend) postCreateNoteForm(c echo.Context) error {
	var noteForm createNoteForm
	root := b.contextGetRootCollection(c)

	if err := b.decodePostForm(c, &noteForm); err != nil {
		return b.renderBadRequestErr(c, "invalid form data")
	}
	colxnID, err := uuid.Parse(noteForm.CollectionID)
	if err != nil {
		return b.renderBadRequestErr(c, "invalid parent id")
	}
	userID := b.contextGetUser(c).ID

	note := &data.Note{
		ID: uuid.New(), UserId: userID,
		CollectionId: colxnID,
		Title:        noteForm.Title,
		Content:      noteForm.Content,
	}
	validate.ValidateNote(&noteForm.Validator, note)
	if !noteForm.Valid() {
		return b.renderErrFormData(c, "create_note.gohtml",
			&noteForm,
			&createNoteData{
				CollectionID: colxnID, IsRoot: root.ID == colxnID,
			},
		)
	}
	if err = b.models.Note.Insert(c.Request().Context(), note); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateNote):
			noteForm.AddFieldError("name", "title already exists")
			return b.renderErrFormData(c, "create_note.gohtml",
				&noteForm,
				&createNoteData{
					CollectionID: colxnID, IsRoot: root.ID == colxnID,
				},
			)
		default:
			b.zlog.Err(err).
				Str("handler", "postCreateCollectionForm").
				Msg("insert collection into db")
			return b.renderInternalServerErr(c)
		}
	}
	return c.Redirect(http.StatusSeeOther, "/v1/note/"+note.ID.String())
}

func (b *backend) viewNotePage(c echo.Context) error {
	noteID, err := uuid.Parse(c.Param("note-id"))
	if err != nil {
		b.zlog.Err(err).
			Str("handler", "viewNotePage").
			Msg("failed to parse note-id from parameter")
		return b.renderInternalServerErr(c)
	}
	note, err := b.models.Note.GetByID(c.Request().Context(), noteID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return b.renderNotFoundErr(c, "note not found")
		default:
			return b.renderInternalServerErr(c)
		}
	}
	rendered, err := b.renderer.Markdown(note.Content)
	if err != nil {
		return b.renderInternalServerErr(c)
	}
	noteView := &viewNoteData{
		ID: note.ID, CollectionID: note.CollectionId,
		Heading: note.Title, Rendered: rendered,
		CreatedAt: note.CreatedAt, UpdatedAt: note.UpdatedAt,
	}
	pd := b.NewPageData(c)
	pd.Data = noteView
	return c.Render(http.StatusOK, "notes.gohtml", pd)
}
